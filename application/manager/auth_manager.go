package manager

import (
	"time"

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
	Login(email, password string) (*models.Session, *fcerror.Error)
	Logout(token models.Token) *fcerror.Error
	CreateUser(authCtx *authorization.Context, user *models.User) (*models.Session, *fcerror.Error)
	GetUserByID(authCtx *authorization.Context, userID models.UserID) (*models.User, *fcerror.Error)
	VerifyToken(token models.Token) (*models.User, *fcerror.Error)
}

func NewAuthManager(cfg config.Config) AuthManager {
	return &authManager{
		cfg:             cfg,
		userPersistence: persistence.GetUserPersistenceController(cfg),
		authPersistence: persistence.GetAuthPersistenceController(cfg),
		done:            make(chan struct{}),
	}
}

type authManager struct {
	cfg             config.Config
	userPersistence persistence.UserPersistenceController
	authPersistence persistence.AuthPersistenceController
	done            chan struct{}
}

func (mgr *authManager) Close() {
	mgr.done <- struct{}{}
}

func (mgr *authManager) cleanupExpiredSessionsRoutine() {
	interval := mgr.cfg.GetSessionCleanupInterval()
	logrus.WithField("interval", interval).Trace("Starting session cleanup")

	mgr.cleanupExpiredSessions()
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-mgr.done:
			return
		case <-ticker.C:
			mgr.cleanupExpiredSessions()
		}
	}
}

func (mgr *authManager) cleanupExpiredSessions() {
	logrus.Trace("Cleaning expired sessions")

	trans, fcerr := mgr.authPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}

	fcerr = trans.DeleteExpiredSessions()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to delete expired sessions")
		trans.Rollback()
		return
	}

	fcerr = trans.Commit()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to commit transaction")
		return
	}
}

func (mgr *authManager) Login(email, password string) (token *models.Session, fcerr *fcerror.Error) {
	trans, fcerr := mgr.userPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer trans.Close()

	user, fcerr := trans.GetUserByEmail(email)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("email", email).Warn("User with given email not found for login")
		fcerr = fcerror.NewError(fcerror.ErrUnauthorized, nil)
		return
	}

	valid, err := utils.ValidateScryptPassword(password, user.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to validate password")
		fcerr = fcerror.NewError(fcerror.ErrUnauthorized, nil)
		return
	}

	if !valid {
		logrus.WithError(fcerr).WithField("email", email).Warn("Unsuccessful login attempt for user")
		fcerr = fcerror.NewError(fcerror.ErrUnauthorized, nil)
		return
	}

	return mgr.createNewSession(user.ID)
}

func (mgr *authManager) Logout(token models.Token) (fcerr *fcerror.Error) {
	trans, fcerr := mgr.authPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}

	fcerr = trans.DeleteSessionByToken(token)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to delete token")
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

func (mgr *authManager) CreateUser(authCtx *authorization.Context, user *models.User) (token *models.Session, fcerr *fcerror.Error) {
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

	// Make first user an admin
	if user.ID == 0 {
		// TODO
	}

	fcerr = trans.Commit()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to commit transaction")
		return
	}
	return mgr.createNewSession(user.ID)
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
	defer trans.Close()

	user, fcerr = trans.GetUserByID(userID)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get user")
		return
	}
	user.Password = ""

	return
}

func (mgr *authManager) VerifyToken(token models.Token) (user *models.User, fcerr *fcerror.Error) {
	authTrans, fcerr := mgr.authPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer authTrans.Close()

	session, fcerr := authTrans.GetSessionByToken(token)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Token not found or failed to verify")
		return
	}

	if time.Now().After(session.ValidUntil) {
		fcerr = fcerror.NewError(fcerror.ErrSessionExpired, nil)
		return
	}

	return mgr.GetUserByID(authorization.NewSystem(), session.UserID)
}

func (mgr *authManager) createNewSession(userID models.UserID) (session *models.Session, fcerr *fcerror.Error) {
	session = &models.Session{
		Token:      models.Token(utils.GenerateRandomString(mgr.cfg.GetSessionTokenLength())),
		ValidUntil: utils.GetTimeIn(mgr.cfg.GetSessionExpirationDuration()),
		UserID:     userID,
	}

	trans, fcerr := mgr.authPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	fcerr = trans.SaveSession(session)
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
	return session, nil
}
