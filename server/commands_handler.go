package server

import (
	"encoding/json"
	"net/http"
)

type PassiveCommandService interface {
	SubmitResult(time int64, statusCode uint, hostname, serviceName, body string) error
}

type passiveServiceResult struct {
	Time        int64  `json:"time"`
	Hostname    string `json:"hostname"`
	ServiceName string `json:"service_name"`
	Status      uint   `json:"status"`
	Body        string `json:"body"`
}

// RegisterPassiveCommandService sets up /submit route for submitting
// passive service check results.
func (s *Server) RegisterPassiveCommandService(svc PassiveCommandService) {
	s.mux.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		var res passiveServiceResult
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		defer r.Body.Close()

		if err := svc.SubmitResult(
			res.Time,
			res.Status,
			res.Hostname,
			res.ServiceName,
			res.Body,
		); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`"ok"`))
	})
}
