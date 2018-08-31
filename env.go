package zebra

import (
	"fmt"
	"strconv"
	"sync"
)

//Simple implement for zebra evironment, you can use this to save your configs.
//It's thread-safe, so you can have all your fun in erverywhere.
type Environment struct {
	sync.RWMutex
	data map[string]string
}

// Set env variable with a pair of key-value, you cann't use key with prefix "Zebra:", which
// is reserved for zebra system
func (e *Environment) Set(key, value string) {
	e.Lock()
	defer e.Unlock()

	e.data[key] = value
}

//Get a env variable with key, "" will return if not exist
func (e *Environment) Get(key string) string {
	e.RLock()
	defer e.RUnlock()

	if value, ok := e.data[key]; ok {
		return value
	}
	return ""
}

//Delete a env variable with key
func (e *Environment) Del(key string) {
	e.Lock()
	defer e.Unlock()

	delete(e.data, key)
}

//Enable HTTPS to create a secure channel for private data, you must provide
//valid cert file, key file and HTTPS port
func (e *Environment) EnableTLS(cert, key, host string, port int) {

	e.Set("Zebra:TLS", "true")
	e.Set("Zebra:TLSCERT", cert)
	e.Set("Zebra:TLSKEY", key)
	e.Set("Zebra:TLSHOST", host)
	p := fmt.Sprintf("%d", port)
	e.Set("Zebra:TLSPORT", p)
}

//Disable HTTPS
func (e *Environment) DisableTLS() {
	e.Set("Zebra:TLS", "false")
}

//Check if HTTPS is enabled
func (e *Environment) TLS() bool {

	tls := e.Get("Zebra:TLS")
	if tls == "true" {
		return true
	}

	return false
}

//Get https cert file
func (e *Environment) TLSCert() string {
	return e.Get("Zebra:TLSCERT")
}

//Get https key file
func (e *Environment) TLSKey() string {
	return e.Get("Zebra:TLSKEY")
}

//Get https host
func (e *Environment) TLSHost() string {
	return e.Get("Zebra:TLSHOST")
}

//Get https port, if port not set or error happened, default port 443 will be returned
func (e *Environment) TLSPort() int {
	port := e.Get("Zebra:TLSPORT")
	if port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			return p
		}
	}

	return 443
}

//Set http host and port, if not set, "localhost:8080" will be used
func (e *Environment) Http(host string, port int) {

	p := fmt.Sprintf("%d", port)

	e.Set("Zebra:HOST", host)
	e.Set("Zebra:PORT", p)
}

//Get host, works for both http and https
func (e *Environment) Host() string {
	return e.Get("Zebra:HOST")
}

//Get http listen port
func (e *Environment) Port() int {
	port := e.Get("Zebra:PORT")
	if port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			return p
		}
	}

	return 8080
}
