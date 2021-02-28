package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Подготовка переменных окружения.
	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return 1
			}
		} else {
			err := os.Setenv(key, val.Value)
			if err != nil {
				return 1
			}
		}
	}

	comm := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	var in bytes.Buffer
	var out bytes.Buffer
	var stderr bytes.Buffer
	comm.Stdin = &in
	comm.Stdout = &out
	comm.Stderr = &stderr
	err := comm.Run()
	fmt.Println(out.String()) //nolint:forbidigo
	if err != nil {
		fmt.Println("Error:", err, stderr.String())     //nolint:forbidigo
		if exitError, ok := err.(*exec.ExitError); ok { //nolint:errorlint
			return exitError.ExitCode()
		}
	}
	return
}
