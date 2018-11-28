package main

import (
	"bytes"
	"encoding/gob"
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
	Count int        `json:"count"`
	Time  string     `json:"time"`
	Res   []Document `json:"results"`
}

type Document struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
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
		cont := ""
		if len(strings) > 0 {
			cont = strings[0]
		}
		doc := Document{ID: fileCount, Filename: f.Filename, Content: cont}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(doc)
		if err != nil {
			log.Printf("Error encoding file %v", f.Filename)
			return
		}
		s.search.AddContent(strings, buf.String())
	}, parser.ExtensionFilter([]string{"pdf"}))
	if err != nil {
		return err
	}
	log.Printf("Parsed %v files", fileCount)
	return err
}

func (s *server) start() error {
	router := mux.NewRouter()

	router.HandleFunc("/api/documents", s.searchHandler)

	http.Handle("/", router)

	log.Print((http.ListenAndServe(":"+strconv.Itoa(s.port), nil)))
	return nil
}

func (s *server) searchHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["q"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing search query"))
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	searchKey := keys[0]

	log.Printf("Searching for '%v'", searchKey)
	start := time.Now()
	res := s.search.Search(searchKey, true)
	elapsed := time.Since(start)

	var docs []Document
	for _, str := range res.GetResSlice() {
		buf := bytes.NewBufferString(str)
		var doc Document
		dec := gob.NewDecoder(buf)
		err := dec.Decode(&doc)
		if err != nil {
			log.Printf("Error decoding value")
			continue
		}
		docs = append(docs, doc)
	}

	result := SearchResult{Count: len(res.GetRes()), Time: elapsed.String(), Res: docs}
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
