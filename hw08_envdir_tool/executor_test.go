package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("map comparison", func(t *testing.T) {
		dir := "./testdata/env"
		mapEnv, err := ReadDir(dir)
		require.Nil(t, err)

		cmd := []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}
		exitCode := RunCmd(cmd, mapEnv)
		require.Equal(t, exitCode, 0, "maps don't match")

		cmd = []string{"/bin/bash", "./testdata/echo.s", "arg1=1", "arg2=2"}
		exitCode = RunCmd(cmd, mapEnv)
		require.Equal(t, exitCode, 127, "maps don't match")
	})
}
