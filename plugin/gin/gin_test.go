package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freecloudio/server/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination ../../mock/auth_manager.go -package mock github.com/freecloudio/server/application/manager AuthManager,UserManager

func TestNewRouter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
	mockUserMgr := mock.NewMockUserManager(mockCtrl)

	router := NewRouter(mockAuthMgr, mockUserMgr, ":8080")

	assert.NotNil(t, router.engine, "Router engine is nil")
	assert.NotNil(t, router.srv, "Router srv is nil")
	assert.Equal(t, mockAuthMgr, router.authMgr, "Authmanager is not the inserted mock manager")
}

func TestHealthEndpoint(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthMgr := mock.NewMockAuthManager(mockCtrl)
	mockUserMgr := mock.NewMockUserManager(mockCtrl)

	router := NewRouter(mockAuthMgr, mockUserMgr, ":8080")

	testSrv := httptest.NewServer(router.engine)
	defer testSrv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/health", testSrv.URL))

	assert.Nil(t, err, "Error calling health endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint does not return OK status")
}
