package parser

import (
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/reusing-code/dochan/pdf"
)

type ParserCallback func(filename string, strings []string)

func getFileNames(dir string) ([]string, error) {

	fileList := []string{}
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if filepath.Ext(path) == ".pdf" {
					//path = strings.TrimPrefix(path, dir)
					fileList = append(fileList, path)
				}
			}
			return nil
		},
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func ParseDir(dir string, cb ParserCallback) error {
	files, err := getFileNames(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		doc, err := pdf.ParsePDF(file)
		if err != nil {
			continue
		}
		cb(file, doc.GetText())
	}

	return nil
}
