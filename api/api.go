package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func StartServer(port int) error {
	router := mux.NewRouter()

	http.Handle("/", router)
	go log.Print((http.ListenAndServe(":"+strconv.Itoa(port), nil)))
	return nil
}
