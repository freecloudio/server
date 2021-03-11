package utils_test

import (
	"testing"

	"github.com/freecloudio/server/utils"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomStringLength(t *testing.T) {
	var l = []int{1, 5, 10, 20}
	for _, v := range l {
		assert.Equal(t, v, len(utils.GenerateRandomString(v)), "Random string has wrong length")
	}
}

func TestGenerateRandomStringUnique(t *testing.T) {
	assert.NotEqual(t, utils.GenerateRandomString(10), utils.GenerateRandomString(10), "Two different random string are the same")
}
