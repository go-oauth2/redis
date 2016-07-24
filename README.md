Redis Storage for OAuth2
========================

[![GoDoc](https://godoc.org/github.com/go-oauth2/redis?status.svg)](https://godoc.org/github.com/go-oauth2/redis)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-oauth2/redis)](https://goreportcard.com/report/github.com/go-oauth2/redis)

Install
-------

``` bash
$ go get -u -v github.com/go-oauth2/redis
```

Usage
-----

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

License
-------

```
Copyright (c) 2016, OAuth 2.0
All rights reserved.
```