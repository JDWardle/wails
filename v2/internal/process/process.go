package process

import (
	"context"
	"os"
	"os/exec"
)

// Process defines a process that can be executed
type Process struct {
	cmd         *exec.Cmd
	exitChannel chan struct{}
	Running     bool
}

// NewProcess creates a new process struct
func NewProcess(cmd string, args ...string) *Process {
	result := &Process{
		cmd:         exec.Command(cmd, args...),
		exitChannel: make(chan struct{}),
	}
	result.cmd.Stdout = os.Stdout
	result.cmd.Stderr = os.Stderr
	return result
}

// Start the process
func (p *Process) Start(exitCodeChannel chan int) error {
	err := p.cmd.Start()
	if err != nil {
		return err
	}

	p.Running = true

	go func(cmd *exec.Cmd, running *bool, exitCodeChannel chan int) {
		defer close(p.exitChannel)

		err := cmd.Wait()
		if err == nil {
			exitCodeChannel <- 0
		}
		*running = false
	}(p.cmd, &p.Running, exitCodeChannel)

	return nil
}

// Kill the process
func (p *Process) Kill() error {
	if !p.Running {
		return nil
	}
	err := p.cmd.Process.Kill()
	if err != nil {
		return err
	}
	err = p.cmd.Process.Release()
	if err != nil {
		return err
	}

	return err
}

// KillWait kills the process once the provided context is canceled. If the
// process terminates while waiting this returns early without error.
func (p *Process) KillWait(ctx context.Context) error {
	select {
	case <-ctx.Done():
	case <-p.exitChannel:
		return nil
	}

	return p.Kill()
}

// PID returns the process PID
func (p *Process) PID() int {
	return p.cmd.Process.Pid
}

func (p *Process) SetDir(dir string) {
	p.cmd.Dir = dir
}
