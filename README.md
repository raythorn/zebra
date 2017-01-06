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

falcon.Get("/user", user_get)
```
### Context
falcon provides a context which contains http.RespondWriter and http.Request for http, and a simple cache to store temporary
data, such as http request header, request parametr along with the url and form, named regexps and, of course, custom variables.
And it has several convenient APIs to handle http related jobs.
### Routers
falcon supports fixed route and regular expression route.

```go
falcon.Get("/user", handler) //fixed route

falcon.Get("/user/:id", handler) //regexp route, match /user/123 ..., and id will be set in context

falcon.Get("/user/(?P<name>exp)", handler) //named regexp route, name will be set in context
```
### Groups
falcon supports group api with same function.
```go
falcon.Group("/user"
	falcon.GGet("", handler),	//Match "/user"
	falcon.GPut("", handler),
	falcon.GSub("/friends",
		falcon.GGet("", handler), //Match "/user/friends"
		falcon.GPut("", handler),
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
Logger for falcon, it can print log to both console and file. It default logs to console, and you can user log.Add("file", "/tmp/log") to add log to file. Log file will be named with current date, and rotate to new file in 12:00pm, also cleanup older log files, system cached the latest files in a month.

## LICENSE

falcon source code is licensed under the [Apache Licence](http://www.apache.org/licenses/LICENSE-2.0.html), Version 2.0.
