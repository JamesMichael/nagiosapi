package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jamesmichael/nagiosapi/nagios/cmd"
)

type CommandsHandler struct {
	w io.Writer
}

func NewCommandsHandler(writer io.Writer) (*CommandsHandler, error) {
	return &CommandsHandler{
		w: writer,
	}, nil
}

type SubmitResultRequest struct {
	Time        int64  `json:"time"`
	Hostname    string `json:"hostname"`
	ServiceName string `json:"service_name"`
	Status      uint   `json:"status"`
	Body        string `json:"body"`
}

func (h *CommandsHandler) SubmitResult(w http.ResponseWriter, r *http.Request) {
	var req SubmitResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer r.Body.Close()

	if req.Time == 0 {
		req.Time = time.Now().Unix()
	}

	command := fmt.Sprintf("[%d] PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s",
		req.Time,
		cmd.Sanitize(req.Hostname),
		cmd.Sanitize(req.ServiceName),
		req.Status,
		cmd.Sanitize(req.Body),
	)
	if _, err := h.w.Write([]byte(command)); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`"ok"`))
}
