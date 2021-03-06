package utils_test

import (
	"testing"

	"github.com/freecloudio/server/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetPathSegments(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedSegments []string
	}{
		{"Empty path", "", []string{}},
		{"Empty path with spaces", "  \t  ", []string{}},
		{"Only Slash", "/", []string{}},
		{"One segment", "folder", []string{"folder"}},
		{"One segment preceding slash", "/folder", []string{"folder"}},
		{"One segment following slash", "folder/", []string{"folder"}},
		{"One segment slashes preceding and following", "/folder/", []string{"folder"}},
		{"Multiple segments preceding slash", "/folder/file.txt", []string{"folder", "file.txt"}},
		{"Multiple segments following slash", "folder/second_folder/", []string{"folder", "second_folder"}},
		{"Multiple segments slashes preceding and following", "/folder/second_folder/", []string{"folder", "second_folder"}},
		{"Multiple segments slashes preceding and following with spaces", "  /folder /second_folder /", []string{"folder", "second_folder"}},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			actual := utils.GetPathSegments(test.input)
			assert.Equal(t, test.expectedSegments, actual)
		})
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedPath     string
		expectedFilename string
	}{
		{"Empty", "", "", ""},
		{"Only Slash", "/", "/", ""},
		{"Only Filename", "file.txt", "", "file.txt"},
		{"Slash and filename", "/file.txt", "/", "file.txt"},
		{"Folder and filename", "folder/file.txt", "folder/", "file.txt"},
		{"Folder with preceding slash and filename", "/folder/file.txt", "/folder/", "file.txt"},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			actualPath, actualFilename := utils.SplitPath(test.input)
			assert.Equal(t, test.expectedPath, actualPath)
			assert.Equal(t, test.expectedFilename, actualFilename)
		})
	}
}
