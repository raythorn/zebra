## Zebra

zebra is a simple but convenient framework for quickly developing RESTful APIs web applications/services in Go.

## Quick Start
#### Download and install

	go get github.com/raythorn/zebra

###### Create file `test.go`
```go
package main

import "github.com/raythorn/zebra"

func main() {
	zebra.Run()
}

```
###### Build and run
```bash
go build test.go
./test
```
######Congratulations! 
You just built your first zebra app.

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
zebra provides RESTful APIs, you can easily writing RESTful hanlder with Get/Put/Post/Delete/Patch/Options/Head/Any.
Take Get for example:
```go
import (
	"github.com/raythorn/zebra"
	"github.com/raythorn/zebra/context"
)

func user_get(ctx *context.Context) {
	ctx.WriteString("zebra")	
}

zebra.Get("/user", user_get)
```
### Context
zebra provides a context which contains http.RespondWriter and http.Request for http, and a simple cache to store temporary
data, such as http request header, request parametr along with the url and form, named regexps and, of course, custom variables.
And it has several convenient APIs to handle http related jobs.
### Routers
zebra supports fixed route and regular expression route.

```go
zebra.Get("/user", handler) //fixed route

zebra.Get("/user/:id", handler) //regexp route, match /user/123 ..., and id will be set in context

zebra.Get("/user/(?P<name>exp)", handler) //named regexp route, name will be set in context
```
### Groups
zebra supports group api with same function.
```go
zebra.Group("/user"
	zebra.GGet("", handler),	//Match "/user"
	zebra.GPut("", handler),
	zebra.GSub("/friends",
		zebra.GGet("", handler), //Match "/user/friends"
		zebra.GPut("", handler),
	),
)
```
GGet/GGPut/... is same as Get/Put... APIs, which add related route to group, and GSub can add a sub-group to current group.

## Authority
### API Signature
### Token

## Cache
### Ant
### Redis

## Database
### MongoDB

## Log
Logger for zebra, it can print log to both console and file. It default logs to console, and you can user log.Add("file", "/tmp/log") to add log to file. Log file will be named with current date, and rotate to new file in 12:00pm, also cleanup older log files, system cached the latest files in a month.

## LICENSE

zebra source code is licensed under the [Apache Licence](http://www.apache.org/licenses/LICENSE-2.0.html), Version 2.0.
