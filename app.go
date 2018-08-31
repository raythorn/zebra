package zebra

import (
	"fmt"
	"github.com/raythorn/zebra/log"
	"github.com/raythorn/zebra/router"
	"net/http"
	// "os"
	// "os/exec"
	// "path/filepath"
	"time"
)

type app struct {
	router.Router
	g *router.Group
}

func (a *app) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.Handle(rw, req)
}

func (a *app) run() {

	finish := make(chan bool, 1)

	go func() {
		host := Env.Host()
		port := Env.Port()
		addr := fmt.Sprintf("%s:%d", host, port)

		log.Info("Server listen at %s", addr)

		if err := http.ListenAndServe(addr, a); err != nil {
			log.Error("ListenAndServe fail")
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
			if err := http.ListenAndServeTLS(addr, cert, key, a); err != nil {
				log.Error("ListenAndServeTLS fail")
				time.Sleep(100 * time.Microsecond)
				finish <- true
			}
		}()
	}

	<-finish
}
