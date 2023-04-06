package acache

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const FileExtension string = ".yaml"

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

	filePaths, err := storage.findFiles()
	if err != nil {
		return Routes{}, err
	}

	for _, i := range filePaths {
		data, err := os.ReadFile(i)
		if err != nil {
			return routes, err
		}
        var route Route

		err = yaml.Unmarshal(data, &route)
		if err != nil {
			return routes, err
		}

        routes = append(routes, route)
	}

	return routes, nil
}

func (storage *yamlStorage) findFiles() ([]string, error) {
    filesGlob := path.Join(storage.path, "*" + FileExtension)
	filePaths, err := filepath.Glob(filesGlob)
	if err != nil {
		return nil, err
	}
	return filePaths, nil
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

	data, err := yaml.Marshal(&route)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(storage.path, route.ID) + FileExtension, data, 0644)
}

func (storage *yamlStorage) Clear() error {
    filePaths, err := storage.findFiles()
    if err != nil {
        return err
    }
    for _, i := range filePaths {
        err := os.Remove(i)
        if err != nil {
            return err
        }
    }
    return nil
}
