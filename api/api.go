package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/reusing-code/dochan/refuel"

	"github.com/reusing-code/dochan/db"
	"github.com/reusing-code/dochan/parser"

	"github.com/reusing-code/dochan/searchTree"

	"github.com/gorilla/mux"
)

type server struct {
	port      int
	dir       string
	search    *searchTree.SearchTree
	db        *db.DB
	assetPath string
}

type SearchResult struct {
	Count int        `json:"count"`
	Time  string     `json:"time"`
	Res   []Document `json:"results"`
}

type Document struct {
	ID       uint64 `json:"id"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type ResponseDocument struct {
	ID         uint64 `json:"id"`
	Filename   string `json:"filename"`
	RawContent []byte `json:"content"`
}

type FuelRecord struct {
	refuel.RefuelRecord
	ID uint64 `json:"id"`
}

func main() {
	serv := &server{}
	var dbPath string
	flag.IntVar(&serv.port, "port", 8092, "Listening port")
	flag.StringVar(&serv.dir, "path", "", "Document storage path")
	flag.StringVar(&dbPath, "dbFile", "dochan.db", "DB File storage")
	flag.StringVar(&serv.assetPath, "assetPath", "assets/", "Static assets to serve")
	flag.Parse()

	var err error
	serv.db, err = db.New(dbPath)

	if err != nil {
		log.Fatal(err)
	}

	err = serv.init()
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
	err := parser.ParseDir(s.dir, func(f parser.File, strings []string, rawData []byte) {
		s.db.AddFile(f.Filename, f.Hash, rawData, strings)
		fileCount++

	}, parser.ExtensionFilter([]string{"pdf"}, func(f parser.File) bool {
		return s.db.Contains(f.Hash)
	}))
	if err != nil {
		return err
	}
	log.Printf("Added %v new files", fileCount)

	err = s.db.GetAllFiles(func(key uint64, file db.DBFile) {
		cont := ""
		if len(file.Content) > 0 {
			cont = file.Content[0]
		}
		doc := Document{ID: key, Filename: file.Name, Content: cont}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(doc)
		if err != nil {
			log.Printf("Error encoding file %v", file.Path)
			return
		}
		s.search.AddContent(file.Content, buf.String())
	})
	if err != nil {
		return err
	}
	return err
}

func (s *server) start() error {
	router := mux.NewRouter()

	clientSideRoutes := []string{"/about", "/document", "/search", "/fuel"}
	router.HandleFunc("/api/documents", s.searchHandler)
	router.HandleFunc("/api/documents/{key:[0-9]+}", s.documentHandler)
	router.HandleFunc("/api/documents/{key:[0-9]+}/download", s.downloadHandler)
	router.HandleFunc("/api/fuel", s.fuelHandler)
	router.HandleFunc("/api/fuel/submit", s.fuelSubmitHandler)
	router.HandleFunc("/api/fuel/csv", s.fuelCSVHandler)

	for _, route := range clientSideRoutes {
		router.PathPrefix(route).HandlerFunc(s.indexHandler)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(s.assetPath)))

	http.Handle("/", router)

	log.Print((http.ListenAndServe(":"+strconv.Itoa(s.port), nil)))
	return nil
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(s.assetPath, "index.html"))
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

func (s *server) documentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	keyStr := mux.Vars(r)["key"]
	i, err := strconv.ParseInt(keyStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	key := uint64(i)
	f, err := s.db.GetFile(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	doc := &ResponseDocument{ID: key, Filename: f.Name, RawContent: f.RawData}
	js, err := json.Marshal(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *server) downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	keyStr := mux.Vars(r)["key"]
	i, err := strconv.ParseInt(keyStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	key := uint64(i)
	f, err := s.db.GetFile(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(f.RawData)
}

func (s *server) fuelHandler(w http.ResponseWriter, r *http.Request) {
	records := make([]FuelRecord, 0)
	err := refuel.GetAllFuelRecords(s.db, func(key uint64, record *refuel.RefuelRecord) {
		newRecord := FuelRecord{*record, key}
		records = append(records, newRecord)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *server) fuelSubmitHandler(w http.ResponseWriter, r *http.Request) {
	var record refuel.RefuelRecord
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(buf, &record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = refuel.AddFuelRecord(s.db, &record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) fuelCSVHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = refuel.ParseCSV(s.db, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
