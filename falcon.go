// Copyright 2016 Derek Ray. All rights reserved.
// Use of this source code is governed by Apache License 2.0
// that can be found in the LICENSE file.

// Package falcon is a simple wrap implement for develop http server.
package falcon

import (
	"fmt"
	"github.com/raythorn/falcon/router"
	"log"
	"net/http"
)

var (
	falcon *Falcon
	Env    *Environment
)

func init() {
	falcon = New()
	Env = &Environment{data: make(map[string]string)}
}

type Falcon struct {
	router.Router
	namespace *router.NameSpace
}

func New() *Falcon {
	return &Falcon{router.New(), router.NewNameSpace()}
}

func (f *Falcon) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f.Handle(rw, req)
}

func (f *Falcon) run() {
	log.Println("Server listen at 192.168.1.107:8080")

	http.ListenAndServe("192.168.1.107:8080", f)
	if Env.TLS() {
		go func() {
			cert := Env.TLSCert()
			key := Env.TLSKey()
			port := Env.TLSPort()

			addr := fmt.Sprintf(":%d", port)
			http.ListenAndServeTLS(addr, cert, key, f)
		}()
	}

}

func Run() {
	falcon.run()
}

func Use(handler router.Handler) {
	falcon.Use(handler)
}

func Get(pattern string, handler router.Handler) {
	falcon.Get(pattern, handler)
}

func Patch(pattern string, handler router.Handler) {
	falcon.Patch(pattern, handler)
}

func Put(pattern string, handler router.Handler) {
	falcon.Put(pattern, handler)
}

func Post(pattern string, handler router.Handler) {
	falcon.Post(pattern, handler)
}

func Delete(pattern string, handler router.Handler) {
	falcon.Delete(pattern, handler)
}

func Head(pattern string, handler router.Handler) {
	falcon.Head(pattern, handler)
}

func Options(pattern string, handler router.Handler) {
	falcon.Options(pattern, handler)
}

func Any(pattern string, handler router.Handler) {
	falcon.Any(pattern, handler)
}

func NotFound(handler router.Handler) {
	falcon.NotFound(handler)
}

func NotAllowed(handler router.Handler) {
	falcon.NotAllowed(handler)
}

func NameSpace(prefix string, routes ...*router.Route) *router.NameSpace {

	ns := router.NewNameSpace()
	ns.NameSpace(prefix, routes...)

	falcon.Router.NameSpace(ns.Route)

	return ns
}

func NSGet(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Get(pattern, handler)
}

func NSPatch(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Patch(pattern, handler)
}

func NSPut(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Put(pattern, handler)
}

func NSPost(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Post(pattern, handler)
}

func NSDelete(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Delete(pattern, handler)
}

func NSHead(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Head(pattern, handler)
}

func NSOptions(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Options(pattern, handler)
}

func NSAny(pattern string, handler router.Handler) *router.Route {
	return falcon.namespace.Any(pattern, handler)
}
