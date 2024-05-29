package executil

import (
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	ExecuteTime = 60
)

func ExecuteTimeoutEnvByUser(binary string, args []string, seconds int, env []string, executeUser string) (string, error) {
	var output []byte
	var err error
	cmd := exec.Command(binary, args...)
	if executeUser != "" {
		targetUser, err := user.Lookup(executeUser)
		if err != nil {
			return "", err
		}
		uid, _ := strconv.Atoi(targetUser.Uid)
		gid, _ := strconv.Atoi(targetUser.Gid)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}
	cmd.Env = env
	done := make(chan struct{})

	go func() {
		output, err = cmd.CombinedOutput()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(time.Duration(seconds) * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("timeout executing: %v %v, output %v, error %v", binary, args, string(output), err)
	}

	if err != nil {
		if !strings.Contains(err.Error(), "no child processes") {
			return "", fmt.Errorf("failed to execute: %v %v, output %v, error %v", binary, args, string(output), err)
		} else {
			fmt.Printf("execute: %v %v, output %v, error %v", binary, args, string(output), err)
		}
	}
	return string(output), nil
}

func Execute(binary string, args []string) (string, error) {
	return ExecuteTimeout(binary, args, ExecuteTime)
}

func ExecuteByUser(binary string, args []string, executeUser string) (string, error) {
	return ExecuteTimeoutByUser(binary, args, 600, executeUser)
}

func ExecuteTimeout(binary string, args []string, seconds int) (string, error) {
	for _, arg := range args {
		if strings.Contains(arg, "`") {
			return "", fmt.Errorf("timeout executing: %v,error: %s contain special symbols", binary, arg)
		}
	}
	return ExecuteTimeoutEnv(binary, args, seconds, nil)
}

func ExecuteTimeoutByUser(binary string, args []string, seconds int, targetUser string) (string, error) {
	for _, arg := range args {
		if strings.Contains(arg, "`") {
			return "", fmt.Errorf("timeout executing: %v,error: %s contain special symbols", binary, arg)
		}
	}
	return ExecuteTimeoutEnvByUser(binary, args, seconds, nil, targetUser)
}

func CheckIfCmdlineArgvIsValid(args []string) bool {
	for _, arg := range args {
		if strings.Contains(arg, "`") {
			return false
		}
	}
	return true
}

func ExecuteTimeoutEnv(binary string, args []string, seconds int, env []string) (string, error) {
	var output []byte
	var err error
	cmd := exec.Command(binary, args...)
	cmd.Env = env

	done := make(chan struct{}, 1)

	go func() {
		output, err = cmd.CombinedOutput()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(time.Duration(seconds) * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("timeout executing: %v %v, output %v, error %v", binary, args, string(output), err)
	}

	if err != nil {
		if !strings.Contains(err.Error(), "no child processes") {
			return string(output), fmt.Errorf("failed to execute: %v %v, output %v, error %v", binary, args, string(output), err)
		} else {
			fmt.Printf("execute: %v %v, output %v, error %v", binary, args, string(output), err)
		}
	}
	return string(output), nil
}

func ExecuteCmdSplitStdoutStderr(binary string, arg []string, seconds int) (string, string, error) {
	if seconds == 0 {
		seconds = 60
	}
	var output []byte
	var err error
	var stderr bytes.Buffer
	cmd := exec.Command(binary, arg...)
	cmd.Stderr = &stderr
	done := make(chan struct{})

	go func() {
		output, err = cmd.Output()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(time.Duration(seconds) * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", "", fmt.Errorf("timeout executing: %v, arg %v,error %v", binary, arg, err)
	}

	if err != nil {
		if !strings.Contains(err.Error(), "no child processes") {
			return string(output), stderr.String(), fmt.Errorf("failed to execute: %v, arg %v, error %v", binary, arg, err)
		}
	}

	return string(output), stderr.String(), nil
}
