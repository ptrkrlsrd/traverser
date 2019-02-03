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
	"encoding/json"
	"fmt"
)

// Route Route
type Route struct {
	ID          string `json:"key"`
	URL         string `json:"url"`
	Alias       string `json:"alias"`
	Data        []byte `json:"data"`
	ContentType string `json:"contentType"`
}

// NewRouteFromBytes RouteFromBytes...
func NewRouteFromBytes(bytes []byte) (Route, error) {
	var routes Route
	err := json.Unmarshal(bytes, &routes)
	if err != nil {
		return routes, err
	}

	return routes, nil
}

// Routes Routes
type Routes []Route

//ContainsURL ContainsURL returns true if the slice of routes contains an URL
func (routes *Routes) ContainsURL(url string) (bool, error) {
	for _, v := range *routes {
		if v.URL == url {
			return true, nil
		}
	}

	return false, nil
}
func (routes Routes) ToString() string {
	var output string
	for i, v := range routes {
		output += fmt.Sprintf("%d) %s -> %s\n", i, v.URL, v.Alias)
	}

	return output
}

func (routes Routes) Print() {
	fmt.Print(routes.ToString())
}

//PrintAll PrintAll...
func (routes *Routes) PrintInfo() {
	for i, v := range *routes {
		fmt.Printf("%d) %s\n\tAlias: %s\n\tKey: %s\n\tContent-Type: %s\n", i, v.URL, v.Alias, v.ID, v.ContentType)
	}
}
