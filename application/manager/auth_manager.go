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
	VerifyToken(token models.Token) (*models.Session, *fcerror.Error)
	CreateNewSession(userID models.UserID) (*models.Session, *fcerror.Error)
	Close()
}

func NewAuthManager(cfg config.Config, authPersistence persistence.AuthPersistenceController, managers *Managers) AuthManager {
	authMgr := &authManager{
		cfg:             cfg,
		authPersistence: authPersistence,
		managers:        managers,
		done:            make(chan struct{}),
	}
	go authMgr.cleanupExpiredSessionsRoutine()

	managers.Auth = authMgr
	return authMgr
}

type authManager struct {
	cfg             config.Config
	authPersistence persistence.AuthPersistenceController
	managers        *Managers
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
	defer func() { _ = trans.Finish(fcerr) }()

	fcerr = trans.DeleteExpiredSessions()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to delete expired sessions")
		return
	}
}

func (mgr *authManager) Login(email, password string) (token *models.Session, fcerr *fcerror.Error) {
	user, fcerr := mgr.managers.User.GetUserByEmail(authorization.NewSystem(), email)
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

	return mgr.CreateNewSession(user.ID)
}

func (mgr *authManager) Logout(token models.Token) (fcerr *fcerror.Error) {
	trans, fcerr := mgr.authPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer func() { fcerr = trans.Finish(fcerr) }()

	fcerr = trans.DeleteSessionByToken(token)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to delete token")
		return
	}

	return
}

func (mgr *authManager) VerifyToken(token models.Token) (session *models.Session, fcerr *fcerror.Error) {
	authTrans, fcerr := mgr.authPersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer authTrans.Close()

	session, fcerr = authTrans.GetSessionByToken(token)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Token not found or failed to verify")
		return
	}

	if time.Now().After(session.ValidUntil) {
		fcerr = fcerror.NewError(fcerror.ErrSessionExpired, nil)
		return
	}

	return
}

func (mgr *authManager) CreateNewSession(userID models.UserID) (session *models.Session, fcerr *fcerror.Error) {
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
	defer func() { fcerr = trans.Finish(fcerr) }()

	fcerr = trans.SaveSession(session)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to save user")
		return
	}
	return session, nil
}
