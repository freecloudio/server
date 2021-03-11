package gin

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authContextKey = "authentication_context"
	authTokenKey   = "authentication_token"

	authHeaderName = "Authorization"
	authPrefix     = "Bearer "
)

func getAuthContext(c *gin.Context) *authorization.Context {
	authContextInt, found := c.Get(authContextKey)
	if !found {
		logrus.Warn("AuthContext not found in gin context")
		return authorization.NewAnonymous()
	}
	authContext, ok := authContextInt.(*authorization.Context)
	if !ok {
		logrus.Warn("AuthContext in gin context is not of correct type")
		return authorization.NewAnonymous()
	}
	return authContext
}

func getAuthMiddleware(authMgr manager.AuthManager, userMgr manager.UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authContext *authorization.Context

		authHeader := c.GetHeader(authHeaderName)
		if len(authHeader) <= len(authPrefix) {
			authContext = authorization.NewAnonymous()
		} else {
			tokenString := models.Token(authHeader[len(authPrefix):])
			session, fcerr := authMgr.VerifyToken(tokenString)
			var user *models.User
			if fcerr == nil {
				user, fcerr = userMgr.GetUserByID(authorization.NewUser(&models.User{ID: session.UserID}), session.UserID)
			}
			if fcerr == nil {
				authContext = authorization.NewUser(user)
				c.Set(authTokenKey, tokenString)
			} else {
				authContext = authorization.NewAnonymous()
			}
		}

		c.Set(authContextKey, authContext)
		c.Next()
	}
}
