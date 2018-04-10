# Acache
Simple API cacher and server written in Go

## Usage
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
Note: You must supply *https://* to the first URL
```
acache add https://api.coinmarketcap.com/v1/ticker/ eth
```

### Serve:
```
acache serve
```

# Tech
- [Go](https://golang.org/) <3
- [Cobra](https://github.com/spf13/cobra)
- [BoltDB](https://github.com/coreos/bbolt)
