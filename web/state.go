package web

import (
	"encoding/json"
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

func (s *StateHandler) GetState(w http.ResponseWriter, r *http.Request) {
	state, err := s.state.GetState(s.Project(r))
	if err == tf.ErrNoState {
		// Terraform doesn't like 404 responses
		state = terraform.NewState()
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

func (s *StateHandler) SetState(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *StateHandler) DeleteState(w http.ResponseWriter, r *http.Request) {
	err := s.state.DeleteState(s.Project(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
