package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNodeInfoByPath(t *testing.T) {
	var (
		token models.Token = "token"
	)

	tests := []struct {
		name         string
		input        string
		success      bool
		expectedCode int
	}{
		{name: "Empty Path Success", input: "", success: true, expectedCode: http.StatusOK},
		{name: "Slash Path Success", input: "/", success: true, expectedCode: http.StatusOK},
		{name: "Extended Path Success", input: "/folder/file.txt", success: true, expectedCode: http.StatusOK},
		{name: "Extended Path Failure", input: "/folder/file.txt", success: false, expectedCode: http.StatusNotFound},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			mockAuthMgr.EXPECT().VerifyToken(token).Return(&models.User{}, nil).Times(1)
			mockNodeMgr := mock.NewMockNodeManager(mockCtrl)
			param := test.input
			if param == "" {
				param = "/"
			}
			if test.success {
				mockNodeMgr.EXPECT().GetNodeByPath(gomock.Any(), param).Return(&models.Node{}, nil).Times(1)
			} else {
				mockNodeMgr.EXPECT().GetNodeByPath(gomock.Any(), param).Return(nil, fcerror.NewError(fcerror.ErrNodeNotFound, nil)).Times(1)
			}

			managers := &manager.Managers{Node: mockNodeMgr, Auth: mockAuthMgr}
			router := NewRouter(managers, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/node/info/path%s", testSrv.URL, test.input), nil)
			require.Nil(t, err, "Failed creating get node info by path request")
			req.Header.Add("Authorization", "Bearer "+string(token))

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling get node info by path endpoint")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Get node info by path endpoint does not return expected status")
		})
	}
}

func TestGetNodeInfoByID(t *testing.T) {
	var (
		token models.Token = "token"
	)

	tests := []struct {
		name         string
		input        models.NodeID
		success      bool
		expectedCode int
	}{
		{name: "Success", input: models.NodeID(1), success: true, expectedCode: http.StatusOK},
		{name: "Failure", input: models.NodeID(2), success: false, expectedCode: http.StatusNotFound},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			mockAuthMgr.EXPECT().VerifyToken(token).Return(&models.User{}, nil).Times(1)
			mockNodeMgr := mock.NewMockNodeManager(mockCtrl)
			if test.success {
				mockNodeMgr.EXPECT().GetNodeByID(gomock.Any(), test.input).Return(&models.Node{}, nil).Times(1)
			} else {
				mockNodeMgr.EXPECT().GetNodeByID(gomock.Any(), test.input).Return(nil, fcerror.NewError(fcerror.ErrNodeNotFound, nil)).Times(1)
			}

			managers := &manager.Managers{Node: mockNodeMgr, Auth: mockAuthMgr}
			router := NewRouter(managers, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/node/info/id/%d", testSrv.URL, test.input), nil)
			require.Nil(t, err, "Failed creating get node info by id request")
			req.Header.Add("Authorization", "Bearer "+string(token))

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling get node info by id endpoint")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Get node info by id endpoint does not return expected status")
		})
	}
}
