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
