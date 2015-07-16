package web

import (
	"net/http"
)

func GetState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func SetState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(http.StatusText(http.StatusAccepted)))
}

func DeleteState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
