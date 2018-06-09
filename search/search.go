package search

import (
	"fmt"

	"github.com/reusing-code/dochan/parser"
	"github.com/reusing-code/dochan/searchTree"
)

type Search struct {
	tree *searchTree.SearchTree
}

func NewDirectorySearch(dir string) (*Search, error) {
	res := &Search{}
	res.tree = searchTree.MakeSearchTree()
	counter := 0
	parser.ParseDir(dir, func(filename string, strings []string) {
		res.tree.AddContent(strings, filename)
		counter++
		if counter%10 == 0 {
			fmt.Printf("Parsed %d documents\n", counter)
		}
	})

	return res, nil
}

func (s *Search) Search(query string, prefix bool) []string {
	resMap := s.tree.Search(query, prefix)
	res := make([]string, 0)
	for k, _ := range resMap.GetRes() {
		res = append(res, k.(string))
	}
	return res
}
