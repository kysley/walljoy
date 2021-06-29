package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

type Store struct {
	db *bolt.DB
}

var (
	bucketName = []byte("wj")
)

func Open(path string) (*Store, error) {
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}
	if db, err := bolt.Open(path, 0600, opts); err != nil {
		return nil, err
	} else {
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketName)
			return err
		})
		if err != nil {
			return nil, err
		} else {
			return &Store{db: db}, nil
		}
	}
}

func (db *Store) Set(key string, value interface{}) error {
	if value == nil {
		return errors.New("no value")
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return db.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), buf.Bytes())
	})
}

func (db *Store) Get(key string, value interface{}) error {
	return db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, v := c.Seek([]byte(key)); k == nil || string(k) != key {
			return errors.New("no value")
		} else if value == nil {
			return nil
		} else {
			d := gob.NewDecoder(bytes.NewReader(v))
			return d.Decode(value)
		}
	})
}
