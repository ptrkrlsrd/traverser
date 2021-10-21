package acache

import "net/http"

func NewRouteFromRequest(url string, alias string) (Route, error) {
	res, err := http.Get(url)
	return NewRoute(url, alias, http.MethodGet, res), err
}
