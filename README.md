# Acache
Simple API cacher and server written in Go

## Usage
```
Usage:
  acache [command]

Available Commands:
  add         Add a new route
  clear       Clear the database
  help        Help about any command
  info        Info about the routes
  init        Init BoltDB
  list        List all routes(aliases)
  serve       

Flags:
      --config string     Config file (default "~/.config/acache/acache.json")
  -d, --database string   Database (default "~/.config/acache/acache.db")
  -h, --help              help for acache

Use "acache [command] --help" for more information about a command.
```

### Installation
```
go get https://github.com/ptrkrlsrd/acache/
```

### Init
```
acache init
```

### List routes
```
acache list
```

### Add routes
```
acache add https://api.coinmarketcap.com/v1/ticker/ /eth
```

### Serve:
```
acache serve
```

# Tech
- [Go](https://golang.org/) <3
- [Cobra](https://github.com/spf13/cobra)
- [BoltDB](https://github.com/coreos/bbolt)
