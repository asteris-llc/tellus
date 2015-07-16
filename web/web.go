package web

import (
	"github.com/codegangsta/negroni"
	negronilogrus "github.com/meatballhat/negroni-logrus" // package name doesn't match! >:(
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tylerb/graceful"
	// "gopkg.in/unrolled/secure.v1" // TODO!
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"time"
)

func Serve(addr string) {
	router := mux.NewRouter()
	state := router.PathPrefix(`/state/{projects}`).Subrouter()
	state.Methods("GET").HandlerFunc(GetState)
	state.Methods("POST", "PUT").HandlerFunc(SetState)
	state.Methods("DELETE").HandlerFunc(DeleteState)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	logrus.WithField("addr", addr).Info("started listening")
	graceful.Run(addr, 10*time.Second, n)
}
