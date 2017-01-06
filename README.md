## Falcon

falcon is a simple but convenient framework for quickly developing RESTful APIs web applications/services in Go.

## Quick Start
#### Download and install

	go get github.com/raythorn/falcon

###### Create file `test.go`
```go
package main

import "github.com/raythorn/falcon"

func main() {
	falcon.Run()
}

```
###### Build and run
```bash
go build test.go
./test
```
######Congratulations! 
You just built your first falcon app.

## Features
* [RESTful API](#restful-api)
	* [Context](#context)
	* [Routers](#routers)
	* [Groups](#groups)
* [Authority](#authority)
	* [API Signature](#api-signature)
	* [Token](#token)
* [Cache](#cache)
	* [Ant](#ant)
	* [Redis](#redis)
* [Database](#database)
	* [MongoDB](#mongodb)
* [Log](#log)

## RESTful API
falcon provides RESTful APIs, you can easily writing RESTful hanlder with Get/Put/Post/Delete/Patch/Options/Head/Any.
Take Get for example:
```go
import (
	"github.com/raythorn/falcon"
	"github.com/raythorn/falcon/context"
)

func user_get(ctx *context.Context) {
	ctx.WriteString("Falcon")	
}

falcon.Get("/user", user_create)
```
### Context
### Routers
### Groups

## Authority
### API Signature
### Token

## Cache
### Ant
### Redis

## Database
### MongoDB

## Log

## LICENSE

falcon source code is licensed under the [Apache Licence](http://www.apache.org/licenses/LICENSE-2.0.html), Version 2.0.
