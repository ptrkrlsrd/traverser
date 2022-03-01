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

type Store interface {
	LoadRoutes() (routes Routes, err error)
	AddRoute(Route) error
}

// Storage contains details about the Bolt DB
type Storage struct {
	db *badger.DB
}

// NewDB creates a new Bolt DB
func NewDB(path string) (*badger.DB, error) {
	expandedPath, err := tilde.Expand(path)
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(expandedPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewStorage creates a new Storage struct
func NewStorage(path string, db *badger.DB) (Store, error) {
	return &Storage{db: db}, nil
}

//LoadRoutes loads the routes from the storage
func (storage *Storage) LoadRoutes() (routes Routes, err error) {
	err = storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		routeData, err := readBytesFromIterator(it)
		for _, v := range routeData {
			route, err := NewRouteFromBytes(v)
			if err != nil {
				return fmt.Errorf("failed reading route from bytes: %v", err)
			}

			routes = append(routes, route)
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	return routes, nil
}

func readBytesFromIterator(it *badger.Iterator) ([][]byte, error) {
	data := [][]byte{}
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			data = append(data, v)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// AddRoute adds a route the database
func (storage *Storage) AddRoute(route Route) error {
	jsonData, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %v", err)
	}

	err = storage.db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(route.Alias), jsonData)
		return err
	})

	return err
}
