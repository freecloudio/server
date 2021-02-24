package application

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/sirupsen/logrus"
)

// AuthManager contains all use cases related to authentication and user management
type AuthManager interface {
	CreateUser(authCtx *authorization.Context, user *models.User) (*models.Token, *fcerror.Error)
	GetUserByID(authCtx *authorization.Context, userID models.UserID) (*models.User, *fcerror.Error)
	VerifyToken(token models.TokenValue) (*models.User, *fcerror.Error)
}

func NewAuthManager(cfg config.Config) AuthManager {
	return &authManager{
		cfg:             cfg,
		userPersistence: persistence.GetUserPersistenceController(cfg),
		authPersistence: persistence.GetAuthPersistenceController(cfg),
	}
}

type authManager struct {
	cfg             config.Config
	userPersistence persistence.UserPersistenceController
	authPersistence persistence.AuthPersistenceController
}

// TODO: Session cleanup

func (mgr *authManager) CreateUser(authCtx *authorization.Context, user *models.User) (token *models.Token, fcerr *fcerror.Error) {
	// TODO: Input Validation

	trans, fcerr := mgr.userPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}

	existingUser, fcerr := trans.GetUserByEmail(user.Email)
	if fcerr != nil && fcerr.ID != fcerror.ErrUserNotFound {
		logrus.WithError(fcerr).Error("Could not verify if user with this email already exists")
		return
	} else if fcerr == nil && existingUser != nil {
		fcerr = fcerror.NewError(fcerror.ErrUserNotFound, nil)
		return
	}

	user.Password, fcerr = utils.HashScrypt(user.Password)
	if fcerr != nil {
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

	// Make first user an admin
	if user.ID == 0 {
		// TODO
	}

	fcerr = trans.Commit()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to commit transaction")
		return
	}
	return mgr.createNewToken(user.ID)
}

func (mgr *authManager) GetUserByID(authCtx *authorization.Context, userID models.UserID) (user *models.User, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceSelf(authCtx, userID)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.userPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	user, fcerr = trans.GetUserByID(userID)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get user")
		return
	}
	return
}

func (mgr *authManager) VerifyToken(token models.TokenValue) (user *models.User, fcerr *fcerror.Error) {
	// TODO
	return mgr.GetUserByID(authorization.NewSystem(), models.UserID(0))
}

func (mgr *authManager) createNewToken(userID models.UserID) (token *models.Token, fcerr *fcerror.Error) {
	token = &models.Token{
		Value:      models.TokenValue(utils.GenerateRandomString(mgr.cfg.GetTokenValueLength())),
		ValidUntil: utils.GetTimeIn(mgr.cfg.GetTokenExpirationDuration()),
		UserID:     userID,
	}

	trans, fcerr := mgr.authPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	fcerr = trans.SaveToken(token)
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
	return token, nil
}
