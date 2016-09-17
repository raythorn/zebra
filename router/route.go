package router

import (
	"fmt"
	"github.com/raythorn/falcon/context"
	"regexp"
	"strings"
)

type Route struct {
	actions map[string]Handler
	pattern string
	regexp  *regexp.Regexp
	before  []Handler
	after   []Handler
	routes  map[string]*Route
}

func newRoute() *Route {
	return &Route{make(map[string]Handler), "", nil, nil, nil, make(map[string]*Route)}
}

func (r *Route) match(method, path string, handlers map[string][]Handler) *Route {

	//Only namespace route can use this operation
	if handlers != nil {
		if r.before != nil && len(r.before) > 0 {

			if h, ok := handlers["before"]; !ok {
				h = make([]Handler, 0)
				h = append(h, r.before...)
				handlers["before"] = h
			} else {
				h = append(h, r.before...)
			}
		}

		if r.after != nil && len(r.after) > 0 {

			if h, ok := handlers["after"]; !ok {
				h = make([]Handler, 0)
				h = append(h, r.after...)
				handlers["after"] = h
			} else {
				h = append(h, r.after...)
			}
		}
	}

	if _, ok := r.actions[method]; ok {

		if path == r.pattern {
			return r
		}

		matches := r.regexp.FindStringSubmatch(path)
		if len(matches) > 0 && matches[0] == path {

			return r
		}
	}

	//Handle sub-routes
	if len(r.routes) > 0 {
		for _, route := range r.routes {
			subroute := route.match(method, path, handlers)
			if subroute != nil {
				return subroute
			}
		}
	}

	return nil
}

func (r *Route) params(ctx *context.Context) {
	matchs := r.regexp.FindStringSubmatch(ctx.URI())
	for i, name := range r.regexp.SubexpNames() {
		if len(name) > 0 {
			ctx.Set(name, matchs[i])
		}
	}
}

func (r *Route) regexpCompile() {
	routeExp := regexp.MustCompile(`:[^/#?()\.\\]+`)
	r.pattern = routeExp.ReplaceAllStringFunc(r.pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})

	pattern := r.pattern + `\/?`
	r.regexp = regexp.MustCompile(pattern)
	if strings.Contains(r.pattern, "(?P") {
		r.pattern = pattern
	}
}

// cleanPath is the URL version of path.Clean, it returns a canonical URL path
// for example, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//	5. Omit regexp characters in the path, all regexp will be in a pair of "()"
//	6. Eliminate the trailing slash
//
// If the result of this process is an empty string, "/" is returned

func (r *Route) cleanPath() {

	path := r.pattern
	if path == "" {
		r.pattern = "/"
		return
	}

	read := 1
	write := 1

	size := len(path)
	var buf []byte = nil

	if path[0] != '/' {
		buf = make([]byte, size+1)
		// Must start with a single slash
		buf[0] = '/'
		read = 0
	}

	for read < size {
		switch {
		case path[read] == '/':
			//Eliminate trailing slash and multiple slash
			read++
		case path[read] == '.' && (read+1 == size || path[read+1] == '/'):
			//Eliminate trailing '.'
			read++
		case path[read] == '.' && path[read+1] == '.' && (read+2 == size || path[read+2] == '/'):
			read += 2

			if buf == nil {
				for write > 1 && path[write] != '/' {
					write--
				}
			} else {
				for write > 1 && buf[write] != '/' {
					write--
				}
			}

		default:

			if buf == nil && read != write {
				buf = make([]byte, size+1)
				copy(buf, path[:write])
			}

			if write > 1 && buf != nil {
				buf[write] = '/'
				write++
			}

			wildcard := path[read] == '('

			for read < size && ((wildcard && path[read] != ')') || path[read] != '/') {
				if buf != nil {

					buf[write] = path[read]
				}

				read++
				write++
			}

		}
	}

	if buf != nil {
		r.pattern = string(buf[:write])
	} else {
		r.pattern = string(r.pattern[:write])
	}
}
