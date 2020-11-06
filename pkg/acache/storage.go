// Copyright © 2020 github.com/ptrkrlsrd
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
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

// Storage contains details about the Bolt DB
type Storage struct {
	Path       string
	BucketName string
	Routes     Routes
	db         *badger.DB
}

// NewDB creates a new Bolt DB
func NewDB(path string) (*badger.DB, error) {
	expandedPath, err := tilde.Expand(path)
	db, err := badger.Open(badger.DefaultOptions(expandedPath))
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewStorage creates a new Storage struct
func NewStorage(bucketName, path string, db *badger.DB) (Storage, error) {
	return Storage{BucketName: bucketName, db: db}, nil
}

//LoadRoutes loads the routes from the storage
func (storage *Storage) LoadRoutes() (routes Routes, err error) {
	err = storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				route, err := NewRouteFromBytes(v)
				if err != nil {
					return fmt.Errorf("failed reading route from bytes: %v", err)
				}

				routes = append(routes, route)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	storage.Routes = routes
	return routes, nil
}

// AddRoute adds a route the database
func (storage *Storage) AddRoute(route Route) error {
	jsonData, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %v", err)
	}

	return storage.db.Update(func(tx *badger.Txn) error {
		if err := tx.Set([]byte(route.ID), jsonData); err != nil {
			return fmt.Errorf("failed marshaling JSON: %v", err)
		}

		return tx.Commit()
	})
}
