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
	"log"
	"net/http"

	"github.com/coreos/bolt"
	"github.com/gin-gonic/gin"
)

const (
	BoltBucketName = "acache"
)

type CacheItem struct {
	ID    string `json:"key"`
	URL   string `json:"url"`
	Alias string `json:"alias"`
	Data  []byte `json:"data"`
}

func CacheItemFromBytes(bytes []byte) (CacheItem, error) {
	var cacheItem CacheItem
	err := json.Unmarshal(bytes, &cacheItem)
	if err != nil {
		return cacheItem, err
	}

	return cacheItem, nil
}

type CacheStore struct {
	DB *bolt.DB
}

func NewCache(db *bolt.DB) CacheStore {
	cacheStore := CacheStore{DB: db}

	return cacheStore
}

func (cacheStore *CacheStore) InitBucket() {
	cacheStore.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("acache"))
		if err != nil {
			return fmt.Errorf("Failed creating bucket with error: %s", err)
		}
		return nil
	})
}

func (cacheStore *CacheStore) ListRoutes() {
	cacheItems, _ := cacheStore.GetCacheItems()
	for i, v := range cacheItems {
		fmt.Printf("%d) %s -> %s\n", i, v.URL, v.Alias)
	}
}

func (cacheStore *CacheStore) Info() {
	cacheItems, _ := cacheStore.GetCacheItems()
	for i, v := range cacheItems {
		fmt.Printf("%d) %s\n\tAlias: %s\n\tKey: %s\n", i, v.URL, v.Alias, v.ID)
	}
}

func (cacheStore CacheStore) GetCacheItems() ([]CacheItem, error) {
	var cacheItems []CacheItem

	err := cacheStore.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			cacheItem, err := CacheItemFromBytes(v)
			if err != nil {
				log.Fatal(err)
			}

			cacheItems = append(cacheItems, cacheItem)
		}

		return nil
	})

	return cacheItems, err
}

func (cacheStore *CacheStore) AddRoute(url string, alias string) error {
	data := fetchJSON(url)
	key := md5Hash(alias)

	cacheItem := CacheItem{ID: key, URL: url, Alias: alias, Data: data}
	jsonData, err := json.Marshal(cacheItem)

	if err != nil {
		log.Println(err)
		return err
	}

	err = cacheStore.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBucketName))
		if b == nil {
			err = fmt.Errorf("Failed to update the DB. Have you run acache init?")
			log.Println(err)
			return err
		}

		err := b.Put([]byte(key), jsonData)
		return err
	})

	return err
}

func (cacheStore *CacheStore) ClearDB() error {
	err := cacheStore.DB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(BoltBucketName))
		return err
	})

	return err
}

func (cacheStore *CacheStore) StartServer(port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	cacheItems, _ := cacheStore.GetCacheItems()

	for _, v := range cacheItems {
		router.GET(v.Alias, func(c *gin.Context) {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.String(http.StatusOK, string(v.Data))
		})
	}

	router.Run(":" + port)
}
