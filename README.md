# Redis Storage for OAuth 2.0

> Based on the redis token storage

[![License][License-Image]][License-Url] 
[![ReportCard][ReportCard-Image]][ReportCard-Url] 
[![GoDoc][GoDoc-Image]][GoDoc-Url]

## Install

``` bash
$ go get -u github.com/go-oauth2/redis
```

## Usage

``` go
package main

import (
	"github.com/go-oauth2/redis"
	"gopkg.in/oauth2.v3/manage"
)

func main() {
	manager := manage.NewDefaultManager()
	// use redis token store
	manager.MustTokenStorage(redis.NewTokenStore(&redis.Config{
		Addr: "127.0.0.1:6379",
	}))

	// ...
}
```

## MIT License

```
Copyright (c) 2016 Lyric
```

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/go-oauth2/redis
[ReportCard-Image]: https://goreportcard.com/badge/github.com/go-oauth2/redis
[GoDoc-Url]: https://godoc.org/github.com/go-oauth2/redis
[GoDoc-Image]: https://godoc.org/github.com/go-oauth2/redis?status.svg