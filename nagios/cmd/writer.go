package cmd

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// Number of 'in-flight' commands to store in the command writer channel.
const DefaultWriteBufferSize = 1000

const DefaultExternalCommandsFile = "/usr/share/nagios/rw/nagios.cmd"

// Writer writes to the Nagios external commands file.
//
// The writes are buffered in a queue, to handle cases where the nagios
// external commands file is temporarily unavailable.
//
// Writer implements io.Writer by pushing a command onto the channel.
// A consumer coroutine must be started via the Run method.
//
// The code assumes the external command file is a fifo, such that writes
// to the file will return an error when the underlying pipe has been closed.
type Writer struct {
	filename    string
	logger      *zap.Logger
	nonBlocking bool
	queue       chan string
}

type writerConfig struct {
	bufferSize  int
	filename    string
	logger      *zap.Logger
	nonBlocking bool
}

// NewWriter constructs an instance of Writer.
func NewWriter(opts ...WriterOption) (*Writer, error) {
	cfg := writerConfig{
		bufferSize: DefaultWriteBufferSize,
		filename:   DefaultExternalCommandsFile,
		logger:     zap.NewNop(),
	}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	return &Writer{
		filename: cfg.filename,
		logger: cfg.logger.Named("ExternalCommands").With(
			zap.String("command file", cfg.filename),
		),
		queue: make(chan string, cfg.bufferSize),
	}, nil
}

// Run launches an infinite loop. It should be run as a goroutine.
func (w *Writer) Run() {
	var f *os.File

	log := w.logger
	for {
		// the nagios command file is only available when nagios is running.
		// so, it is expected to not exist, for example when nagios is
		// in the middle of restarting.
		if f == nil {
			newFile, err := os.OpenFile(w.filename, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Warn("unable to open command file",
					zap.String("retry interval", "1s"),
				)
				time.Sleep(time.Second)
				continue
			}
			f = newFile
		}

		// read the next command from the queue and attempt to write to the
		// external command file
		cmd, ok := <-w.queue
		if !ok {
			continue
		}

		// if the external commands file goes away, requeue the command for
		// later and redo the loop.
		_, err := f.Write([]byte(cmd + "\n"))
		if err != nil {
			err = f.Sync()
		}
		if err != nil {
			f = nil
			log.Warn("unable to write to command file",
				zap.String("retry interval", "1s"),
				zap.Error(err),
			)
			go func() {
				w.queue <- cmd
			}()
			time.Sleep(time.Second)
			continue
		}

		log.Debug("wrote to command file",
			zap.String("command", cmd),
		)
	}
}

// Write appends a command string to the queue, for later writing.
//
// If the queue is in non-blocking mode, an error will be returned when the
// queue is full.
//
// If the queue is in blocking mode, the write will block until the queue
// empties.
//
// It is assumed that the input is a full, valid, nagios command, not
// terminated by a new-line.
func (w *Writer) Write(cmd []byte) (n int, err error) {
	if !w.nonBlocking {
		w.queue <- string(cmd)
		return len(cmd), nil
	}

	select {
	case w.queue <- string(cmd):
		return len(cmd), nil
	default:
		return 0, fmt.Errorf("channel blocked")
	}
}

// WriterOption passes in parameters to NewWriter().
type WriterOption func(*writerConfig) error

func WithFilename(f string) WriterOption {
	return func(cfg *writerConfig) error {
		cfg.filename = f
		return nil
	}
}

// WithLogger passes in a zap Logger to NewWriter().
func WithLogger(l *zap.Logger) WriterOption {
	return func(cfg *writerConfig) error {
		cfg.logger = l
		return nil
	}
}

// WithNonBlocking sets the blocking behaviour of the Write() method.
//
// Writers are blocking by default.
func WithNonBlocking(b bool) WriterOption {
	return func(cfg *writerConfig) error {
		cfg.nonBlocking = b
		return nil
	}
}
