package utils

import (
	"testing"
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
	h1, err := HashScrypt("testpassword")
	if err != nil {
		t.Errorf("Error while hashing password: %v", err)
		return
	}
	h2, err := HashScrypt("testpassword")
	if err != nil {
		t.Errorf("Error while hashing password: %v", err)
		return
	}
	if h1 == h2 {
		t.Errorf("Hashing the same password twice yielded the same result: %s", h1)
		return
	}
	// first, hash a password, then test it against itself
	hash, err := HashScrypt("h4x0r!")
	if err != nil {
		t.Errorf("Got error while hashing password: %v", err)
		return
	}
	valid, err := ValidateScryptPassword("h4x0r!", hash)
	if err != nil {
		t.Errorf("Got error while validating password: %v", err)
		return
	}
	if !valid {
		t.Errorf("Expected a valid password verification, but it failed")
	}
}
