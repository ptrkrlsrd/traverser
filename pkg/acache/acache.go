// Copyright © 2018 Petter Karlsrud petterkarlsrud@me.com
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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/bolt"
	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/utilities/pkg/ucrypt"
	"github.com/ptrkrlsrd/utilities/pkg/unet"
)

const (
	BoltBucketName = "acache"
)

type Route struct {
	ID          string `json:"key"`
	URL         string `json:"url"`
	Alias       string `json:"alias"`
	Data        []byte `json:"data"`
	ContentType string `json:"contentType"`
}

func RouteFromBytes(bytes []byte) (Route, error) {
	var cacheItem Route
	err := json.Unmarshal(bytes, &cacheItem)
	if err != nil {
		return cacheItem, err
	}

	return cacheItem, nil
}

type Store struct {
	DB *bolt.DB
}

func NewCache(db *bolt.DB) Store {
	store := Store{DB: db}

	return store
}

func (store *Store) InitBucket() error {
	return store.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("acache"))
		if err != nil {
			return fmt.Errorf("Failed creating bucket with error: %s", err)
		}

		return nil
	})
}

func (store *Store) ListRoutes() (string, error) {
	var output string
	cacheItems, err := store.GetRoutes()
	if err != nil {
		return "", err
	}

	for i, v := range cacheItems {
		output += fmt.Sprintf("%d) %s -> %s\n", i, v.URL, v.Alias)
	}
	return output, nil
}

func (store *Store) Info() {
	cacheItems, err := store.GetRoutes()
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range cacheItems {
		fmt.Printf("%d) %s\n\tAlias: %s\n\tKey: %s\n\tContent-Type: %s\n", i, v.URL, v.Alias, v.ID, v.ContentType)
	}
}

func (store Store) GetRoutes() ([]Route, error) {
	var cacheItems []Route

	err := store.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBucketName))
		if b == nil {
			return fmt.Errorf("Could not find bucket %s", BoltBucketName)
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			cacheItem, err := RouteFromBytes(v)
			if err != nil {
				return err
			}

			cacheItems = append(cacheItems, cacheItem)
		}

		return nil
	})

	return cacheItems, err
}

func (store *Store) HasRoute(url string) (bool, error) {
	var hasRoute = false

	routes, err := store.GetRoutes()
	if err != nil {
		return false, err
	}

	for _, v := range routes {
		if v.URL == url {
			hasRoute = true
		}
	}

	return hasRoute, nil
}

func generateAlias(url string) (string, error) {
	splitURL, err := unet.GetURLPath(url)
	if err != nil {
		return "", nil
	}

	if len(splitURL) > 1 {
		return splitURL[1], nil
	}

	hash := ucrypt.MD5Hash(url)
	return hash, nil
}

func fetchItem(url string) ([]byte, *http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, res, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, res, err
	}

	return body, res, nil
}

func (store *Store) AddRoute(url string, alias string) error {
	data, resp, err := fetchItem(url)
	key := ucrypt.MD5Hash(alias)

	cacheItem := Route{
		ID:          key,
		URL:         url,
		Alias:       alias,
		Data:        data,
		ContentType: resp.Header.Get("Content-Type"),
	}

	jsonData, err := json.Marshal(cacheItem)
	if err != nil {
		return err
	}

	err = store.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBucketName))
		if b == nil {
			err = fmt.Errorf("failed to update the DB. Have you run 'acache init' yet?")
			return err
		}

		err := b.Put([]byte(key), jsonData)
		return err
	})

	return err
}

func (store *Store) ClearDB() error {
	err := store.DB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(BoltBucketName))
		return err
	})

	return err
}

func (store *Store) StartServer(port string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	cacheItems, err := store.GetRoutes()
	if err != nil {
		return err
	}

	for _, v := range cacheItems {
		router.GET(v.Alias, func(c *gin.Context) {
			c.Header("Content-Type", v.ContentType)
			c.String(http.StatusOK, string(v.Data))
		})
	}

	if err := router.Run(":" + port); err != nil {
		return err
	}

	return nil
}
