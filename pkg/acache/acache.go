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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/utilities/ucrypt"
)

// Service Service..
type Service struct {
	Storage Storage
}

func NewService(storage Storage) Service {
	return Service{Storage: storage}
}

//AddRoute AddRoute...
func (service *Service) FetchRoute(url string, alias string) (Route, error) {
	data, resp, err := fetchItem(url)
	if err != nil {
		return Route{}, nil
	}

	key := ucrypt.MD5Hash(alias)

	route := Route{
		ID:          key,
		URL:         url,
		Alias:       alias,
		Data:        data,
		ContentType: resp.Header.Get("Content-Type"),
	}

	return route, err
}

//AddRoute AddRoute...
func (service *Service) AddRoute(route Route) error {
	jsonData, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %v", err)
	}

	return service.Storage.Add(route.ID, jsonData)
}

//StartServer Start the API server
func (service *Service) StartServer(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	for _, v := range service.Storage.Routes {
		router.GET(v.Alias, func(c *gin.Context) {
			c.Header("Content-Type", v.ContentType)
			c.String(http.StatusOK, string(v.Data))
		})
	}

	return router.Run(addr)
}
