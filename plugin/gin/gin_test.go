package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freecloudio/server/application/manager"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination ../../mock/manager.go -package mock github.com/freecloudio/server/application/manager AuthManager,UserManager,NodeManager

func TestNewRouter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	managers := &manager.Managers{}
	router := NewRouter(managers, ":8080")

	assert.NotNil(t, router.engine, "Router engine is nil")
	assert.NotNil(t, router.srv, "Router srv is nil")
	assert.Equal(t, managers, router.managers, "managers is not the inserted managers")
}

func TestHealthEndpoint(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	router := NewRouter(&manager.Managers{}, ":8080")

	testSrv := httptest.NewServer(router.engine)
	defer testSrv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/health", testSrv.URL))

	assert.Nil(t, err, "Error calling health endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint does not return OK status")
}
