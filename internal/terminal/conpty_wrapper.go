//go:build windows

package terminal

import "github.com/UserExistsError/conpty"

// ConPtyWrapper wraps the ConPTY to implement ConPtyIO interface
type ConPtyWrapper struct {
	*conpty.ConPty
}

func (c *ConPtyWrapper) Resize(width, height int) error {
	return c.ConPty.Resize(width, height)
}
