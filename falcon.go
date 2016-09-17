package falcon

import (
	"github.com/raythorn/falcon/router"
	"log"
	"net/http"
)

var (
	falcon *Falcon
)

func init() {
	falcon = New()
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
