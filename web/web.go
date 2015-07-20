package web

import (
	"github.com/codegangsta/negroni"
	negronilogrus "github.com/meatballhat/negroni-logrus" // package name doesn't match! >:(
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tylerb/graceful"
	// "gopkg.in/unrolled/secure.v1" // TODO!
	"github.com/asteris-llc/tellus/tf"
	"github.com/gorilla/mux"
	"time"
)

func Serve(addr string, tf *tf.Terraformer) {
	router := mux.NewRouter()

	state := router.PathPrefix(`/state/{project}`).Subrouter()
	stateHandler := StateHandler{tf}
	state.Methods("GET").HandlerFunc(stateHandler.GetState)
	state.Methods("POST", "PUT").HandlerFunc(stateHandler.SetState)
	state.Methods("DELETE").HandlerFunc(stateHandler.DeleteState)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	graceful.Run(addr, 10*time.Second, n)
}
