package db

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
)

// Bolt DB reference
type Database struct {
	db *bolt.DB
}

var defaultBucket = []byte("default")

// Return newDatabase
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {

	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, nil, fmt.Errorf("failed to open database: %v", err)
	}

	db = &Database{db: boltDb}
	closeFunc = boltDb.Close

	if err = db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, err
	}

	return db, closeFunc, err
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

func (db *Database) SetKey(key string, value string) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), []byte(value))
	})
}

func (db *Database) GetKey(key string) ([]byte, error) {

	var result []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, nil
	}

	return nil, err
}
