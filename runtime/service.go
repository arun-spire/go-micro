package runtime

import (
	"io"
	"strings"
	"sync"

	packager "github.com/arun-spire/go-micro/runtime/package"
	"github.com/arun-spire/go-micro/runtime/process"
	proc "github.com/arun-spire/go-micro/runtime/process/os"
	"github.com/arun-spire/go-micro/util/log"
)

type service struct {
	sync.RWMutex

	running bool
	closed  chan bool
	err     error

	// output for logs
	output io.Writer

	// service to manage
	*Service
	// process creator
	Process *proc.Process
	// Exec
	Exec *process.Executable
	// process pid
	PID *process.PID
}

func newService(s *Service, c CreateOptions) *service {
	var exec string
	var args []string

	if len(s.Exec) > 0 {
		parts := strings.Split(s.Exec, " ")
		exec = parts[0]
		args = []string{}

		if len(parts) > 1 {
			args = parts[1:]
		}
	} else {
		// set command
		exec = c.Command[0]
		// set args
		if len(c.Command) > 1 {
			args = c.Command[1:]
		}
	}

	return &service{
		Service: s,
		Process: new(proc.Process),
		Exec: &process.Executable{
			Binary: &packager.Binary{
				Name: s.Name,
				Path: exec,
			},
			Env:  c.Env,
			Args: args,
		},
		closed: make(chan bool),
		output: c.Output,
	}
}

func (s *service) streamOutput() {
	go io.Copy(s.output, s.PID.Output)
	go io.Copy(s.output, s.PID.Error)
}

// Running returns true is the service is running
func (s *service) Running() bool {
	s.RLock()
	defer s.RUnlock()
	return s.running
}

// Start stars the service
func (s *service) Start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	// reset
	s.err = nil
	s.closed = make(chan bool)

	// TODO: pull source & build binary
	log.Debugf("Runtime service %s forking new process", s.Service.Name)
	p, err := s.Process.Fork(s.Exec)
	if err != nil {
		return err
	}

	// set the pid
	s.PID = p
	// set to running
	s.running = true

	if s.output != nil {
		s.streamOutput()
	}

	// wait and watch
	go s.Wait()

	return nil
}

// Stop stops the service
func (s *service) Stop() error {
	s.Lock()
	defer s.Unlock()

	select {
	case <-s.closed:
		return nil
	default:
		close(s.closed)
		s.running = false
		if s.PID == nil {
			return nil
		}
		return s.Process.Kill(s.PID)
	}
}

// Error returns the last error service has returned
func (s *service) Error() error {
	s.RLock()
	defer s.RUnlock()
	return s.err
}

// Wait waits for the service to finish running
func (s *service) Wait() {
	// wait for process to exit
	err := s.Process.Wait(s.PID)

	s.Lock()
	defer s.Unlock()

	// save the error
	if err != nil {
		s.err = err
	}

	// no longer running
	s.running = false
}
