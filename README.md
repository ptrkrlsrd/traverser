# Acache
![Go](https://github.com/ptrkrlsrd/acache/workflows/Go/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ptrkrlsrd_acache&metric=alert_status)](https://sonarcloud.io/dashboard?id=ptrkrlsrd_acache)

## What is Acache?

Acache is a tool used for storing responses from endpoints locally, and then serving them from your own computer. This is useful when you want to work on your solutions without access to a certain API when you're for example offline.  

## Usage
```
API response recorder

Usage:
  acache [command]

Available Commands:
  add         Add a new route. 
                Example: "acache add https://pokeapi.co/api/v2/pokemon/ditto /ditto"
                Here the first argument is the path to the endpoint you want to cache, 
                and the last is the alias
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  info        Print route information
  list        List all routes(aliases)
  proxy       Start the server as a proxy between you and another API
  serve       Load the stored routes from cache and serve the API

Flags:
      --config string     Config file (default "~/.config/acache/acache.json")
  -d, --database string   Database (default "~/.config/acache/")
  -h, --help              help for acache

Use "acache [command] --help" for more information about a command.
```

### Installation
```
go get github.com/ptrkrlsrd/acache
```

### Add routes
```
acache add <url> <alias>
acache add https://api.coinmarketcap.com/v1/ticker/eth /v1/eth
```


### Server
```
$ acache serve
$ curl localhost:3000/v1/eth
```

# Tech
- [Go](https://golang.org/) <3
- [Cobra](https://github.com/spf13/cobra)
- [Badger DB](https://github.com/dgraph-io/badger)
