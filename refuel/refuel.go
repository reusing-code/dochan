package refuel

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
)

const (
	fuelBucket = "fuel"
)

type RefuelRecord struct {
	Date     time.Time `json:"date"`
	CostCent int       `json:"costCent"`
	FuelKG   float32   `json:"fuelKG"`
	TotalKM  int       `json:"totalKM"`
	Lat      float32   `json:"lat"`
	Lon      float32   `json:"lon"`
	IgnoreKM int       `json:"ignoreKM"`
}

type DB struct {
	Handle *bolt.DB
}

const (
	bucket = "fuel"
)

func New(path string) (*DB, error) {
	result := &DB{}
	var err error
	result.Handle, err = bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = result.Handle.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket %q: %q", bucket, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *DB) Close() error {
	if db != nil && db.Handle != nil {
		return db.Handle.Close()
	}
	return errors.New("No DB")
}

func (db *DB) AddFuelRecord(record *RefuelRecord) error {
	err := db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fuelBucket))

		keyInt, _ := bucket.NextSequence()

		key := Itob(keyInt)

		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		err := enc.Encode(record)
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
	return nil
}

func (db *DB) GetFuelRecord(key uint64) (*RefuelRecord, error) {
	r := &RefuelRecord{}
	err := db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fuelBucket))
		b := bucket.Get(Itob(key))
		if b == nil {
			return fmt.Errorf("Document with key %v not found", key)
		}
		buf := bytes.NewBuffer(b)
		dec := gob.NewDecoder(buf)

		err := dec.Decode(r)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (db *DB) GetAllFuelRecords(cb func(key uint64, record *RefuelRecord)) error {
	err := db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fuelBucket))
		err := bucket.ForEach(func(k, v []byte) error {
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			var r RefuelRecord
			err := dec.Decode(&r)
			if err != nil {
				return err
			}
			cb(Btoi(k), &r)
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

func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func Btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
