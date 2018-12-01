package acache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Route Route
type Route struct {
	ID          string `json:"key"`
	URL         string `json:"url"`
	Alias       string `json:"alias"`
	Data        []byte `json:"data"`
	ContentType string `json:"contentType"`
}

//ContainsURL ContainsURL returns true if the slice of routes contains an URL
func (routes *Routes) ContainsURL(url string) (bool, error) {
	for _, v := range *routes {
		if v.URL == url {
			return true, nil
		}
	}

	return false, nil
}

// Routes Routes
type Routes []Route

// NewRouteFromBytes RouteFromBytes...
func NewRouteFromBytes(bytes []byte) (Route, error) {
	var routes Route
	err := json.Unmarshal(bytes, &routes)
	if err != nil {
		return routes, err
	}

	return routes, nil
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

func fetchItem(url string) ([]byte, *http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	return body, res, err
}
