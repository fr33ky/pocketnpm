package db

import (
	"log"
	"path/filepath"

	"encoding/binary"

	"github.com/boltdb/bolt"
)

// PocketBase type is a frontend for BoltDB
type PocketBase struct {
	db   *bolt.DB
	path string
}

// NewPocketBase creates a new PocketBase object
func NewPocketBase(config *DatabaseConfig) *PocketBase {
	path, _ := filepath.Abs(config.Path)
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatalf("Failed to load database file: %s", err)
	}

	pb := &PocketBase{
		db:   db,
		path: path,
	}

	return pb
}

// Close method closes database connection
func (pb *PocketBase) Close() {
	pb.db.Close()
}

// Init method initializes scheme
//
// Globals bucket contains global variables such as "sequence"
// Packages bucket contains the map of "id": "revision" hash
// Marks bucket contains the state to determine whether package currently downloaded or not
// Documents bucket contains full document of the package
// Files bucket con
func (pb *PocketBase) Init() {
	pb.db.Update(func(tx *bolt.Tx) error {
		global, _ := tx.CreateBucketIfNotExists([]byte("Globals"))
		tx.CreateBucketIfNotExists([]byte("Packages"))
		tx.CreateBucketIfNotExists([]byte("Marks"))
		tx.CreateBucketIfNotExists([]byte("Documents"))
		tx.CreateBucketIfNotExists([]byte("Files"))

		defaultSequence := make([]byte, 4)
		binary.LittleEndian.PutUint32(defaultSequence, 0)
		global.Put([]byte("sequence"), defaultSequence)
		return nil
	})
}

func (pb *PocketBase) IsInitialized() bool {
	initialized := false
	pb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Globals"))
		if b != nil {
			v := b.Get([]byte("sequence"))
			initialized = v != nil
		}
		return nil
	})

	return initialized
}
