package parser

import (
	"path/filepath"
	"runtime"
	"sync"

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

func GetFileCount(dir string) (int, error) {
	fileList, err := getFileNames(dir)
	if err != nil {
		return 0, err
	}
	return len(fileList), nil
}

func concurrentParse(input chan string, cb ParserCallback, resultMtx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range input {
		doc, err := pdf.ParsePDF(file)
		if err != nil {
			continue
		}
		resultMtx.Lock()
		cb(file, doc.GetText())
		resultMtx.Unlock()
	}
}

func ParseDir(dir string, cb ParserCallback) error {
	resultMtx := &sync.Mutex{}
	files, err := getFileNames(dir)
	if err != nil {
		return err
	}

	inputChan := make(chan string, 10)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go concurrentParse(inputChan, cb, resultMtx, &wg)
	}

	for _, file := range files {
		inputChan <- file
	}
	close(inputChan)
	wg.Wait()

	return nil
}
