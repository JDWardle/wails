//go:build darwin || linux
// +build darwin linux

package process

import "syscall"

// Stop sends a SIGTERM signal to the process.
func (p *Process) Stop() error {
	if !p.Running {
		return nil
	}

	if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	return nil
}
