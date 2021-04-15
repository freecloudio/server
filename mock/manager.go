// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/freecloudio/server/application/manager (interfaces: AuthManager,UserManager,NodeManager)

// Package mock is a generated GoMock package.
package mock

import (
	io "io"
	reflect "reflect"

	authorization "github.com/freecloudio/server/application/authorization"
	models "github.com/freecloudio/server/domain/models"
	fcerror "github.com/freecloudio/server/domain/models/fcerror"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthManager is a mock of AuthManager interface.
type MockAuthManager struct {
	ctrl     *gomock.Controller
	recorder *MockAuthManagerMockRecorder
}

// MockAuthManagerMockRecorder is the mock recorder for MockAuthManager.
type MockAuthManagerMockRecorder struct {
	mock *MockAuthManager
}

// NewMockAuthManager creates a new mock instance.
func NewMockAuthManager(ctrl *gomock.Controller) *MockAuthManager {
	mock := &MockAuthManager{ctrl: ctrl}
	mock.recorder = &MockAuthManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthManager) EXPECT() *MockAuthManagerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockAuthManager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockAuthManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAuthManager)(nil).Close))
}

// CreateNewSession mocks base method.
func (m *MockAuthManager) CreateNewSession(arg0 models.UserID) (*models.Session, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewSession", arg0)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// CreateNewSession indicates an expected call of CreateNewSession.
func (mr *MockAuthManagerMockRecorder) CreateNewSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewSession", reflect.TypeOf((*MockAuthManager)(nil).CreateNewSession), arg0)
}

// Login mocks base method.
func (m *MockAuthManager) Login(arg0, arg1 string) (*models.Session, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthManagerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthManager)(nil).Login), arg0, arg1)
}

// Logout mocks base method.
func (m *MockAuthManager) Logout(arg0 models.Token) *fcerror.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0)
	ret0, _ := ret[0].(*fcerror.Error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthManagerMockRecorder) Logout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthManager)(nil).Logout), arg0)
}

// VerifyToken mocks base method.
func (m *MockAuthManager) VerifyToken(arg0 models.Token) (*models.User, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockAuthManagerMockRecorder) VerifyToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockAuthManager)(nil).VerifyToken), arg0)
}

// MockUserManager is a mock of UserManager interface.
type MockUserManager struct {
	ctrl     *gomock.Controller
	recorder *MockUserManagerMockRecorder
}

// MockUserManagerMockRecorder is the mock recorder for MockUserManager.
type MockUserManagerMockRecorder struct {
	mock *MockUserManager
}

// NewMockUserManager creates a new mock instance.
func NewMockUserManager(ctrl *gomock.Controller) *MockUserManager {
	mock := &MockUserManager{ctrl: ctrl}
	mock.recorder = &MockUserManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserManager) EXPECT() *MockUserManagerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockUserManager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockUserManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockUserManager)(nil).Close))
}

// CountUsers mocks base method.
func (m *MockUserManager) CountUsers(arg0 *authorization.Context) (int64, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUsers", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// CountUsers indicates an expected call of CountUsers.
func (mr *MockUserManagerMockRecorder) CountUsers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUsers", reflect.TypeOf((*MockUserManager)(nil).CountUsers), arg0)
}

// CreateUser mocks base method.
func (m *MockUserManager) CreateUser(arg0 *authorization.Context, arg1 *models.User) (*models.Session, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserManagerMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserManager)(nil).CreateUser), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockUserManager) GetUserByEmail(arg0 *authorization.Context, arg1 string) (*models.User, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserManagerMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserManager)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockUserManager) GetUserByID(arg0 *authorization.Context, arg1 models.UserID) (*models.User, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserManagerMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserManager)(nil).GetUserByID), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockUserManager) UpdateUser(arg0 *authorization.Context, arg1 models.UserID, arg2 *models.UserUpdate) (*models.User, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserManagerMockRecorder) UpdateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserManager)(nil).UpdateUser), arg0, arg1, arg2)
}

// MockNodeManager is a mock of NodeManager interface.
type MockNodeManager struct {
	ctrl     *gomock.Controller
	recorder *MockNodeManagerMockRecorder
}

// MockNodeManagerMockRecorder is the mock recorder for MockNodeManager.
type MockNodeManagerMockRecorder struct {
	mock *MockNodeManager
}

