package eml

import (
	"os"
	"testing"
)

func TestExtractAttachments(t *testing.T) {
	f, err := os.Open("testdata/document.eml")
	if err != nil {
		t.Fatal(err)
	}
	fileFound := false
	err = ExtractAttachments(f, func(filename string, content []byte, messageID string) error {
		if filename != "IncredibleDocument.pdf" {
			t.Errorf("Want filename %q, got %q", "IncredibleDocument.pdf", filename)
		}
		if len(content) < 10000 {
			t.Errorf("File %q not large enough: %v B", filename, len(content))
		}
		fileFound = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !fileFound {
		t.Errorf("Expected file not extracted")
	}
}

func TestExtractAttachmentsFromDirRec(t *testing.T) {
	fileFound := false
	err := ExtractAttachmentsFromDirRec("testdata", func(filename string, content []byte, messageID string) error {
		if filename != "IncredibleDocument.pdf" {
			t.Errorf("Want filename %q, got %q", "IncredibleDocument.pdf", filename)
		}
		if len(content) < 10000 {
			t.Errorf("File %q not large enough: %v B", filename, len(content))
		}
		fileFound = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !fileFound {
		t.Errorf("Expected file not extracted")
	}
}
