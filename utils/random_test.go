package utils_test

import (
	"testing"

	"github.com/freecloudio/server/utils"
)

func TestGenerateRandomStringLength(t *testing.T) {
	var l = []int{1, 5, 10, 20}
	for _, v := range l {
		if length := len(utils.GenerateRandomString(v)); length != v {
			t.Errorf("Expected string of length %d, but got %d", v, length)
		}
	}
}

func TestGenerateRandomStringUnique(t *testing.T) {
	if utils.GenerateRandomString(10) == utils.GenerateRandomString(10) { //nolint:staticcheck
		t.Error("Expected two different random strings but got two times the same")
	}
}
