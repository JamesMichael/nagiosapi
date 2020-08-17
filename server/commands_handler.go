package server

import (
	"encoding/json"
	"net/http"

	"github.com/jamesmichael/nagiosapi/service/submission"
)

type CommandsHandler struct {
	svc *submission.Service
}

func NewCommandsHandler(svc *submission.Service) (*CommandsHandler, error) {
	return &CommandsHandler{
		svc: svc,
	}, nil
}

func (h *CommandsHandler) SubmitResult(w http.ResponseWriter, r *http.Request) {
	var res submission.ServiceResult
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	if err := h.svc.SubmitPassiveResult(&res); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`"ok"`))
}
