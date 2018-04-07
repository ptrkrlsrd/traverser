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
acache add https://github.com github.com
```
- Serve:
```
acache serve
```

# Tech
- Go <3
- Cobra 
- BoltDB
