package acache

import (
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

// badgerStorage is a wrapper around badger.DB that implements the Store interface
type badgerStorage struct {
	db *badger.DB
}

// NewBadgerDB creates a new Badger DB
func NewBadgerDB(path string) (*badger.DB, error) {
	expandedPath, err := tilde.Expand(path)
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(expandedPath)
	opts.Logger = nil
	return badger.Open(opts)
}

// NewBadgerStorage creates a new Storage struct
func NewBadgerStorage(db *badger.DB) (RouteStorer, error) {
	return &badgerStorage{db: db}, nil
}

//GetRoutes loads the routes from the storage and returns them
func (storage *badgerStorage) GetRoutes() (routes Routes, err error) {
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
func (storage *badgerStorage) AddRoute(route Route) error {
	jsonData, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %v", err)
	}

	return storage.db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(route.Alias), jsonData)
		return err
	})
}
