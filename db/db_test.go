package db

import (
	"encoding/json"
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

type fileStruct struct {
	Name string
	Bla  string
}

var testFiles = []struct {
	key  string
	file fileStruct
}{
	{"key1", fileStruct{"abc", "def"}},
	{"key2", fileStruct{"123", "123"}},
	{"key3", fileStruct{"F", "X"}},
}

func TestStoreFiles(t *testing.T) {
	defer os.Remove("test.db")
	db, err := New("test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	in := make(map[string][]byte)
	for _, tc := range testFiles {
		b, err := json.Marshal(tc.file)
		in[tc.key] = b
		if err != nil {
			t.Fatal(err)
		}
	}

	err = db.AddFiles(in)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string][]byte)
	err = db.GetAllFiles(out)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testFiles {
		var file fileStruct
		err := json.Unmarshal(out[tc.key], &file)
		if err != nil {
			t.Fatal(err)
		}
		if file != tc.file {
			t.Errorf("Value for key %q expected %q was %q", tc.key, tc.file, file)
		}
	}
}
