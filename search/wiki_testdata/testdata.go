package wiki_testdata

import (
	"os"
	"runtime"

	wikiparse "github.com/dustin/go-wikiparse"
)

type StringCallback func(data string)

func ParseDataEN(callback StringCallback) error {
	envData := os.Getenv("WIKI_EN_DATAFILE")
	envIndex := os.Getenv("WIKI_EN_INDEXFILE")
	if envData != "" && envIndex != "" {
		return parseData(callback, envIndex, envData)
	}
	return parseData(callback, "/data/enwiki-20180420-pages-articles-multistream-index.txt.bz2",
		"/data/enwiki-20180420-pages-articles-multistream.xml.bz2")
}

func ParseDataDE(callback StringCallback) error {
	envData := os.Getenv("WIKI_DE_DATAFILE")
	envIndex := os.Getenv("WIKI_DE_INDEXFILE")
	if envData != "" && envIndex != "" {
		return parseData(callback, envIndex, envData)
	}
	return parseData(callback, "/data/dewiki-20180420-pages-articles-multistream-index.txt.bz2",
		"/data/dewiki-20180420-pages-articles-multistream.xml.bz2")
}

func parseData(callback StringCallback, indexFile string, dataFile string) error {
	p, err := wikiparse.NewIndexedParser(indexFile, dataFile, runtime.NumCPU())
	if err != nil {
		return err
	}

	count := 0
	for err == nil {
		var page *wikiparse.Page
		page, err = p.Next()
		if err == nil {
			callback(page.Revisions[0].Text)
			callback(page.Title)
		}
		if count > 5000 {
			break
		}
		count++
	}
	return nil
}
