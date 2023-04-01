package acache

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlStorage struct {
	path string
}

func NewYAMLStorage(path string) (RouteStorer, error) {
	return &yamlStorage{path: path}, nil
}

func (storage *yamlStorage) GetRoute(filter RouteFilter) (Route, error) {
	routes, err := storage.GetRoutes()
	if err != nil {
		return Route{}, err
	}

	for _, route := range routes {
		if (route.Alias == filter.Alias) || (route.URL == filter.URL) {
			return route, nil
		}
	}

	return Route{}, errors.New("route not found")
}

func (storage *yamlStorage) GetRoutes() (Routes, error) {
	routes := Routes{}

	data, err := ioutil.ReadFile(storage.path)
	if err != nil {
		return routes, err
	}

	err = yaml.Unmarshal(data, &routes)
	if err != nil {
		return routes, err
	}

	return routes, nil
}

func (storage *yamlStorage) AddRoute(route Route) error {
	routes, err := storage.GetRoutes()
	if err != nil {
		return err
	}

	for _, r := range routes {
		if r.Alias == route.Alias {
			return errors.New("route already exists")
		}
	}

	routes = append(routes, route)

	data, err := yaml.Marshal(&routes)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(storage.path, data, 0644)
}

func (storage *yamlStorage) Clear() error {
	return ioutil.WriteFile(storage.path, []byte{}, 0644)
}
