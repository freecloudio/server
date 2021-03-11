package manager

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/sirupsen/logrus"
)

type UserManager interface {
	CreateUser(authCtx *authorization.Context, user *models.User) *fcerror.Error
	GetUserByID(authCtx *authorization.Context, userID models.UserID) (*models.User, *fcerror.Error)
	UpdateUser(authCtx *authorization.Context, userID models.UserID, updateUser *models.UserUpdate) (*models.User, *fcerror.Error)
	CountUsers(authCtx *authorization.Context) (int64, *fcerror.Error)
	Close()
}

func NewUserManager(cfg config.Config, userPersistence persistence.UserPersistenceController) UserManager {
	userMgr := &userManager{
		cfg:             cfg,
		userPersistence: userPersistence,
		done:            make(chan struct{}),
	}

	return userMgr
}

type userManager struct {
	cfg             config.Config
	userPersistence persistence.UserPersistenceController
	done            chan struct{}
}

func (mgr *userManager) Close() {
	mgr.done <- struct{}{}
}

func (mgr *userManager) CreateUser(authCtx *authorization.Context, user *models.User) (fcerr *fcerror.Error) {
	// TODO: Input Validation

	trans, fcerr := mgr.userPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}

	existingUser, fcerr := trans.GetUserByEmail(user.Email)
	if fcerr != nil && fcerr.ID != fcerror.ErrUserNotFound {
		logrus.WithError(fcerr).Error("Could not verify if user with this email already exists")
		trans.Rollback()
		return
	} else if fcerr == nil && existingUser != nil {
		fcerr = fcerror.NewError(fcerror.ErrEmailAlreadyRegistered, nil)
		trans.Rollback()
		return
	}

	var err error
	user.Password, err = utils.HashScrypt(user.Password)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrPasswordHashingFailed, err)
		logrus.WithError(fcerr).Error("Failed to hash new user password")
		return
	}
	user.IsAdmin = false

	fcerr = trans.SaveUser(user)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to save user")
		trans.Rollback()
		return
	}

	fcerr = trans.Commit()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to commit transaction")
		return
	}

	count, fcerr := mgr.CountUsers(authorization.NewSystem())
	if fcerr == nil && count == 1 {
		isAdmin := true
		_, fcerr = mgr.UpdateUser(authorization.NewSystem(), user.ID, &models.UserUpdate{IsAdmin: &isAdmin})
		if fcerr == nil {
			user.IsAdmin = true
		} else {
			logrus.WithError(fcerr).Error("Failed to set first user an admin - ignore for now")
			fcerr = nil
		}
	} else if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to count users to determine whether created user should be an admin - ignore for now")
		fcerr = nil
	}

	return
}

func (mgr *userManager) GetUserByID(authCtx *authorization.Context, userID models.UserID) (user *models.User, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceSelf(authCtx, userID)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.userPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer trans.Close()

	user, fcerr = trans.GetUserByID(userID)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get user")
		return
	}
	if authCtx.Type != authorization.ContextTypeSystem {
		user.Password = ""
	}

	return
}

func (mgr *userManager) UpdateUser(authCtx *authorization.Context, userID models.UserID, updateUser *models.UserUpdate) (user *models.User, fcerr *fcerror.Error) {
	// TODO: Input Validation

	fcerr = authorization.EnforceSelf(authCtx, userID)
	if fcerr != nil {
		return
	}

	user, fcerr = mgr.GetUserByID(authCtx, userID)
	if fcerr != nil {
		return
	}

	if updateUser.FirstName != nil {
		user.FirstName = *updateUser.FirstName
	}
	if updateUser.LastName != nil {
		user.LastName = *updateUser.LastName
	}
	if updateUser.Email != nil {
		user.Email = *updateUser.Email
	}
	if updateUser.Password != nil {
		hashedPassword, err := utils.HashScrypt(*updateUser.Password)
		if err != nil {
			fcerr = fcerror.NewError(fcerror.ErrPasswordHashingFailed, err)
			logrus.WithError(err).WithField("userID", userID).Error("Failed to hash password for UpdateUser")
			return
		}
		user.Password = hashedPassword
	}
	if updateUser.IsAdmin != nil && authorization.EnforceAdmin(authCtx) == nil {
		user.IsAdmin = *updateUser.IsAdmin
	}

	trans, fcerr := mgr.userPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer func() { fcerr = trans.Finish(fcerr) }()

	fcerr = trans.UpdateUser(user)
	if fcerr != nil {
		logrus.WithError(fcerr).WithFields(logrus.Fields{"userID": userID, "updateUser": updateUser}).Error("Failed to update user")
		return
	}

	return
}

func (mgr *userManager) CountUsers(authCtx *authorization.Context) (count int64, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceAdmin(authCtx)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.userPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
	}
	defer trans.Close()

	count, fcerr = trans.CountUsers()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to count users")
	}
	return
}
