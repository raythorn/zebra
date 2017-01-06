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
* [RESTful API](#RESTful API)
	* [Routers](#Routers)
	* [Groups](#Groups)
* [Authority](#Authority)
	* [API Signature](#API Signature)
	* [Token](#Token)
* [Cache](#Cache)
	* [Ant](#Ant)
	* [Redis](#Redis)
* [Database](#Database)
	* [MongoDB](#MongoDB)
* [Log](#Log)

## RESTful API
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
