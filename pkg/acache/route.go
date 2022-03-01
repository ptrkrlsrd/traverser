// Copyright © 2021 Petter Karlsrud
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
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Route contains a route that will be served by the Server
type Route struct {
	ID       string           `json:"key"`
	URL      string           `json:"url"`
	Alias    string           `json:"alias"`
	Method   string           `json:"method"`
	Response StorableResponse `json:"response"`
}

func NewRouteFromResponse(url, alias, method string, res *http.Response) (Route, error) {
	key := createKey(alias)
	response, err := NewStorableResponse(res)
	if err != nil {
		return Route{}, err
	}

	return Route{
		ID:       key,
		URL:      url,
		Alias:    alias,
		Method:   http.MethodGet,
		Response: response,
	}, nil
}

func NewRouteFromURL(url string, alias string) (Route, error) {
	res, err := http.Get(url)
	if err != nil {
		return Route{}, err
	}

	key := createKey(alias)
	response, err := NewStorableResponse(res)
	if err != nil {
		return Route{}, err
	}

	return Route{
		ID:       key,
		URL:      url,
		Alias:    alias,
		Method:   http.MethodGet,
		Response: response,
	}, nil
}

func createKey(alias string) string {
	sha1er := sha1.New()
	sha1er.Write([]byte(alias))
	key := hex.EncodeToString(sha1er.Sum(nil))
	return key
}

// NewRouteFromBytes RouteFromBytes creates a new route from a slice of bytes.
// Used to retrive a route from the database
func NewRouteFromBytes(bytes []byte) (Route, error) {
	var route Route
	err := json.Unmarshal(bytes, &route)
	if err != nil {
		return route, err
	}

	return route, nil
}

// Routes is a type that represents a slice of Routes
type Routes []Route

//ContainsURL ContainsURL returns true if the slice of routes contains an URL
func (routes *Routes) ContainsURL(url string) bool {
	for _, v := range *routes {
		if v.URL == url {
			return true
		}
	}

	return false
}

// ToString converts a slice of routes to a string with a newline to make printing easier
func (routes Routes) ToString() string {
	var output string
	for i, v := range routes {
		output += fmt.Sprintf("%d) %s -> %s\n", i, v.URL, v.Alias)
	}

	return output
}

// Print prints info about a slice of routes
func (routes Routes) Print() {
	fmt.Print(routes.ToString())
}

//PrintInfo prints info about all the routes
func (routes Routes) PrintInfo() {
	for i, v := range routes {
		fmt.Printf("%d) %s\n\tAlias: %s\n\tKey: %s\n\tMethod: %s\n\tHeaders:\n", i, v.URL, v.Alias, v.ID, v.Method)
		for k, h := range v.Response.Header {
			fmt.Printf("\t\t%s: %s\n", k, strings.Join(h, " "))
		}
	}
}
