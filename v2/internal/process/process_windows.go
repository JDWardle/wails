//go:build windows
// +build windows

package process

// Stop is a no-op on Windows. Interrupts are not supported on Windows for
// detached processes https://github.com/golang/go/issues/46345.
func (p *Process) Stop() error { return nil }
