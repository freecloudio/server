package utils_test

import (
	"testing"
	"time"

	"github.com/freecloudio/server/utils"
)

func TestGetCurrentTimeIsUTC(t *testing.T) {
	if utils.GetCurrentTime().Location() != time.UTC {
		t.Error("Current time is not in UTC")
	}
}

func TestGetTimeInIsUTC(t *testing.T) {
	if utils.GetTimeIn(time.Hour).Location() != time.UTC {
		t.Error("Current time is not in UTC")
	}
}
