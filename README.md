# Redis Storage for [OAuth 2.0](https://github.com/go-oauth2/oauth2)

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Install

``` bash
$ go get -u -v gopkg.in/go-oauth2/redis.v4
```

## Usage

``` go
package main

import (
	"github.com/go-redis/redis"
	oredis "gopkg.in/go-oauth2/redis.v4"
	"gopkg.in/oauth2.v3/manage"
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

## MIT License

```
Copyright (c) 2020 Lyric
```

[Build-Status-Url]: https://travis-ci.org/go-oauth2/redis
[Build-Status-Image]: https://travis-ci.org/go-oauth2/redis.svg?branch=master
[codecov-url]: https://codecov.io/gh/go-oauth2/redis
[codecov-image]: https://codecov.io/gh/go-oauth2/redis/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/gopkg.in/go-oauth2/redis.v4
[reportcard-image]: https://goreportcard.com/badge/gopkg.in/go-oauth2/redis.v4
[godoc-url]: https://godoc.org/gopkg.in/go-oauth2/redis.v4
[godoc-image]: https://godoc.org/gopkg.in/go-oauth2/redis.v4?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg