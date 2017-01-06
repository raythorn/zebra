package router

import (
	"fmt"
	"github.com/raythorn/falcon/context"
	// "github.com/raythorn/falcon/log"
	"regexp"
	"strings"
)

type Route struct {
	pattern string
	regexp  *regexp.Regexp
	actions map[string]Handler
	group   *Group
}

func newRoute() *Route {
	return &Route{"", nil, make(map[string]Handler), nil}
}

func (r *Route) match(ctx *context.Context) bool {

	if _, ok := r.actions[ctx.Method()]; !ok {
		if _, ok := r.actions["ANY"]; !ok {
			return false
		}
	}

	if ctx.URL() == r.pattern {
		return true
	}

	matches := r.regexp.FindStringSubmatch(ctx.URL())

	if len(matches) > 0 && matches[0] == ctx.URL() {
		for i, name := range r.regexp.SubexpNames() {
			// log.Println(name)
			if len(name) > 0 {
				ctx.Set(name, matches[i])
			}
		}
		return true
	}

	return false
}

func (r *Route) regexpCompile() {
	routeExp := regexp.MustCompile(`:[^/#?()\.\\]+`)
	r.pattern = routeExp.ReplaceAllStringFunc(r.pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})

	pattern := r.pattern
	if !strings.HasSuffix(pattern, `\/?`) {
		pattern += `\/?`
	}

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

func cleanPath(pattern string) string {

	path := pattern
	if path == "" {
		return "/"
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
		return string(buf[:write])
	} else {
		return string(pattern[:write])
	}
}
