package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/reusing-code/dochan/parser"

	"github.com/reusing-code/dochan/searchTree"

	"github.com/gorilla/mux"
)

type server struct {
	port   int
	dir    string
	search *searchTree.SearchTree
}

type SearchResult struct {
	Count int
	Time  string
	Res   []string
}

func main() {
	serv := &server{}
	flag.IntVar(&serv.port, "port", 8092, "Listening port")
	flag.StringVar(&serv.dir, "path", "", "Document storage path")
	flag.Parse()

	err := serv.init()
	if err != nil {
		log.Fatal(err)
	}

	err = serv.start()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *server) init() error {
	s.search = searchTree.MakeSearchTree()
	fileCount := 0
	err := parser.ParseDir(s.dir, func(f parser.File, strings []string) {
		fileCount++
		s.search.AddContent(strings, f.Filename)
	}, parser.ExtensionFilter([]string{"pdf"}))
	if err != nil {
		return err
	}
	log.Printf("Parsed %v files", fileCount)
	return err
}

func (s *server) start() error {
	router := mux.NewRouter()

	router.HandleFunc("/search", s.searchHandler)

	http.Handle("/", router)

	log.Print((http.ListenAndServe(":"+strconv.Itoa(s.port), nil)))
	return nil
}

func (s *server) searchHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["query"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing search query"))
		return
	}

	searchKey := keys[0]

	log.Printf("Searching for '%v'", searchKey)
	start := time.Now()
	res := s.search.Search(searchKey, true)
	elapsed := time.Since(start)

	result := SearchResult{Count: len(res.GetRes()), Time: elapsed.String(), Res: res.GetResSlice()}
	log.Print(result)
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
