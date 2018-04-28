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

func parsePDF(path string) (*Document, error) {
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
