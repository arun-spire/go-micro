// Package process executes a binary
package process

import (
	"io"

	"github.com/arun-spire/go-micro/runtime/package"
)

// Process manages a running process
type Process interface {
	// Executes a process to completion
	Exec(*Executable) error
	// Creates a new process
	Fork(*Executable) (*PID, error)
	// Kills the process
	Kill(*PID) error
	// Waits for a process to exit
	Wait(*PID) error
}

type Executable struct {
	// The executable binary
	Binary *packager.Binary
	// The env variables
	Env []string
	// Args to pass
	Args []string
}

// PID is the running process
type PID struct {
	// ID of the process
	ID string
	// Stdin
	Input io.Writer
	// Stdout
	Output io.Reader
	// Stderr
	Error io.Reader
}
