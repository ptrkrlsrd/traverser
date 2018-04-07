# Acache
Simple API cacher and server written in Go

Usage
- Init
```
acache init
```

- List routes:
```
acache list
```

- Add routes:
  - Note: You must supply *https://*
```
acache add https://api.coinmarketcap.com/v1/ticker/ eth
```
- Serve:
```
acache serve
```

# Tech
- [Go](https://golang.org/) <3
- [Cobra](https://github.com/spf13/cobra)
- [BoltDB](https://github.com/coreos/bbolt)
