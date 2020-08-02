package statusdata

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jamesmichael/nagiosapi/encoding/xdata"
	"go.uber.org/zap"
)

var (
	// ErrUnknownHost is used to indicate that no service status can be found for the given host.
	ErrUnknownHost = errors.New("unknown host")

	// ErrUnknownService is used to indicate that no service status can be found for the given service.
	ErrUnknownService = errors.New("unknown service")
)

// Repository provides access to the data in Nagios' status.dat file.
//
// The status.dat is periodically reloaded to ensure the data remains fresh.
type Repository struct {
	filename        string
	mux             sync.RWMutex
	refreshInterval time.Duration
	log             *zap.Logger

	services map[string]map[string]*xdata.ServiceStatus
}

// NewRepository constructs an instance of statusdata.Repository.
//
// The filepath of the Nagios status.dat file is required.
//
// The contents of the status file is loaded into memory to cache the results.  An error is returned if it is not
// possible to read the status file.
//
// The status file is periodically reloaded to ensure the statusdata is relatively fresh. If an error occurs during reload,
// the error is logged.
func NewRepository(filename string, opts ...RepositoryOpt) (*Repository, error) {
	r := &Repository{
		filename:        filename,
		refreshInterval: time.Minute,
		log:             zap.NewNop(),
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	if err := r.load(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repository) load() error {
	f, err := os.Open(r.filename)
	if err != nil {
		r.log.Error("unable to open nagios status file",
			zap.String("filename", r.filename),
			zap.Error(err),
		)
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			r.log.Error("unable to close nagios status file",
				zap.String("filename", r.filename),
				zap.Error(err),
			)
		}
	}()

	var raw xdata.Status
	if err := xdata.NewDecoder(f).Decode(&raw); err != nil {
		r.log.Error("unable to decode nagios status file",
			zap.String("filename", r.filename),
			zap.Error(err),
		)
		return fmt.Errorf("unable to decode nagios status file")
	}

	statuses := make(map[string]map[string]*xdata.ServiceStatus)
	for _, check := range raw.ServiceStatus {
		services, ok := statuses[check.HostName]
		if !ok {
			services = make(map[string]*xdata.ServiceStatus)
			statuses[check.HostName] = services
		}

		services[check.ServiceDescription] = check
	}

	r.mux.Lock()
	defer r.mux.Unlock()
	r.services = statuses

	r.log.Info("loaded nagios status file",
		zap.String("filename", r.filename),
	)

	return nil
}

// ServiceStatus looks up a Nagios service check result by host and service description.
//
// ErrUnknownHost and ErrUnknownService are returned if the respective Host and Service are not found in the Nagios
// statusdata file.
func (r *Repository) ServiceStatus(host, name string) (*xdata.ServiceStatus, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	services, ok := r.services[host]
	if !ok {
		return nil, ErrUnknownHost
	}

	service, ok := services[name]
	if !ok {
		return nil, ErrUnknownService
	}

	return service, nil
}

// RepositoryOpt is used to customise the functionality of statusdata.Repository.
type RepositoryOpt func(r *Repository) error

// WithLog can be used to pass a zap.Logger into the repository.
func WithLog(log *zap.Logger) RepositoryOpt {
	return func(r *Repository) error {
		r.log = log
		return nil
	}
}

// WithRefresh configures the repository to periodically reload the Nagios statusdata file.
func WithRefresh(interval time.Duration) RepositoryOpt {
	return func(r *Repository) error {
		r.refreshInterval = interval

		go func() {
			for range time.Tick(r.refreshInterval) {
				r.load()
			}
		}()

		return nil
	}
}
