package application

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/sirupsen/logrus"
)

// UserManager contains all use cases related to user management
type UserManager interface {
	CreateUser(authCtx *authorization.Context, user *models.User) *fcerror.Error
	GetUser(authCtx *authorization.Context, userID models.UserID) (*models.User, *fcerror.Error)
}

func NewUserManager(cfg config.Config) UserManager {
	return &userManager{
		cfg:             cfg,
		userPersistence: persistence.GetUserPersistenceController(cfg),
	}
}

type userManager struct {
	cfg             config.Config
	userPersistence persistence.UserPersistenceController
}

func (mgr *userManager) CreateUser(authCtx *authorization.Context, user *models.User) (fcerr *fcerror.Error) {
	fcerr = authorization.EnforceAdmin(authCtx)
	if fcerr != nil {
		return fcerr
	}

	trans, fcerr := mgr.userPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
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
	return
}

func (mgr *userManager) GetUser(authCtx *authorization.Context, userID models.UserID) (user *models.User, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceSelf(authCtx, userID)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.userPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	user, fcerr = trans.GetUser(userID)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get user")
		return
	}
	return
}
