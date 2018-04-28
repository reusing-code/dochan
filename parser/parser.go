package parser

import (
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
)

func ParseDir(dir string) ([]string, error) {

	fileList := []string{}
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if filepath.Ext(path) == ".pdf" {
					path = strings.TrimPrefix(path, dir)
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
