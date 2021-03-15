package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Настройка переменных окружения.
func settingEnv(key, value string, needRemove bool) error {
	if needRemove {
		return os.Unsetenv(key)
	}
	return os.Setenv(key, value)
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Подготовка переменных окружения.
	for key, val := range env {
		err := settingEnv(key, val.Value, val.NeedRemove)
		if err != nil {
			return 1
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
