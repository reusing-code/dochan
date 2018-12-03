package refuel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
	database "github.com/reusing-code/dochan/db"
)

const (
	fuelBucket = "fuel"
)

type RefuelRecord struct {
	Date     time.Time
	CostCent int
	FuelKG   float32
	TotalKM  int
	Lat      float32
	Lon      float32
}

func AddFuelRecord(db database.DB, date time.Time, costCent int,
	fuelKG float32, totalKM int, lat float32, lon float32) error {
	err := db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fuelBucket))

		keyInt, _ := bucket.NextSequence()

		key := database.Itob(keyInt)
		r := &RefuelRecord{Date: date, CostCent: costCent,
			FuelKG: fuelKG, TotalKM: totalKM, Lat: lat, Lon: lon}

		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		err := enc.Encode(r)
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

func GetFuelRecord(db database.DB, key uint64) (*RefuelRecord, error) {
	r := &RefuelRecord{}
	err := db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fuelBucket))
		b := bucket.Get(database.Itob(key))
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

func GetAllFiles(db database.DB, cb func(key uint64, record *RefuelRecord)) error {
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
			cb(database.Btoi(k), &r)
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
