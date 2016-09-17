package router

import (
	"github.com/raythorn/falcon/context"
	"log"
	"net/http"
	"strings"
)

type Handler func(*context.Context)

type Router interface {

	// Add middleware to router, these handler will called before every request handler
	Use(Handler)

	NameSpace(*Route)

	// Get adds a route for a HTTP GET request to the specified matching pattern.
	Get(string, Handler)

	// Patch adds a route for a HTTP PATCH request to the specified matching pattern.
	Patch(string, Handler)

	// Put adds a route for a HTTP PUT request to the specified matching pattern.
	Put(string, Handler)

	// Post adds a route for a HTTP POST request to the specified matching pattern.
	Post(string, Handler)

	// Delete adds a route for a HTTP DELETE request to the specified matching pattern.
	Delete(string, Handler)

	// Head adds a route for a HTTP HEAD request to the specified matching pattern.
	Head(string, Handler)

	// Options adds a route for a HTTP OPTIONS request to the specified matching pattern.
	Options(string, Handler)

	// Any adds a route for any HTTP method request to the specified matching pattern.
	Any(string, Handler)

	// NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.
	NotFound(Handler)

	// NotAllowed sets the handler that are called when a not allowed http method request
	NotAllowed(Handler)

	// Handle is the entry point for routing.
	Handle(http.ResponseWriter, *http.Request)
}

type router struct {
	route      *Route
	namespaces map[string]*Route
	handlers   []Handler
}

func New() Router {

	r := &router{route: newRoute(), namespaces: make(map[string]*Route), handlers: make([]Handler, 0)}
	r.route.pattern = "/"
	r.route.cleanPath()
	r.route.regexpCompile()

	return r
}

func (r *router) Use(handler Handler) {
	r.handlers = append(r.handlers, handler)
}

func (r *router) NameSpace(route *Route) {
	if _, ok := r.namespaces[route.pattern]; ok {
		panic("NameSpace exist!")
	}

	for prefix, _ := range r.namespaces {
		if strings.HasPrefix(route.pattern, prefix) || strings.HasPrefix(prefix, route.pattern) {
			panic("NameSpace has sub namespace or parent namespace, please unity them together!")
		}
	}

	r.namespaces[route.pattern] = route
}

func (r *router) Get(pattern string, handler Handler) {

	r.add("GET", pattern, handler)
}

func (r *router) Patch(pattern string, handler Handler) {
	r.add("PATCH", pattern, handler)
}

func (r *router) Put(pattern string, handler Handler) {
	r.add("PUT", pattern, handler)
}

func (r *router) Post(pattern string, handler Handler) {
	r.add("POST", pattern, handler)
}

func (r *router) Delete(pattern string, handler Handler) {
	r.add("DELETE", pattern, handler)
}

func (r *router) Head(pattern string, handler Handler) {
	r.add("HEAD", pattern, handler)
}

func (r *router) Options(pattern string, handler Handler) {
	r.add("OPTIONS", pattern, handler)
}

func (r *router) Any(pattern string, handler Handler) {
	r.add("ANY", pattern, handler)
}

func (r *router) NotFound(handler Handler) {

}

func (r *router) NotAllowed(handler Handler) {

}

func (r *router) Handle(rw http.ResponseWriter, req *http.Request) {

	ctx := context.New()
	ctx.Reset(rw, req)

	log.Printf("URI: %s", ctx.URI())
	log.Printf("PATH: %s", ctx.URL())

	//Call all middleware first
	if len(r.handlers) > 0 {
		for _, h := range r.handlers {
			h(ctx)
		}
	}

	//Match namespace first
	if len(r.namespaces) > 0 {
		log.Println("Namespace")
		for pattern, route := range r.namespaces {

			if strings.HasPrefix(ctx.URL(), pattern) {
				handlers := map[string][]Handler{}
				matchroute := route.match(ctx.Method(), ctx.URL(), handlers)

				//Found route
				if matchroute != nil {
					matchroute.params(ctx)
					if handler, ok := matchroute.actions[ctx.Method()]; ok {
						if bfs, ok := handlers["before"]; ok && len(bfs) > 0 {
							for _, bf := range bfs {
								bf(ctx)
							}
						}

						handler(ctx)

						if afs, ok := handlers["after"]; ok && len(afs) > 0 {
							for _, af := range afs {
								af(ctx)
							}
						}
					}

					return
				}
			}
		}
	}

	log.Println("Match route")
	matchroute := r.route.match(ctx.Method(), ctx.URL(), nil)
	if matchroute != nil {
		log.Println("Found route")
		matchroute.params(ctx)
		if handler, ok := matchroute.actions[ctx.Method()]; ok {
			log.Printf("Excute %s", ctx.Method())
			handler(ctx)
		} else {
			log.Println("Panic")
		}
	}

}

func (r *router) add(method, pattern string, handler Handler) *Route {
	route := newRoute()

	route.pattern = pattern
	route.actions[method] = handler
	route.cleanPath()
	route.regexpCompile()

	insertRoute(r.route, route)

	return route
}

func insertRoute(root, route *Route) bool {

	// Check if route is subroute of root
	if ok := strings.HasPrefix(route.pattern, root.pattern); !ok {
		log.Printf("The route to insert is not a sub-route of root!")
		return false
	}

	//Same route with different actions
	if root.pattern == route.pattern {
		for m, h := range route.actions {
			if _, ok := root.actions[m]; ok {
				log.Printf("Method exist for %s, ignored", route.pattern)
			} else {
				root.actions[m] = h
			}
		}

		if route.before != nil && len(route.before) > 0 {
			if root.before == nil {
				root.before = make([]Handler, 0)
			}
			root.before = append(root.before, route.before...)
		}

		if route.after != nil && len(route.after) > 0 {
			if root.after == nil {
				root.after = make([]Handler, 0)
			}
			root.after = append(root.after, route.after...)
		}
	} else { //Sub-route
		log.Println("Sub route")
		if r, ok := root.routes[route.pattern]; ok {
			if ok := insertRoute(r, route); !ok {
				log.Printf("Insert route(%s) to route(%S) failed", route.pattern, root.pattern)
				return false
			}
		} else {
			root.routes[route.pattern] = route
		}
	}
	return true
}
