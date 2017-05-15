// Copyright 2016 Derek Ray. All rights reserved.
// Use of this source code is governed by Apache License 2.0
// that can be found in the LICENSE file.

// Package falcon is a simple wrap implement for develop http server.
package falcon

import (
	"github.com/raythorn/falcon/oss"
	"github.com/raythorn/falcon/router"
)

var (
	falcon *app
	Env    *Environment
)

func init() {
	falcon = &app{router.New(), &router.Group{}}
	Env = &Environment{data: make(map[string]string)}
}

//Run starts a http(s) server
func Run() {
	falcon.run()
}

//Insert midware to http server, which will be called before each request handled.
func Use(handler router.Midware) {
	falcon.Use(handler)
}

func Oss(pattern, root string, archive oss.Archive) {
	falcon.Oss(pattern, root, archive)
}

//Get add a GET handler, which used to get data from server
func Get(pattern string, handler router.Handler) {
	falcon.Get(pattern, handler)
}

//Patch add a PATCH handler, which used to patch existed data
func Patch(pattern string, handler router.Handler) {
	falcon.Patch(pattern, handler)
}

//Put add a PUT handler, which used to update data
func Put(pattern string, handler router.Handler) {
	falcon.Put(pattern, handler)
}

//Post add a POST handler, which used to create resource
func Post(pattern string, handler router.Handler) {
	falcon.Post(pattern, handler)
}

//Delete add a DELETE handler, which used to delete resource from server
func Delete(pattern string, handler router.Handler) {
	falcon.Delete(pattern, handler)
}

//Head add a HEAD handler
func Head(pattern string, handler router.Handler) {
	falcon.Head(pattern, handler)
}

//Options add a OPTIONS handler
func Options(pattern string, handler router.Handler) {
	falcon.Options(pattern, handler)
}

//Any add a ANY handler, which can response to all method
func Any(pattern string, handler router.Handler) {
	falcon.Any(pattern, handler)
}

//NotFound add a not found handler, which used to be the handler when request not found
func NotFound(handler router.Handler) {
	falcon.NotFound(handler)
}

//NotAllowed add a not allowed handler, which used to be the handler when request not allowed
func NotAllowed(handler router.Handler) {
	falcon.NotAllowed(handler)
}

//Group assemble handlers with same prefix together, routes can be routes and sub-groups, with
//group you can add midwares with Before and After, Before add midware to be called before
//handler called and After add midware to be called after handler called
func Group(prefix string, routes ...interface{}) *router.Group {
	return falcon.Router.Group(prefix, routes...)
}

//GSub add a sub-group
func GSub(prefix string, routes ...interface{}) *router.Group {
	return falcon.g.Sub(prefix, routes...)
}

//GGet add a grouped GET handler
func GGet(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Get(pattern, handler)
}

//GPatch add a grouped PATCH handler
func GPatch(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Patch(pattern, handler)
}

//GPut add a grouped PUT handler
func GPut(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Put(pattern, handler)
}

//GPost add a grouped POST handler
func GPost(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Post(pattern, handler)
}

//GDelete add a grouped DELETE handler
func GDelete(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Delete(pattern, handler)
}

//GHead add a grouped HEAD handler
func GHead(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Head(pattern, handler)
}

//GOptions add a grouped OPTIONS handler
func GOptions(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Options(pattern, handler)
}

//GAny add a grouped ANY handler
func GAny(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Any(pattern, handler)
}
