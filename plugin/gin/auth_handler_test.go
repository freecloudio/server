package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func TestLoginEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		willCallLogin bool
		success       bool
		expectedCode  int
	}{
		{name: "Correct login", input: "{\"username\": \"correct\", \"password\": \"correct\"}", success: true, willCallLogin: true, expectedCode: http.StatusOK},
		{name: "Incorrect login", input: "{\"username\": \"wrong\", \"password\": \"wrong\"}", success: false, willCallLogin: true, expectedCode: http.StatusUnauthorized},
		{name: "Incorrect input", input: "some stuff", success: false, willCallLogin: false, expectedCode: http.StatusBadRequest},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			if test.success && test.willCallLogin {
				mockAuthMgr.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&models.Session{}, nil).Times(1)
			} else if test.willCallLogin {
				mockAuthMgr.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, fcerror.NewError(fcerror.ErrUnauthorized, nil)).Times(1)
			}

			router := NewRouter(mockAuthMgr, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			resp, err := http.Post(fmt.Sprintf("%s/api/auth/login", testSrv.URL), "application/json", strings.NewReader(test.input))

			require.Nil(t, err, "Error calling login endpoint")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Login endpoint does not return expected status")
		})
	}
}

func TestLogoutEndpoint(t *testing.T) {
	var (
		good models.Token = "good"
		bad  models.Token = "bad"
	)

	tests := []struct {
		name         string
		input        string
		valid        bool
		success      bool
		expectedCode int
	}{
		{name: "Good Token", input: "Bearer " + string(good), valid: true, success: true, expectedCode: http.StatusNoContent},
		{name: "Good Token + fail logout", input: "Bearer " + string(good), valid: true, success: false, expectedCode: http.StatusInternalServerError},
		{name: "Bad Token", input: "Bearer " + string(bad), valid: false, expectedCode: http.StatusUnauthorized},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			if test.valid {
				mockAuthMgr.EXPECT().VerifyToken(good).Return(&models.User{}, nil).Times(1)
				if test.success {
					mockAuthMgr.EXPECT().Logout(good).Return(nil).Times(1)
				} else {
					mockAuthMgr.EXPECT().Logout(good).Return(fcerror.NewError(fcerror.ErrUnknown, nil)).Times(1)
				}
			} else {
				mockAuthMgr.EXPECT().VerifyToken(bad).Return(nil, fcerror.NewError(fcerror.ErrUnauthorized, nil)).Times(1)
			}

			router := NewRouter(mockAuthMgr, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/auth/logout", testSrv.URL), nil)
			require.Nil(t, err, "Failed creating logout request")
			req.Header.Add("Authorization", test.input)

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling logout endpoint")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Logout endpoint does not return expected status")
		})
	}
}
