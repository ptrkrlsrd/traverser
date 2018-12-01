package acache

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

type Storage struct {
	Path       string
	BucketName string
	DB         *bolt.DB
	Routes     Routes
}

func NewDB(path string) (*bolt.DB, error) {
	expandedPath, err := tilde.Expand(path)
	db, err := bolt.Open(expandedPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func NewStorage(bucketName, path string, db *bolt.DB) (Storage, error) {
	return Storage{BucketName: bucketName, DB: db}, nil
}

func (storage *Storage) Init() error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(storage.BucketName))
		if err != nil {
			return fmt.Errorf("failed creating bucket with error: %v", err)
		}

		return nil
	})
}

//GetRoutes GetRoutes...
func (storage *Storage) GetRoutes() (err error) {
	var routes Routes
	err = storage.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(storage.BucketName))
		if b == nil {
			return fmt.Errorf("could not find bucket %s", storage.BucketName)
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if route, err := NewRouteFromBytes(v); err == nil {
				routes = append(routes, route)
			} else {
				return fmt.Errorf("failed reading route from bytes: %v", err)
			}
		}

		return nil
	})

	storage.Routes = routes
	return err
}

func (storage *Storage) Add(key string, data []byte) error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(storage.BucketName))
		if b == nil {
			return fmt.Errorf("failed to update the DB. Have you run 'acache init' yet?")
		}

		return b.Put([]byte(key), data)
	})
}

//Clear Clear...
func (storage *Storage) Clear() error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(storage.BucketName))
	})
}
