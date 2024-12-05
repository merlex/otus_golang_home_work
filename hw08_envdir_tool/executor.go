package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}
	if cmd[0] == "" {
		return 1
	}
	cmdpath := cmd[0]
	command := exec.Command(cmdpath, cmd[1:]...)
	osEnv := make(map[string]string, len(os.Environ()))
	for _, s := range os.Environ() {
		envStringSlice := strings.Split(s, "=")
		osEnv[envStringSlice[0]] = envStringSlice[1]
	}
	for name, value := range env {
		if !value.NeedRemove {
			osEnv[name] = value.Value
		} else {
			delete(osEnv, name)
		}
	}
	for n, v := range osEnv {
		command.Env = append(command.Env, n+"="+v)
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		var cmdErr *exec.ExitError
		if errors.As(err, &cmdErr) {
			return cmdErr.ExitCode()
		}
	}
	return
}
