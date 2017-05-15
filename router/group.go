package router

import (
	"github.com/raythorn/falcon/context"
	"strings"
)

type Group struct {
	pattern string
	routes  map[string]*Route
	groups  map[string]*Group
	before  []Midware
	after   []Midware
}

func newGroup() *Group {
	return &Group{
		pattern: "",
		routes:  make(map[string]*Route),
		groups:  make(map[string]*Group),
		before:  nil,
		after:   nil,
	}
}

// Before set midwares which will be called before actually http handler.
// All routes in this group will be affected if set
func (g *Group) Before(midwares ...Midware) *Group {
	if g.before == nil {
		g.before = make([]Midware, 0)
	}

	g.before = append(g.before, midwares...)
	return g
}

// After set midwares which will be called after actually http handler.
// All routes in this group will be affected if set
func (g *Group) After(midwares ...Midware) *Group {
	if g.after == nil {
		g.after = make([]Midware, 0)
	}

	g.after = append(g.after, midwares...)
	return g
}

func (g *Group) Sub(prefix string, args ...interface{}) *Group {
	group := newGroup()
	return group.group(prefix, args...)
}

func (g *Group) Get(pattern string, handler Handler) *Route {
	return g.add("GET", pattern, handler)
}

func (g *Group) Patch(pattern string, handler Handler) *Route {
	return g.add("PATCH", pattern, handler)
}

func (g *Group) Put(pattern string, handler Handler) *Route {
	return g.add("PUT", pattern, handler)
}

func (g *Group) Post(pattern string, handler Handler) *Route {
	return g.add("POST", pattern, handler)
}

func (g *Group) Delete(pattern string, handler Handler) *Route {
	return g.add("DELETE", pattern, handler)
}

func (g *Group) Head(pattern string, handler Handler) *Route {
	return g.add("HEAD", pattern, handler)
}

func (g *Group) Options(pattern string, handler Handler) *Route {
	return g.add("OPTIONS", pattern, handler)
}

func (g *Group) Any(pattern string, handler Handler) *Route {

	return g.add("ANY", pattern, handler)
}

func (g *Group) group(pattern string, args ...interface{}) *Group {

	for _, arg := range args {
		switch arg.(type) {
		case *Route:
			route, _ := arg.(*Route)
			route.pattern = cleanPath(pattern + route.pattern)
			route.regexpCompile()
			route.group = g

			if r, ok := g.routes[route.pattern]; ok {
				for m, h := range route.actions {
					r.actions[m] = h
				}

				route = nil
			} else {
				g.routes[route.pattern] = route
			}
		case *Group:
			grp, _ := arg.(*Group)
			grp.pattern = cleanPath(pattern + grp.pattern)
			g.groups[grp.pattern] = grp

			if len(grp.routes) > 0 {
				for _, route := range grp.routes {
					route.pattern = cleanPath(pattern + route.pattern)
					route.regexpCompile()
					g.routes[route.pattern] = route
				}
			}

			if len(grp.groups) > 0 {
				for _, group := range grp.groups {
					group.pattern = cleanPath(pattern + group.pattern)
					g.groups[group.pattern] = group
					if len(grp.before) > 0 {
						if group.before == nil {
							group.before = make([]Midware, 0)
						}
						group.before = append(grp.before, group.before...)
					}
					if len(grp.after) > 0 {
						if group.after == nil {
							group.after = make([]Midware, 0)
						}

						group.after = append(group.after, grp.after...)
					}
				}
			}
		}
	}

	return g
}

func (g *Group) add(method, pattern string, handler Handler) *Route {
	route := newRoute()
	route.pattern = cleanPath(pattern)
	route.actions[method] = handler
	route.regexpCompile()

	return route
}

func (g *Group) insert(method, pattern string, handler Handler) *Route {

	route := newRoute()

	route.pattern = cleanPath(pattern)
	route.actions[method] = handler
	route.regexpCompile()

	if rt, ok := g.routes[route.pattern]; ok {
		for m, h := range route.actions {
			rt.actions[m] = h
		}

		if route.oss != nil {
			rt.oss = route.oss
		}

		route = nil
		return rt
	} else {
		g.routes[route.pattern] = route
		return route
	}
}

func (g *Group) match(ctx *context.Context) *Route {

	if r, ok := g.routes[ctx.URL()]; ok {
		if r.match(ctx) {
			return r
		}
	} else {

		for p, r := range g.routes {
			if strings.Contains(p, "(?P") {
				if r.match(ctx) {
					return r
				}
			}
		}
	}

	return nil
}
