package pdf

import (
	"os"
	"os/exec"
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

func ParsePDF(path string) (*Document, error) {
	os.MkdirAll(tempDir, 0777)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("pdftohtml", "-xml", path, tempDir+"temp.xml")
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	doc, err := ParseFile(tempDir + "temp.xml")
	if err != nil {
		return nil, err
	}
	return doc, nil
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
