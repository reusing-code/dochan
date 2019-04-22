package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	bolt "github.com/coreos/bbolt"
)

type SessionDB struct {
	secret string
	Handle *bolt.DB
}

const (
	sessionBucket = "session"
	sessionBytes  = 64
)

func NewSessionHandler(path string, secret string) (*SessionDB, error) {
	result := &SessionDB{secret: secret}
	var err error
	result.Handle, err = bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = result.Handle.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(sessionBucket))
		if err != nil {
			return fmt.Errorf("create bucket %q: %q", sessionBucket, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *SessionDB) Close() error {
	if db != nil && db.Handle != nil {
		return db.Handle.Close()
	}
	return errors.New("No DB")
}

func (db *SessionDB) sessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	secret := string(buf)
	if err == nil && secret == db.secret {
		session, err := db.CreateSession()
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
		}
		w.Header().Add("X-Session-Token", session)
		w.Write([]byte("{}"))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func (db *SessionDB) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/session/create" {
			next.ServeHTTP(w, r)
			return
		}

		session := r.Header.Get("X-Session-Token")
		if session == "" || !db.GetSession(session) {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
			return
		}

	})
}

func (db *SessionDB) CreateSession() (string, error) {
	key := make([]byte, sessionBytes)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	err = db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(sessionBucket))
		err := bucket.Put(key, []byte(""))
		return err
	})
	if err != nil {
		return "", err
	}

	baseEncoded := base64.URLEncoding.EncodeToString(key)

	return baseEncoded, nil
}

func (db *SessionDB) DestroySession(baseKey string) error {
	key, err := base64.URLEncoding.DecodeString(baseKey)
	if err != nil {
		return err
	}
	err = db.Handle.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(sessionBucket))
		return bucket.Delete(key)
	})

	if err != nil {
		return err
	}
	return nil
}

func (db *SessionDB) GetSession(baseKey string) bool {
	key, err := base64.URLEncoding.DecodeString(baseKey)
	if err != nil {
		return false
	}
	sessionFound := false
	db.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(sessionBucket))
		b := bucket.Get(key)
		if b == nil {
			return nil
		}
		sessionFound = true
		return nil

	})

	return sessionFound
}
