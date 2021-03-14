package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseScryptStub(t *testing.T) {
	type data struct {
		Name        string
		Password    string
		ShouldError bool
	}

	var testdata = []data{
		{"empty password", "", true},
		{"only preemble", "$s1$", true},
		{"wrong preemble", "$s3$16384$8$1$SaltString$Base64Password=", true},
		{"missing parts", "$s1$16384$8$1$", true},
		{"valid password", "$s1$16384$8$1$VmFsaWRTYWx0$VmFsaWRQYXNzd29yZA0K", false},
		{"N no number", "$s1$ABC$8$1$VmFsaWRTYWx0$VmFsaWRQYXNzd29yZA0K", true},
		{"r no number", "$s1$16384$c$1$VmFsaWRTYWx0$VmFsaWRQYXNzd29yZA0K", true},
		{"p no number", "$s1$16384$8$b$VmFsaWRTYWx0$VmFsaWRQYXNzd29yZA0K", true},
		{"salt no b64", "$s1$16384$8$1$VmF*aWRTYWx0$VmFsaWRQYXNzd29yZA0K", true},
		{"hash no b64", "$s1$16384$8$1$VmFsaWRTYWx0$VmFsaWRQ)XNzd29yZA0K", true},
	}

	for _, td := range testdata {
		_, _, _, _, _, err := parseScryptStub(td.Password)
		if err == nil && td.ShouldError {
			t.Errorf("Expected error for %s but got none", td.Name)
		} else if err != nil && !td.ShouldError {
			t.Errorf("Expected no error for %s but got: %v", td.Name, err)
		}
	}

}

func TestPasswordHashing(t *testing.T) {
	// verify that hashing the same password two times does not yield the same result
	password := "testpassword"
	h1, err := HashScrypt(password)
	assert.Nil(t, err, "Error while hashing password")

	h2, err := HashScrypt(password)
	assert.Nil(t, err, "Error while hashing password")

	assert.NotEqual(t, h1, h2, "Hashing the same password twice yielded the same result")

	// first, hash a password, then test it against itself
	hash, err := HashScrypt(password)
	require.Nil(t, err, "Error while hashing password")

	err = ValidateScryptPassword(password, hash)
	assert.Nil(t, err, "Error while validating password")

	// bad hash
	badHash := "badStuff$" + hash
	err = ValidateScryptPassword(password, badHash)
	assert.NotNil(t, err, "No error while validating with bad hash")

	// wrong password
	err = ValidateScryptPassword(password+"123", hash)
	assert.NotNil(t, err, "No error while validating with bad password")
}
