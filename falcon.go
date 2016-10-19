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
	"time"
)

var (
	falcon *app
	Env    *Environment
)

func init() {
	falcon = &app{router.New(), &router.Group{}}
	Env = &Environment{data: make(map[string]string)}
}

type app struct {
	router.Router
	g *router.Group
}

func (f *app) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f.Handle(rw, req)
}

func (f *app) run() {

	finish := make(chan bool, 1)

	go func() {
		host := Env.Host()
		port := Env.Port()
		addr := fmt.Sprintf("%s:%d", host, port)

		log.Printf("Server listen at %s", addr)

		if err := http.ListenAndServe(addr, f); err != nil {
			log.Println("ListenAndServe fail")
			time.Sleep(100 * time.Microsecond)
			finish <- true
		}
	}()

	if Env.TLS() {
		go func() {
			cert := Env.TLSCert()
			key := Env.TLSKey()
			host := Env.TLSHost()
			port := Env.TLSPort()

			addr := fmt.Sprintf("%s:%d", host, port)
			if err := http.ListenAndServeTLS(addr, cert, key, f); err != nil {
				log.Println("ListenAndServeTLS fail")
				time.Sleep(100 * time.Microsecond)
				finish <- true
			}
		}()
	}

	<-finish
}

func Run() {
	falcon.run()
}

func Use(handler router.Midware) {
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

func Group(prefix string, routes ...interface{}) *router.Group {

	return falcon.Router.Group(prefix, routes...)
}

func GGet(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Get(pattern, handler)
}

func GPatch(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Patch(pattern, handler)
}

func GPut(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Put(pattern, handler)
}

func GPost(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Post(pattern, handler)
}

func GDelete(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Delete(pattern, handler)
}

func GHead(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Head(pattern, handler)
}

func GOptions(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Options(pattern, handler)
}

func GAny(pattern string, handler router.Handler) *router.Route {
	return falcon.g.Any(pattern, handler)
}
