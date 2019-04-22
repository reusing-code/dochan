package refuel

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

type FuelRecord struct {
	RefuelRecord
	DrivenKM int    `json:"drivenKM"`
	ID       uint64 `json:"id"`
}

type Handler struct {
	DataBase *DB
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Register(dbpath string, router *mux.Router) error {
	db, err := New(dbpath)
	if err != nil {
		return err
	}
	h := &Handler{db}

	router.HandleFunc("/submit", h.fuelSubmitHandler)
	router.HandleFunc("", h.fuelHandler)
	return nil
}

func (h *Handler) fuelHandler(w http.ResponseWriter, r *http.Request) {
	page := int(1)
	limit := int(math.MaxInt32)
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	if pageParam != "" && limitParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	records := make([]FuelRecord, 0)
	err := h.DataBase.GetAllFuelRecords(func(key uint64, record *RefuelRecord) {
		newRecord := FuelRecord{*record, 0, key}
		records = append(records, newRecord)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].TotalKM > records[j].TotalKM
	})
	for i := 0; i < len(records)-1; i++ {
		records[i].DrivenKM = records[i].TotalKM - records[i+1].TotalKM - records[i].IgnoreKM
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(len(records)))

	start := min((page-1)*limit, len(records))
	end := min(page*limit, len(records))
	pagedRecords := records[start:end]

	js, err := json.Marshal(pagedRecords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *Handler) fuelSubmitHandler(w http.ResponseWriter, r *http.Request) {
	var record RefuelRecord
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
	err = h.DataBase.AddFuelRecord(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
