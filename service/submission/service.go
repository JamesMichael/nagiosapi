package submission

import (
	"fmt"
	"io"
	"time"

	"github.com/jamesmichael/nagiosapi/nagios/cmd"
)

// Service is used to write passive results to the nagios external commands
// file.
type Service struct {
	externalCommandsFile io.Writer
}

// NewService constructs an instance of Service.
func NewService(opts ...ServiceOption) (*Service, error) {
	s := Service{}
	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	if s.externalCommandsFile == nil {
		return nil, fmt.Errorf("must set external commands file")
	}

	return &s, nil
}

// SubmitPassiveResult constructs a nagios command for submitting a passive
// service check result and queues the command for writing to the nagios
// external commands file.
//
// If the ServiceResult does not have a time set, the current unix timestamp
// is used instead.
func (s *Service) SubmitPassiveResult(r *ServiceResult) error {
	t := r.Time
	if t == 0 {
		t = time.Now().Unix()
	}

	command := fmt.Sprintf("[%d] PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s",
		t,
		cmd.Sanitize(r.Hostname),
		cmd.Sanitize(r.ServiceName),
		r.Status,
		cmd.Sanitize(r.Body),
	)
	if _, err := s.externalCommandsFile.Write([]byte(command)); err != nil {
		return err
	}

	return nil
}

// ServiceOption passes parameters to NewService().
type ServiceOption func(s *Service) error

// WithExternalCommandsWriter sets the external commands writer.
//
// It should be an instance of nagios/cmd, but for testing, anything which
// implements io.Writer would work.
func WithExternalCommandsWriter(w io.Writer) ServiceOption {
	return func(s *Service) error {
		s.externalCommandsFile = w
		return nil
	}
}
