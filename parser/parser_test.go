package parser

import (
	"os"
	"path/filepath"
	"strings"
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

func TestGetFilenames(t *testing.T) {
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

	fileList, err := getFileNames(tempDir)
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
			if parsedFile == filepath.Join(tempDir, tc.filePath) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("File %q not found.", tc.filePath)
		}
	}

}

var parseDirTestCases = []struct {
	filename      string
	containedText string
}{
	{filename: "A.pdf", containedText: "TestStringA"},
	{filename: "B.pdf", containedText: "OtherStringB"},
	{filename: "C.pdf", containedText: "ThirdStringC"},
}

func TestParseDir(t *testing.T) {

	results := make(map[string][]string)

	ParseDir("testdata", func(filename string, data []string) {
		filename = filepath.Base(filename)
		results[filename] = data
	})

	for _, tc := range parseDirTestCases {
		resultData := results[tc.filename]
		if resultData == nil || len(resultData) == 0 {
			t.Errorf("File %q not parsed correctly. No data found", tc.filename)
			continue
		}

		if !strings.Contains(resultData[0], tc.containedText) {
			t.Errorf("Text %q not found in %q", tc.containedText, resultData[0])
		}
	}

}
