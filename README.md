# Traverser [trah-vair-say]

![Go](https://github.com/ptrkrlsrd/traverser/workflows/Go/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ptrkrlsrd_traverser&metric=alert_status)](https://sonarcloud.io/dashboard?id=ptrkrlsrd_traverser)

## What is Traverser?

Traverser is a tool used for storing responses from endpoints locally, and then serving them from your own computer. This is useful when you want to work on your solutions without access to a certain API when you're for example offline.  

## CLI
```
API response recorder

Usage:
  traverser [command]

Available Commands:
  add         Add a new route.
                Example: "traverser add https://pokeapi.co/api/v2/pokemon/ditto /ditto"
                Here the first argument is the path to the endpoint you want to cache,
                and the last is the alias. Note that you can also add from a json file by replacing
                the first URL with a relative path to a json file.
  clear       Clears the database containing the stored routes
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  info        Print route information
  list        List all routes
  proxy       Start Traverser as a proxy between you and another API and save the responses locally
  serve       Load the stored routes from cache and serve the API

Flags:
      --config string      Config file (default "~/.config/traverser/traverser.json")
      --d string           Database (default "~/.config/traverser/")
  -h, --help               help for traverser
  -y, --use-yaml           Use YAML storage
      --yaml-path string   Yaml storage path (default "./routes/")

```

### Installation
```
go install github.com/ptrkrlsrd/traverser@latest
```

### Usage
#### Approach 1. Manually entering routes
* Add routes
```
traverser add <url> <alias>
traverser add https://api.coinmarketcap.com/v1/ticker/eth /v1/eth
```
#### Approach 2. Using Traverser as a proxy to store incoming API requests
```
traverser proxy https://api.coinmarketcap.com/
```
This will create a proxy between you and the API which you can call by for example running `curl localhost:3000/v1/eth` which internally fetches `https://api.coinmarketcap.com/v1/ticker/eth` and stores the response.

#### Approach 3. Add from JSON files
You can also add a route from a JSON file containing the body of the response you want to add. The Content-Type header will be set to "application/json".
```
traverser add ./route.json /route
```



### Serving the stored API endpoints
* Start the server by running:
```
$ traverser serve
```

* Perform curl against aliased routes served by Traverser
```
$ curl localhost:3000/v1/eth
```

# Tech
- [Go](https://golang.org/) <3
- [Cobra](https://github.com/spf13/cobra)
- [Badger DB](https://github.com/dgraph-io/badger)
