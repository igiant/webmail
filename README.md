# Kerio Connect API
[![Go Reference](https://pkg.go.dev/badge/github.com/igiant/webmail.svg)](https://pkg.go.dev/github.com/igiant/webmail)
## Overview
Client for [Webmail Kerio API Connect (JSON-RPC 2.0)](https://manuals.gfi.com/en/kerio/api/connect/client/reference/index.html)

Implemented all Client API for Kerio Connect methods

## Installation
```go
go get github.com/igiant/webmail
```

## Example
```go
package main

import (
	"fmt"
	"log"

	"github.com/igiant/connect"
)

func main() {
	config := connect.NewConfig("server_addr")
	conn, err := config.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	app := &connect.Application{
		Name:    "MyApp",
		Vendor:  "Me",
		Version: "v0.0.1",
	}
	err = conn.Login("user_name", "user_password", app)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Logout()
		if err != nil {
			log.Println(err)
		}
	}()
	info, err := conn.ContactsGetPersonal()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Name: %s\nComment: %s\n",
		info.CommonName,
		info.Comment,
	)
}
```
## Documentation
* [GoDoc](http://godoc.org/github.com/igiant/webmail)

## RoadMap
* Add tests and search errors
