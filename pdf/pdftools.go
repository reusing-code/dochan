package pdf

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Document struct {
	pages []page
}

type page struct {
	sizeX  int32
	sizeY  int32
	blocks []textBlock
}

type textBlock struct {
	posX  int32
	posY  int32
	sizeX int32
	sizeY int32
	text  string
}

var tempDir string = "temp/"

func ParsePDF(path string) (doc *Document, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in file %v: %v\n", path, r)
			doc = nil
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	os.MkdirAll(tempDir, 0777)
	base := filepath.Base(path)
	tmpFile := filepath.Join(tempDir, base+"temp.xml")
	defer os.Remove(tmpFile)

	cmd := exec.Command("pdftohtml", "-xml", path, tmpFile)
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	doc, err = ParseFile(tmpFile)
	if err != nil {
		return nil, err
	}
	return
}

func (d *Document) GetText() []string {
	result := make([]string, 0)
	for _, page := range d.pages {
		for _, tb := range page.blocks {
			result = append(result, tb.text)
		}
	}
	return result
}
