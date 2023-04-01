package acache

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Route contains a route that will be served by the Server
type Route struct {
	ID       cacheKey         `json:"key"`
	URL      string           `json:"url"`
	Alias    string           `json:"alias"`
	Method   string           `json:"method"`
	Response StorableResponse `json:"response"`
}

func NewRouteFromResponse(url, alias, method string, res *http.Response) (Route, error) {
	response, err := NewStorableResponse(res)
	if err != nil {
		return Route{}, err
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return Route{}, err
	}

	key, err := NewCacheKey(alias, req)
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
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Route{}, err
	}

	res, err := client.Do(request)
	if err != nil {
		return Route{}, err
	}

	response, err := NewStorableResponse(res)
	if err != nil {
		return Route{}, err
	}

	key, err := NewCacheKey(alias, request)
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

func contains[T Route](routes []T, fn func(item T) bool) bool {
	for _, value := range routes {
		if fn(value) {
			return true
		}
	}
	return false
}

// ContainsURL ContainsURL returns true if the slice of routes contains an URL
func (routes *Routes) ContainsURL(url string) bool {
	return contains(*routes, func(item Route) bool {
		return item.URL == url
	})
}

// ContainsAlias ContainsAlias returns true if the slice of routes contains an alias
func (routes *Routes) ContainsAlias(alias string) bool {
	return contains(*routes, func(item Route) bool {
		return item.Alias == alias
	})
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

// PrintInfo prints info about all the routes
func (routes Routes) PrintInfo() {
	for i, v := range routes {
		fmt.Printf("%d) %s\n\tAlias: %s\n\tMethod: %s\n\tHeaders:\n", i, v.URL, v.Alias, v.Method)
		for k, h := range v.Response.Headers {
			fmt.Printf("\t\t%s: %s\n", k, h)
		}
	}
}
