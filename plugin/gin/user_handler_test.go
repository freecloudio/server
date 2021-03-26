package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserEndpoint(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		willCallRegister bool
		success          bool
		expectedCode     int
	}{
		{name: "Successful Register", input: "{\"username\": \"user\"}", success: true, willCallRegister: true, expectedCode: http.StatusCreated},
		{name: "Failing Register", input: "{\"username\": \"user\"}", success: false, willCallRegister: true, expectedCode: http.StatusInternalServerError},
		{name: "Incorrect input", input: "some stuff", success: false, willCallRegister: false, expectedCode: http.StatusBadRequest},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserMgr := mock.NewMockUserManager(mockCtrl)
			if test.success && test.willCallRegister {
				mockUserMgr.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(&models.Session{}, nil).Times(1)
			} else if test.willCallRegister {
				mockUserMgr.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, fcerror.NewError(fcerror.ErrUnknown, nil)).Times(1)
			}

			managers := &manager.Managers{User: mockUserMgr}
			router := NewRouter(managers, nil, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			resp, err := http.Post(fmt.Sprintf("%s/api/user", testSrv.URL), "application/json", strings.NewReader(test.input))

			require.Nil(t, err, "Error calling register endpoint")
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Register endpoint does not return expected status")
		})
	}
}

func TestGetOwnUserEndpoint(t *testing.T) {
	var (
		good models.Token = "good"
		bad  models.Token = "bad"
	)

	tests := []struct {
		name         string
		input        string
		success      bool
		expectedCode int
	}{
		{name: "Logged in user", input: "Bearer " + string(good), success: true, expectedCode: http.StatusOK},
		{name: "Not logged in user", input: "Bearer " + string(bad), success: false, expectedCode: http.StatusUnauthorized},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			if test.success {
				mockAuthMgr.EXPECT().VerifyToken(good).Return(&models.User{}, nil).Times(1)
			} else {
				mockAuthMgr.EXPECT().VerifyToken(bad).Return(nil, fcerror.NewError(fcerror.ErrUnauthorized, nil)).Times(1)
			}

			managers := &manager.Managers{Auth: mockAuthMgr}
			router := NewRouter(managers, nil, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/user", testSrv.URL), nil)
			require.Nil(t, err, "Failed creating get own user request")
			req.Header.Add("Authorization", test.input)

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling get own user endpoint")
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Get own user endpoint does not return expected status")
		})
	}
}

func TestGetUserByIDEndpoint(t *testing.T) {
	tests := []struct {
		name         string
		inputID      models.UserID
		success      bool
		expectedCode int
	}{
		{name: "valid user id", inputID: models.UserID("1"), success: true, expectedCode: http.StatusOK},
		{name: "failed to get", inputID: models.UserID("1"), success: false, expectedCode: http.StatusInternalServerError},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockUserMgr := mock.NewMockUserManager(mockCtrl)
			if test.success && test.inputID != models.UserID("0") {
				mockUserMgr.EXPECT().GetUserByID(gomock.Any(), test.inputID).Return(&models.User{}, nil).Times(1)
			} else if test.inputID != models.UserID("0") {
				mockUserMgr.EXPECT().GetUserByID(gomock.Any(), test.inputID).Return(nil, fcerror.NewError(fcerror.ErrUnknown, nil)).Times(1)
			}

			managers := &manager.Managers{User: mockUserMgr}
			router := NewRouter(managers, nil, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/user/%s", testSrv.URL, test.inputID), nil)
			require.Nil(t, err, "Failed creating get user by id request")
			req.Header.Add("Authorization", string(test.inputID))

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling get user by id endpoint")
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Get user by id endpoint does not return expected status")
		})
	}
}
