package db

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	bolt "github.com/coreos/bbolt"
)

type DB struct {
	Handle    *bolt.DB
	hashTable map[string]bool
}

type DBFile struct {
	Name       string
	Path       string
	ImportDate time.Time
	RawData    []byte
	Content    []string
}

const (
	hashBucket = "hashes"
	hashKey    = "hashes"
	fileBucket = "files"
)

func New(path string) (*DB, error) {
	result := &DB{}
	var err error
	result.Handle, err = bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = result.Handle.Update(func(tx *bolt.Tx) error {
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
	result.hashTable, err = result.loadHashTable()
	return result, nil
}

func (db *DB) Close() error {
	if db != nil && db.Handle != nil {
		return db.Handle.Close()
	}
	return errors.New("No DB")
}

func (db *DB) loadHashTable() (map[string]bool, error) {
	result := make(map[string]bool)
	var b []byte
	err := db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hashBucket))
		b = bucket.Get([]byte(hashKey))
		return nil
	})
	if err != nil {
		return result, err
	}
	if b == nil || len(b) == 0 {
		return result, nil
	}
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err = dec.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *DB) storeHashTable() error {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(db.hashTable)
	if err != nil {
		return err
	}
	err = db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hashBucket))
		err := bucket.Put([]byte(hashKey), buf.Bytes())
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Contains(hash string) bool {
	_, ok := db.hashTable[hash]
	return ok
}

func (db *DB) AddFile(path string, hash string, rawData []byte, content []string) error {
	err := db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fileBucket))

		keyInt, _ := bucket.NextSequence()

		filename := filepath.Base(path)
		key := Itob(keyInt)
		f := &DBFile{Name: filename, Path: path, RawData: rawData, Content: content, ImportDate: time.Now()}

		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		err := enc.Encode(f)
		if err != nil {
			return err
		}
		err = bucket.Put(key, buf.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	db.hashTable[hash] = true
	db.storeHashTable()
	return nil
}

func (db *DB) GetAllFiles(cb func(key uint64, file DBFile)) error {
	err := db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fileBucket))
		err := bucket.ForEach(func(k, v []byte) error {
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			var f DBFile
			err := dec.Decode(&f)
			if err != nil {
				return err
			}
			cb(Btoi(k), f)
			return nil

		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetFile(key uint64) (*DBFile, error) {
	f := &DBFile{}
	err := db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fileBucket))
		b := bucket.Get(Itob(key))
		if b == nil {
			return fmt.Errorf("Document with key %v not found", key)
		}
		buf := bytes.NewBuffer(b)
		dec := gob.NewDecoder(buf)

		err := dec.Decode(f)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func Btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
