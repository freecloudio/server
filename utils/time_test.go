package utils_test

import (
	"testing"
	"time"

	"github.com/freecloudio/server/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentTimeIsUTC(t *testing.T) {
	assert.Equal(t, time.UTC, utils.GetCurrentTime().Location(), "Current time is not in UTC")
}

func TestGetTimeInIsUTC(t *testing.T) {
	assert.Equal(t, time.UTC, utils.GetTimeIn(time.Hour).Location(), "TimeIn is not in UTC")
}