// NewMockNodeManager creates a new mock instance.
func NewMockNodeManager(ctrl *gomock.Controller) *MockNodeManager {
	mock := &MockNodeManager{ctrl: ctrl}
	mock.recorder = &MockNodeManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeManager) EXPECT() *MockNodeManagerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockNodeManager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockNodeManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockNodeManager)(nil).Close))
}

// CreateNode mocks base method.
func (m *MockNodeManager) CreateNode(arg0 *authorization.Context, arg1 *models.Node) (bool, *models.Node, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNode", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*models.Node)
	ret2, _ := ret[2].(*fcerror.Error)
	return ret0, ret1, ret2
}

// CreateNode indicates an expected call of CreateNode.
func (mr *MockNodeManagerMockRecorder) CreateNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNode", reflect.TypeOf((*MockNodeManager)(nil).CreateNode), arg0, arg1)
}

// CreateUserRootFolder mocks base method.
func (m *MockNodeManager) CreateUserRootFolder(arg0 *authorization.Context, arg1 models.UserID) *fcerror.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserRootFolder", arg0, arg1)
	ret0, _ := ret[0].(*fcerror.Error)
	return ret0
}

// CreateUserRootFolder indicates an expected call of CreateUserRootFolder.
func (mr *MockNodeManagerMockRecorder) CreateUserRootFolder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserRootFolder", reflect.TypeOf((*MockNodeManager)(nil).CreateUserRootFolder), arg0, arg1)
}

// DownloadFile mocks base method.
func (m *MockNodeManager) DownloadFile(arg0 *authorization.Context, arg1 models.NodeID) (*models.Node, io.ReadCloser, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", arg0, arg1)
	ret0, _ := ret[0].(*models.Node)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(*fcerror.Error)
	return ret0, ret1, ret2
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockNodeManagerMockRecorder) DownloadFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockNodeManager)(nil).DownloadFile), arg0, arg1)
}

// GetNodeByID mocks base method.
func (m *MockNodeManager) GetNodeByID(arg0 *authorization.Context, arg1 models.NodeID) (*models.Node, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Node)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// GetNodeByID indicates an expected call of GetNodeByID.
func (mr *MockNodeManagerMockRecorder) GetNodeByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeByID", reflect.TypeOf((*MockNodeManager)(nil).GetNodeByID), arg0, arg1)
}

// GetNodeByPath mocks base method.
func (m *MockNodeManager) GetNodeByPath(arg0 *authorization.Context, arg1 string) (*models.Node, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeByPath", arg0, arg1)
	ret0, _ := ret[0].(*models.Node)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// GetNodeByPath indicates an expected call of GetNodeByPath.
func (mr *MockNodeManagerMockRecorder) GetNodeByPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeByPath", reflect.TypeOf((*MockNodeManager)(nil).GetNodeByPath), arg0, arg1)
}

// ListByID mocks base method.
func (m *MockNodeManager) ListByID(arg0 *authorization.Context, arg1 models.NodeID) ([]*models.Node, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByID", arg0, arg1)
	ret0, _ := ret[0].([]*models.Node)
	ret1, _ := ret[1].(*fcerror.Error)
	return ret0, ret1
}

// ListByID indicates an expected call of ListByID.
func (mr *MockNodeManagerMockRecorder) ListByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByID", reflect.TypeOf((*MockNodeManager)(nil).ListByID), arg0, arg1)
}

// UploadFile mocks base method.
func (m *MockNodeManager) UploadFile(arg0 *authorization.Context, arg1 *models.Node, arg2 string) (bool, *models.Node, *fcerror.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*models.Node)
	ret2, _ := ret[2].(*fcerror.Error)
	return ret0, ret1, ret2
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockNodeManagerMockRecorder) UploadFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockNodeManager)(nil).UploadFile), arg0, arg1, arg2)
}

// UploadFileByID mocks base method.
func (m *MockNodeManager) UploadFileByID(arg0 *authorization.Context, arg1 models.NodeID, arg2 string) *fcerror.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFileByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*fcerror.Error)
	return ret0
}

// UploadFileByID indicates an expected call of UploadFileByID.
func (mr *MockNodeManagerMockRecorder) UploadFileByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFileByID", reflect.TypeOf((*MockNodeManager)(nil).UploadFileByID), arg0, arg1, arg2)
}
