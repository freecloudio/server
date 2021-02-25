package gin

import (
	"github.com/freecloudio/server/application"
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/domain/models"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const (
	authContextKey = "authentication_context"

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

func getAuthMiddleware(authMgr application.AuthManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authContext *authorization.Context

		authHeader := c.GetHeader(authHeaderName)
		if len(authHeader) <= len(authPrefix) {
			authContext = authorization.NewAnonymous()
		} else {
			tokenString := authHeader[len(authPrefix):]
			user, fcerr := authMgr.VerifyToken(models.TokenValue(tokenString))
			if fcerr == nil {
				authContext = authorization.NewUser(user)
			} else {
				authContext = authorization.NewAnonymous()
			}
		}

		c.Set(authContextKey, authContext)
		c.Next()
	}
}
