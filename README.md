# Redis Storage for [OAuth 2.0](https://github.com/go-oauth2/oauth2)

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Install

``` bash
$ go get -u -v github.com/go-oauth2/redis/v4
```

## Usage

``` go
package main

import (
	"github.com/go-redis/redis/v9"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
)

func main() {
	manager := manage.NewDefaultManager()
	
	// use redis token store
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 15,
	}))

	// use redis cluster store
	// manager.MapTokenStorage(oredis.NewRedisClusterStore(&redis.ClusterOptions{
	// 	Addrs: []string{"127.0.0.1:6379"},
	// 	DB: 15,
	// }))
}
```

## Testing

Testing requires a redis db. If you already have one, go ahead and run `go test ./...` like normal.
If you don't already have one, or don't want to use it, you can run in docker with `docker compose run --rm test && docker compose down`.

## MIT License

```
Copyright (c) 2020 Lyric
```

[Build-Status-Url]: https://travis-ci.org/go-oauth2/redis
[Build-Status-Image]: https://travis-ci.org/go-oauth2/redis.svg?branch=master
[codecov-url]: https://codecov.io/gh/go-oauth2/redis
[codecov-image]: https://codecov.io/gh/go-oauth2/redis/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/go-oauth2/redis/v4
[reportcard-image]: https://goreportcard.com/badge/github.com/go-oauth2/redis/v4
[godoc-url]: https://godoc.org/github.com/go-oauth2/redis/v4
[godoc-image]: https://godoc.org/github.com/go-oauth2/redis/v4?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
