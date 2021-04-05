package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"
	"github.com/freecloudio/server/utils"
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
			router := NewRouter(managers, nil, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/node/info/path%s", testSrv.URL, test.input), nil)
			require.Nil(t, err, "Failed creating get node info by path request")
			req.Header.Add("Authorization", "Bearer "+string(token))

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling get node info by path endpoint")
			defer resp.Body.Close()
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
		{name: "Success", input: models.NodeID("1"), success: true, expectedCode: http.StatusOK},
		{name: "Failure", input: models.NodeID("2"), success: false, expectedCode: http.StatusNotFound},
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
			router := NewRouter(managers, nil, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/node/info/id/%s", testSrv.URL, test.input), nil)
			require.Nil(t, err, "Failed creating get node info by id request")
			req.Header.Add("Authorization", "Bearer "+string(token))

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling get node info by id endpoint")
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Get node info by id endpoint does not return expected status")
		})
	}
}

func TestCreateNodeByID(t *testing.T) {
	var (
		token models.Token = "token"
	)

	tests := []struct {
		name              string
		inputParentNodeID models.NodeID
		inputFile         bool
		inputNew          bool
		success           bool
		expectedCode      int
	}{
		{name: "Success Existing Folder", inputParentNodeID: models.NodeID("1"), success: true, inputFile: false, inputNew: false, expectedCode: http.StatusOK},
		{name: "Failure Existing Folder", inputParentNodeID: models.NodeID("2"), success: false, inputFile: false, inputNew: false, expectedCode: http.StatusNotFound},
		{name: "Success New File", inputParentNodeID: models.NodeID("2"), success: true, inputFile: true, inputNew: true, expectedCode: http.StatusOK},
		{name: "Failure New File", inputParentNodeID: models.NodeID("2"), success: false, inputFile: true, inputNew: true, expectedCode: http.StatusNotFound},
		{name: "Success Existing File", inputParentNodeID: models.NodeID("2"), success: true, inputFile: true, inputNew: false, expectedCode: http.StatusOK},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
			mockAuthMgr.EXPECT().VerifyToken(token).Return(&models.User{}, nil).Times(1)
			mockNodeMgr := mock.NewMockNodeManager(mockCtrl)
			nodeType := models.NodeTypeFolder
			if test.inputFile {
				nodeType = models.NodeTypeFile
			}
			node := &models.Node{ID: models.NodeID("1"), ParentNodeID: &test.inputParentNodeID, Name: utils.GenerateRandomString(10), Type: nodeType}
			if test.success {
				mockNodeMgr.EXPECT().CreateNode(gomock.Any(), node).Return(test.inputNew, node, nil).Times(1)
			} else {
				mockNodeMgr.EXPECT().CreateNode(gomock.Any(), node).Return(test.inputNew, nil, fcerror.NewError(fcerror.ErrNodeNotFound, nil)).Times(1)
			}

			managers := &manager.Managers{Node: mockNodeMgr, Auth: mockAuthMgr}
			router := NewRouter(managers, nil, ":8080")

			testSrv := httptest.NewServer(router.engine)
			defer testSrv.Close()

			jsonNode, err := json.Marshal(node)
			require.Nil(t, err, "Failed creating json from node")
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/node/create/", testSrv.URL), bytes.NewReader(jsonNode))
			require.Nil(t, err, "Failed creating create node by parent id request")
			req.Header.Add("Authorization", "Bearer "+string(token))

			resp, err := http.DefaultClient.Do(req)

			require.Nil(t, err, "Error calling create node by parent id endpoint")
			defer resp.Body.Close()
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Create node by parent id endpoint does not return expected status")

			resJSON := map[string]interface{}{}
			err = json.NewDecoder(resp.Body).Decode(&resJSON)
			require.Nil(t, err, "Failed to decode response JSON")
			if test.success {
				assert.Equal(t, string(node.ID), resJSON["node_id"], "Returned node id does not match")
				assert.Equal(t, test.inputNew, resJSON["created"], "Returned created flag does not match")
			}
		})
	}
}
