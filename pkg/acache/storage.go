// Copyright © 2019 Petter Karlsrud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Init Initializes the Bolt database
func (storage *Storage) Init() error {
	return storage.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(storage.BucketName))
		if err != nil {
			return fmt.Errorf("failed creating bucket with error: %v", err)
		}

		return nil
	})
}

//GetRoutes gets the routes from the Bbolt bucket and sets it [TODO: Return instead of setting]
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
