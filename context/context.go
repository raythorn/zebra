package context

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	acceptsHTMLRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXMLRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
)

type Context struct {
	rw      http.ResponseWriter
	request *http.Request
	data    map[string]string
}

// Return a new Context instance
func New() *Context {
	return &Context{data: make(map[string]string)}
}

// Initialise Context with HTTP Request and ResponseWriter, it will parse the Request header,
// and it also parse the get/post/put form parameters. NOTE: The Path Regexp param MUST NOT have
// same name with HTTP Request form param, otherwise, it will override the HTTP form param
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.request = r
	c.rw = w

	// Parse Request Header
	for k, v := range c.request.Header {
		c.Set(k, v)
	}

	// Parse Request Form
	c.request.ParseForm()
	for k, v := range c.request.Form {
		c.Set(k, strings.Join(v, ""))
	}
}

// Get data from context
func (c *Context) Get(key interface{}) interface{} {
	if v, ok := c.data[key]; ok {
		return v
	}

	return nil
}

// Set data to context
func (c *Context) Set(key, value interface{}) {
	if c.data == nil {
		c.data = make(map[interface{}]interface{})
	}

	c.data[key] = value
}

//Request relate method

// Protocol returns request protocol name, such as HTTP/1.1 .
func (c *Context) Protocol() string {
	return c.request.Proto
}

// URI returns full request url with query string, fragment.
func (c *Context) URI() string {
	return c.request.RequestURI
}

// URL returns request url path (without query string, fragment).
func (c *Context) URL() string {
	return c.request.URL.Path
}

// Scheme returns Request scheme, "http" or "https"
func (c *Context) Scheme() string {
	if scheme := c.request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}

	if c.request.URL.Scheme != "" {
		return c.request.URL.Scheme
	}

	if c.request.TLS == nil {
		return "http"
	}

	return "https"
}

// Host returns request host name, if no host info in requst, "localhost" will be returned
func (c *Context) Host() string {
	if c.request.Host != "" {
		hostParts := strings.Split(c.request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return c.request.Host
	}
	return "localhost"
}

// Port returns request port, 80 will be returned if error happens
func (c *Context) Port() int {
	parts := strings.Split(c.request.Host, ":")
	if len(parts) == 2 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}

// Site returns site url as scheme://domain
func (c *Context) Site() string {
	return c.Scheme() + "://" + c.Host()
}

// Domain returns host name, alias of Host
func (c *Context) Domain() string {
	return c.Host()
}

// SubDomain returns sub domain string, like api.raythorn.com will return api
func (c *Context) SubDomain() string {
	parts := strings.Split(c.Host(), ".")
	if len(parts) >= 3 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

// Method returns requst method
func (c *Context) Method() string {
	return c.request.Method
}

// RemoteAddr returns reomte address in request
func (c *Context) RemoteAddr() string {
	return c.request.RemoteAddr
}

// Referer returns request referer
func (c *Context) Referer() string {
	return c.request.Referer()
}

// UserAgent returns client agent
func (c *Context) UserAgent() string {
	return c.request.Header.Get("User-Agent")
}

// Proxy returns proxy client ips slice.
func (c *Context) Proxy() []string {
	if ips := c.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}

	return []string{}
}

// IP returns request client ip.
// if in proxy, return first proxy id.
// if error, return 127.0.0.1.
func (c *Context) Ip() string {
	ips := c.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		ip := strings.Split(ips[0], ":")
		return ip[0]
	}

	ip := strings.Split(c.RemoteAddr(), ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}

	return "127.0.01"
}

// AcceptsHTML Checks if request accepts html response
func (c *Context) AcceptsHTML() bool {
	return acceptsHTMLRegex.MatchString(c.Get("Accept"))
}

// AcceptsXML Checks if request accepts xml response
func (c *Context) AcceptsXML() bool {
	return acceptsXMLRegex.MatchString(c.Get("Accept"))
}

// AcceptsJSON Checks if request accepts json response
func (c *Context) AcceptsJSON() bool {
	return acceptsJSONRegex.MatchString(c.Get("Accept"))
}

//ResponseWriter relate method

func (c *Context) Header(key, value string) {
	c.rw.Header().Set(key, value)
}

func (c *Context) WriteHeader(code int) {
	c.rw.WriteHeader(code)
}

func (c *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijack, ok := c.rw.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("Web server doesn't support Hijack!")
	}

	return hijack.Hijack()
}

func (c *Context) Flush() {
	if f, ok := c.rw.(http.Flusher); ok {
		f.Flush()
	}
}

func (c *Context) CloseNotify() <-chan bool {
	if cn, ok := c.rw.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}

	return nil
}

func (c *Context) WriteString(data string) error {
	_, err := c.rw.Write([]byte(data))

	return err
}

func (c *Context) Write(bytes []byte) (int, error) {
	return c.rw.Write(bytes)
}

func (c *Context) ServeJSON() {

}

func (c *Context) ServeXML() {

}

func (c *Context) ServeJSONP() {

}
