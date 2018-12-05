package refuel

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"time"

	bolt "github.com/coreos/bbolt"
	database "github.com/reusing-code/dochan/db"
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

func AddFuelRecord(db *database.DB, record *RefuelRecord) error {
	err := db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(fuelBucket))

		keyInt, _ := bucket.NextSequence()

		key := database.Itob(keyInt)

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

func GetFuelRecord(db *database.DB, key uint64) (*RefuelRecord, error) {
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

func GetAllFuelRecords(db *database.DB, cb func(key uint64, record *RefuelRecord)) error {
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

func ParseCSV(db *database.DB, csv []byte) error {
	buf := bytes.NewBuffer(csv)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ";")
		if len(tokens) < 4 {
			continue
		}
		cents, err := strconv.ParseInt(strings.Replace(tokens[1], ",", "", -1), 10, 32)
		if err != nil {
			return err
		}
		kg, err := strconv.ParseFloat(strings.Replace(tokens[2], ",", ".", -1), 32)
		if err != nil {
			return err
		}
		km, err := strconv.ParseInt(tokens[3], 10, 32)
		if err != nil {
			return err
		}
		if tokens[4] == "" {
			tokens[4] = "0"
		}
		ikm, err := strconv.ParseInt(tokens[4], 10, 32)
		if err != nil {
			return err
		}
		date, err := time.Parse("02.01.2006", tokens[0])
		if err != nil {
			return err
		}
		record := &RefuelRecord{Date: date, CostCent: int(cents),
			FuelKG: float32(kg), TotalKM: int(km), Lat: 0, Lon: 0, IgnoreKM: int(ikm)}
		err = AddFuelRecord(db, record)
		if err != nil {
			return nil
		}
	}
	return nil
}
