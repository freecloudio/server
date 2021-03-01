package viperplg_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/plugin/viperplg"

	"github.com/stretchr/testify/assert"
)

func TestSetCorrectArgsAndRead(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	sessionTokenLength := 4444
	sessionExpiration := 1234

	newArgs := os.Args[1:2]
	newArgs = append(newArgs, fmt.Sprintf("--auth.session.token.length=%v", sessionTokenLength))
	newArgs = append(newArgs, fmt.Sprintf("--auth.session.expiration=%v", sessionExpiration))
	os.Args = newArgs

	cfg := viperplg.InitViperConfig()

	assert.Equal(t, sessionTokenLength, cfg.GetSessionTokenLength(), "Expect given token length to match parsed one")
	assert.Equal(t, time.Duration(sessionExpiration)*time.Hour, cfg.GetSessionExpirationDuration(), "Expect given token expiration to match parsed one")
	assert.Equal(t, config.NeoPersistenceKey, cfg.GetUserPersistencePluginKey(), "Expect not set config to have default")
	assert.Equal(t, config.NeoPersistenceKey, cfg.GetAuthPersistencePluginKey(), "Expect not set config to have default")
	assert.Equal(t, time.Hour, cfg.GetSessionCleanupInterval(), "Expect not set config to have default")
}

func TestSetIncorrectArgs(t *testing.T) {
	if os.Getenv("CALL_CONFIG") == "1" {
		viperplg.InitViperConfig()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSetIncorrectArgs", "--auth.session.token.length=WRONG")
	cmd.Env = append(os.Environ(), "CALL_CONFIG=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.Equal(t, 2, e.ExitCode(), "Expect exit code for wrong input")
		return
	}
	t.Fatalf("Test ran with err %v, want exit status 2", err)
}
