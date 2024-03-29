package web

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/tellus/tf"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
)

type StateHandler struct {
	state tf.StateManipulator
}

func (s *StateHandler) Project(r *http.Request) string {
	return mux.Vars(r)["project"]
}

func (s *StateHandler) Get(w http.ResponseWriter, r *http.Request) {
	state, err := s.state.GetState(s.Project(r))
	switch {
	case err == tf.ErrNoState:
		w.WriteHeader(http.StatusNoContent)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error. Check logs and try again later."))
		logrus.WithField("error", err.Error()).Error("error getting state")
		return
	}

	body, err := json.Marshal(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func (s *StateHandler) Set(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	state, err := terraform.ReadState(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = s.state.SetState(s.Project(r), state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *StateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := s.state.DeleteState(s.Project(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
