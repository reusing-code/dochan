package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

type DB struct {
	handle *bolt.DB
}

const (
	hashBucket = "hashes"
	hashKey    = "hashes"
	fileBucket = "files"
)

func New(path string) (*DB, error) {
	result := &DB{}
	var err error
	result.handle, err = bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = result.handle.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(hashBucket))
		if err != nil {
			return fmt.Errorf("create bucket %q: %q", hashBucket, err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(fileBucket))
		if err != nil {
			return fmt.Errorf("create bucket %q: %q", fileBucket, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) Close() error {
	if db != nil && db.handle != nil {
		return db.handle.Close()
	}
	return errors.New("No DB")
}

func (db *DB) GetHashTable() (map[string]bool, error) {
	result := make(map[string]bool)
	var b []byte
	err := db.handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hashBucket))
		b = bucket.Get([]byte(hashKey))
		return nil
	})
	if err != nil {
		return result, err
	}
	if len(b) == 0 {
		return result, nil
	}
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err = dec.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *DB) SetHashTable(table map[string]bool) error {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(table)
	if err != nil {
		return err
	}
	err = db.handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hashBucket))
		err := bucket.Put([]byte(hashKey), buf.Bytes())
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
