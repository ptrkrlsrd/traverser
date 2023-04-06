# Acache
![Go](https://github.com/ptrkrlsrd/acache/workflows/Go/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ptrkrlsrd_acache&metric=alert_status)](https://sonarcloud.io/dashboard?id=ptrkrlsrd_acache)

## What is Acache?

Acache is a tool used for storing responses from endpoints locally, and then serving them from your own computer. This is useful when you want to work on your solutions without access to a certain API when you're for example offline.  

## CLI
```
API response recorder

Usage:
  acache [command]

Available Commands:
  add         Add a new route.
                Example: "acache add https://pokeapi.co/api/v2/pokemon/ditto /ditto"
                Here the first argument is the path to the endpoint you want to cache,
                and the last is the alias. Note that you can also add from a json file by replacing
        the first URL with a relative path to a json file.
  clear       Clears the database containing the stored routes
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  info        Print route information
  list        List all routes
  proxy       Start Acache as a proxy between you and another API and save the responses locally
  serve       Load the stored routes from cache and serve the API

Flags:
      --config string      Config file (default "~/.config/acache/acache.json")
      --d string           Database (default "~/.config/acache/")
  -h, --help               help for acache
  -y, --use-yaml           Use YAML storage
      --yaml-path string   Yaml storage path (default "./routes/")

```

### Installation
```
go install github.com/ptrkrlsrd/acache@latest
```

### Usage
#### Approach 1. Manually entering routes
* Add routes
```
acache add <url> <alias>
acache add https://api.coinmarketcap.com/v1/ticker/eth /v1/eth
```
#### Approach 2. Using Acache as a proxy to store incoming API requests
```
acache proxy https://api.coinmarketcap.com/
```
This will create a proxy between you and the API which you can call by for example running `curl localhost:3000/v1/eth` which internally fetches `https://api.coinmarketcap.com/v1/ticker/eth` and stores the response into BadgerDB.


### Serving the stored API endpoints
* Start the server by running:
```
$ acache serve
```

* Perform curl against aliased routes served by Acache
```
$ curl localhost:3000/v1/eth
```

# Tech
- [Go](https://golang.org/) <3
- [Cobra](https://github.com/spf13/cobra)
- [Badger DB](https://github.com/dgraph-io/badger)
- I would also thank Github CoPilot for awesome suggestions while creating this tool
