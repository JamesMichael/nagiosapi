package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jamesmichael/nagiosapi/encoding/xdata"
	"github.com/jamesmichael/nagiosapi/nagios/statusdata"
)

type ServiceStatusProvider interface {
	ServiceStatus(host, name string) (*xdata.ServiceStatus, error)
}

// StatusHandler implements functionality for /api/status.
type StatusHandler struct {
	ssp ServiceStatusProvider
}

func NewStatusHandler(ssp ServiceStatusProvider) (*StatusHandler, error) {
	sh := &StatusHandler{
		ssp: ssp,
	}
	return sh, nil
}

type serviceStatusResponse struct {
	IsFound  bool   `json:"is_found"`
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
	Output   string `json:"output"`
	Status   string `json:"status"`
}

// GetServiceStatus returns status information for a single service.
func (sh *StatusHandler) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	host := chi.URLParam(r, "host")
	service := chi.URLParam(r, "service")

	s, err := sh.ssp.ServiceStatus(host, service)
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
		Status:   s.CurrentState.String(),
		Output:   s.PluginOutput,
	})
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(out)
}

type serviceStatusMultiRequest []struct {
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
}

// GetServiceStatusMulti returns status information for multiple services.
func (sh *StatusHandler) GetServiceStatusMulti(w http.ResponseWriter, r *http.Request) {
	var req serviceStatusMultiRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	res := make([]serviceStatusResponse, 0, len(req))
	for _, svc := range req {
		host := svc.Hostname
		service := svc.Service

		s, err := sh.ssp.ServiceStatus(host, service)
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
