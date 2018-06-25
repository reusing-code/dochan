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

	fileList, err := getFiles(tempDir, NoSkip)
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
			if parsedFile.Filename == filepath.Join(tempDir, tc.filePath) {
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

	ParseDir("testdata", func(f File, data []string) {
		filename := filepath.Base(f.Filename)
		results[filename] = data
	}, func(f File) bool { return false })

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

func TestSkipFiles(t *testing.T) {
	ParseDir("testdata", func(f File, data []string) {
		t.Errorf("No file should be parsed, but got callback for %q", f.Filename)
	}, func(f File) bool {
		if f.Hash != "28bac19a4147fdf7225f6c514270aa0867a2a03e" &&
			f.Hash != "8ead9513ba1f8253a30a709b87ae4a7fb386d0d8" &&
			f.Hash != "9f8b3703e131db0429d4497314c63becc7d4b0ec" {
			t.Errorf("File %q has unknown file hash %q", f.Filename, f.Hash)
		}
		return true
	})

	cbCount := 0
	ParseDir("testdata", func(f File, data []string) {
		cbCount++
		if f.Hash != "28bac19a4147fdf7225f6c514270aa0867a2a03e" &&
			f.Hash != "8ead9513ba1f8253a30a709b87ae4a7fb386d0d8" &&
			f.Hash != "9f8b3703e131db0429d4497314c63becc7d4b0ec" {
			t.Errorf("File %q has unknown file hash %q", f.Filename, f.Hash)
		}
	}, func(f File) bool {
		return false
	})
	if cbCount != 3 {
		t.Errorf("Wrong number of callbacks received: Want %d got %d", 3, cbCount)
	}
}

func TestFileCount(t *testing.T) {
	count, err := GetFileCount("testdata")
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("Wrong number of files received: Want %d got %d", 3, count)
	}
}
