package gin

import (
	"net/http"
	"testing"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAuthContext(t *testing.T) {
	c := &gin.Context{}
	logger := logrus.New()

	authContext := getAuthContext(c, logger)
	assert.Equal(t, authorization.ContextTypeAnonymous, authContext.Type, "No auth context does return anonymous")

	c.Set(authContextKey, "not a auth context")
	authContext = getAuthContext(c, logger)
	assert.Equal(t, authorization.ContextTypeAnonymous, authContext.Type, "Wrong context type does return anonymous")

	c.Set(authContextKey, authorization.NewSystem())
	authContext = getAuthContext(c, logger)
	assert.Equal(t, authorization.ContextTypeSystem, authContext.Type, "Wrong context type")
}

func TestAuthMiddleware(t *testing.T) {
	var (
		good models.Token = "good"
		bad  models.Token = "bad"
	)

	tests := []struct {
		name             string
		input            string
		validFormat      bool
		valid            bool
		expectedAuthType authorization.ContextType
	}{
		{name: "Valid Token", input: "Bearer " + string(good), validFormat: true, valid: true, expectedAuthType: authorization.ContextTypeUser},
		{name: "Invalid Token", input: "Bearer " + string(bad), validFormat: true, expectedAuthType: authorization.ContextTypeAnonymous},
		{name: "Too short", input: "short", expectedAuthType: authorization.ContextTypeAnonymous},
		{name: "No Header", input: "", expectedAuthType: authorization.ContextTypeAnonymous},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			if test.validFormat && test.valid {
				mockAuthMgr.EXPECT().VerifyToken(good).Return(&models.User{}, nil).Times(1)
			} else if test.validFormat {
				mockAuthMgr.EXPECT().VerifyToken(bad).Return(nil, fcerror.NewError(fcerror.ErrUnknown, nil)).Times(1)
			}

			authMiddleware := getAuthMiddleware(mockAuthMgr)

			c := &gin.Context{}
			req, err := http.NewRequest(http.MethodGet, "", nil)
			require.Nil(t, err, "Failed to create request")
			req.Header.Add("Authorization", test.input)
			c.Request = req
			logger := logrus.New()

			authMiddleware(c)

			authContext := getAuthContext(c, logger)
			assert.Equal(t, test.expectedAuthType, authContext.Type, "Wrong context type")
			if authContext.Type == authorization.ContextTypeUser {
				tokenInt, ok := c.Get(authTokenKey)
				require.True(t, ok, "AuthTokenKey is not set")
				token, ok := tokenInt.(models.Token)
				require.True(t, ok, "AuthTokenKey is not token type")
				assert.Equal(t, good, token, "Token in context does not match")
			}
		})
	}
}
