package parser

import (
	"os"
	"path/filepath"
	"testing"
)

var tempDir string = "temp/"

var parseTestCases = []struct {
	filePath       string
	shouldBeParsed bool
}{
	{filePath: "test1.pdf", shouldBeParsed: true},
	{filePath: "a/test1.pdf", shouldBeParsed: true},
	{filePath: "a/b/c/d/e/f/test1.pdf", shouldBeParsed: true},
	{filePath: "test1.pd", shouldBeParsed: false},
	{filePath: "test1", shouldBeParsed: false},
}

func createEmptyFile(baseDir string, path string) error {
	filename := filepath.Join(baseDir, path)
	os.MkdirAll(filepath.Dir(filename), 0777)
	_, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	return err
}

func TestParse(t *testing.T) {
	os.MkdirAll(tempDir, 0777)
	defer os.RemoveAll(tempDir)

	expectedCount := 0

	// create files
	for _, tc := range parseTestCases {
		if tc.shouldBeParsed {
			expectedCount++
		}
		err := createEmptyFile(tempDir, tc.filePath)
		if err != nil {
			t.Fatal(err)
		}
	}

	fileList, err := ParseDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(fileList) != expectedCount {
		t.Errorf("Wrong number of parsed files. Want %d, got %d", expectedCount, len(fileList))
	}

	// check files
	for _, tc := range parseTestCases {
		if !tc.shouldBeParsed {
			continue
		}
		found := false
		for _, parsedFile := range fileList {
			if parsedFile == tc.filePath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("File %q not found.", tc.filePath)
		}
	}

}
