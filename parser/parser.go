package parser

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/reusing-code/dochan/pdf"
)

type ParserCallback func(f File, strings []string, rawData []byte)
type SkipCallback func(f File) bool

type File struct {
	Filename string
	Hash     string
}

func NoSkip(f File) bool {
	return false
}

/**
Return a function that filters out all files with extensions
not matching the allowed ones (no dots ('.') in allowedExts!)
*/
func ExtensionFilter(allowedExts []string) func(f File) bool {
	return func(f File) bool {
		for _, ext := range allowedExts {
			if strings.EqualFold(filepath.Ext(f.Filename), "."+ext) {
				return false
			}
		}
		return true
	}
}

func getFiles(dir string, skip SkipCallback) ([]File, error) {
	fileList := []File{}
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				hash, err := hashSum(path)
				if err != nil {
					// TODO log error
					return nil
				}
				f := File{path, hash}
				if !skip(f) {
					fileList = append(fileList, f)
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
	fileList, err := getFiles(dir, NoSkip)
	if err != nil {
		return 0, err
	}
	return len(fileList), nil
}

func concurrentParse(input chan File, cb ParserCallback, resultMtx *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range input {
		doc, err := pdf.ParsePDF(file.Filename)
		if err != nil {
			continue
		}
		b, err := ioutil.ReadFile(file.Filename)
		if err != nil {
			continue
		}
		resultMtx.Lock()
		cb(file, doc.GetText(), b)
		resultMtx.Unlock()
	}
}

func ParseDir(dir string, cb ParserCallback, skip SkipCallback) error {
	resultMtx := &sync.Mutex{}
	files, err := getFiles(dir, skip)
	if err != nil {
		return err
	}

	inputChan := make(chan File, 10)

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

func hashSum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}
