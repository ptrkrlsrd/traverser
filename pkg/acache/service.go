// Copyright © 2018 Petter Karlsrud
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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/utilities/utext"
)

// Service contains the dependencies and handles the logic
type Service struct {
	Storage Storage
}

// NewService returns a new Service
func NewService(storage Storage) Service {
	return Service{Storage: storage}
}

func newRoute(url, alias, method string, data []byte, header http.Header) Route {
	key := utext.MD5Hash(alias)
	return Route{
		ID:     key,
		URL:    url,
		Alias:  alias,
		Method: method,
		Data:   data,
		Header: header,
	}
}

//NewRouteFromPostRequest adds a new route and stores the returned data into the database
func (service *Service) NewRouteFromPostRequest(url string, alias string, data []byte) (Route, error) {
	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Route{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Route{}, err
	}

	route := newRoute(url, alias, http.MethodPost, body, res.Header)
	return route, err
}

//NewRouteFromGetRequest ...
func (service *Service) NewRouteFromGetRequest(url string, alias string) (Route, error) {
	res, err := http.Get(url)
	if err != nil {
		return Route{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Route{}, err
	}

	route := newRoute(url, alias, http.MethodGet, body, res.Header)
	return route, err
}

// StoreRoute adds a route to the Storage
func (service *Service) StoreRoute(route Route) error {
	jsonData, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %v", err)
	}

	return service.Storage.Add(route.ID, jsonData)
}

//StartServer starts the API server
func (service *Service) StartServer(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	for _, v := range service.Storage.Routes {
		router.GET(v.Alias, func(c *gin.Context) {
			for k, h := range v.Header {
				c.Header(k, strings.Join(h, ","))
			}
			c.String(http.StatusOK, string(v.Data))
		})
	}

	return router.Run(addr)
}
