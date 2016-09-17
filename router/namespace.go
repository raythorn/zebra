package router

import (
	"strings"
)

type NameSpace struct {
	Route *Route
}

func NewNameSpace() *NameSpace {
	return &NameSpace{Route: newRoute()}
}

// Before set handlers which will be called before actually http handler.
// All routes in this namespace will be affected if set
func (ns *NameSpace) Before(handlers ...Handler) *NameSpace {
	if ns.Route.before == nil {
		ns.Route.before = make([]Handler, 0)
	}

	ns.Route.before = append(ns.Route.before, handlers...)
	return ns
}

// Before set handlers which will be called before actually http handler.
// All routes in this namespace will be affected if set
func (ns *NameSpace) After(handlers ...Handler) *NameSpace {
	if ns.Route.after == nil {
		ns.Route.after = make([]Handler, 0)
	}

	ns.Route.after = append(ns.Route.after, handlers...)
	return ns
}

func (ns *NameSpace) NameSpace(prefix string, routes ...*Route) *Route {
	ns.Route.pattern = prefix
	ns.Route.cleanPath()
	ns.Route.regexpCompile()

	if strings.HasSuffix(ns.Route.pattern, `\/?`) {
		panic("NameSpace cannot use regexp!!!")
	}

	for _, r := range routes {
		r.pattern = ns.Route.pattern + r.pattern
		r.regexpCompile()
		insertRoute(ns.Route, r)
	}

	return ns.Route
}

func (ns *NameSpace) Get(pattern string, handler Handler) *Route {
	return ns.add("GET", pattern, handler)
}

func (ns *NameSpace) Patch(pattern string, handler Handler) *Route {
	return ns.add("PATCH", pattern, handler)
}

func (ns *NameSpace) Put(pattern string, handler Handler) *Route {
	return ns.add("PUT", pattern, handler)
}

func (ns *NameSpace) Post(pattern string, handler Handler) *Route {
	return ns.add("POST", pattern, handler)
}

func (ns *NameSpace) Delete(pattern string, handler Handler) *Route {
	return ns.add("DELETE", pattern, handler)
}

func (ns *NameSpace) Head(pattern string, handler Handler) *Route {
	return ns.add("HEAD", pattern, handler)
}

func (ns *NameSpace) Options(pattern string, handler Handler) *Route {
	return ns.add("OPTIONS", pattern, handler)
}

func (ns *NameSpace) Any(pattern string, handler Handler) *Route {

	return ns.add("ANY", pattern, handler)
}

func (ns *NameSpace) add(method, pattern string, handler Handler) *Route {
	route := newRoute()
	route.pattern = pattern
	route.actions[method] = handler
	route.cleanPath()
	route.regexpCompile()

	return route
}
