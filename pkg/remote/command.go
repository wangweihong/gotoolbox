package remote

import (
	"fmt"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"golang.org/x/crypto/ssh"
)

// SSHCommand execute command on remote ssh
type SSHCommand struct {
	Session *ssh.Session
}

// Exec executes a command on a specific SSH session and returns the output
func (s *SSHCommand) Output(command string) (string, error) {
	output, err := s.Session.Output(command)
	if err != nil {
		return "", errors.UpdateStack(err)
	}

	return string(output), nil
}

func (s *SSHCommand) Close() error {
	if err := s.Session.Close(); err != nil {
		return fmt.Errorf("failed to close SSH session: %v", err)
	}
	return nil
}
