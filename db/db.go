package db

import (
	"errors"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

type DB struct {
	handle *bolt.DB
}

const (
	hashBucket = "hashes"
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
