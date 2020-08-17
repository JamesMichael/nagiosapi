package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jamesmichael/nagiosapi/encoding/xdata"
	"github.com/jamesmichael/nagiosapi/nagios/statusdata"
)

type StatusService interface {
	ServiceStatus(host, name string) (*xdata.ServiceStatus, error)
}

// RegisterStatusService sets up /status and /status/HOST/SERVICE routes
// for accessing service statuses.
func (s *Server) RegisterStatusService(svc StatusService) {
	s.mux.Route("/status", func(r chi.Router) {
		r.Get("/{host}/{service}", handleServiceStatus(svc))
		r.Post("/", handleMultiServiceStatus(svc))
	})
}

type serviceStatusResponse struct {
	IsFound  bool   `json:"is_found"`
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
	Output   string `json:"output"`
	Status   string `json:"status"`
}

func handleServiceStatus(svc StatusService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host := chi.URLParam(r, "host")
		service := chi.URLParam(r, "service")

		st, err := svc.ServiceStatus(host, service)
		if err != nil {
			if errors.Is(err, statusdata.ErrUnknownHost) || errors.Is(err, statusdata.ErrUnknownService) {
				http.Error(w, http.StatusText(404), 404)
				return
			}

			http.Error(w, http.StatusText(500), 500)
			return
		}

		out, err := json.Marshal(serviceStatusResponse{
			IsFound:  true,
			Hostname: host,
			Service:  service,
			Status:   st.CurrentState.String(),
			Output:   st.PluginOutput,
		})
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write(out)
	}
}

type serviceStatusMultiRequest []struct {
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
}

func handleMultiServiceStatus(svc StatusService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req serviceStatusMultiRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			panic(err)
		}

		res := make([]serviceStatusResponse, 0, len(req))
		for _, service := range req {
			host := service.Hostname
			service := service.Service

			s, err := svc.ServiceStatus(host, service)
			if err != nil {
				if errors.Is(err, statusdata.ErrUnknownHost) || errors.Is(err, statusdata.ErrUnknownService) {
					res = append(res, serviceStatusResponse{
						IsFound:  false,
						Hostname: host,
						Service:  service,
					})
				}
				continue
			}

			res = append(res, serviceStatusResponse{
				IsFound:  true,
				Hostname: host,
				Service:  service,
				Status:   s.CurrentState.String(),
				Output:   s.PluginOutput,
			})
		}

		out, err := json.Marshal(res)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		w.Write(out)
	}
}
