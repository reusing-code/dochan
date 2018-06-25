package db

import (
	"os"
	"testing"
)

func TestOpenCloseDB(t *testing.T) {
	defer os.Remove("test.db")
	db, err := New("test.db")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	db2 := &DB{}
	err = db2.Close()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestSetGetHashTable(t *testing.T) {
	defer os.Remove("test.db")
	db, err := New("test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	table, err := db.GetHashTable()
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 0 {
		t.Error("Hash table not empty")
	}

	table = make(map[string]bool)
	table["test"] = true
	table["bla"] = true

	err = db.SetHashTable(table)
	if err != nil {
		t.Fatal(err)
	}

	table, err = db.GetHashTable()
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 2 {
		t.Error("Hash table does not have 2 entries")
	}
	_, ok := table["bla"]
	if !ok {
		t.Error("Table entry not present")
	}
}
