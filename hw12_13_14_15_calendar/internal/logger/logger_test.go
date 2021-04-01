package logger

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	file := "/tmp/logfile.log"
	t.Run("no such file or directory", func(t *testing.T) {
		_, err := New("INFO", "")
		require.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("incorrect level", func(t *testing.T) {
		_, err := New("test", file)
		require.NotNil(t, err)
		defer os.Remove(file) // очистка
	})

	t.Run("log OK", func(t *testing.T) {
		log, err := New("INFO", file)
		require.NoError(t, err)
		defer os.Remove(file) // очистка
		require.True(t, log != nil)
	})

	t.Run("checking text writing to a file", func(t *testing.T) {
		log, err := New("DEBUG", file)
		require.NoError(t, err)
		defer os.Remove(file) // очистка

		log.Info("test")
		f1, err := ioutil.ReadFile(file)
		if err != nil {
			return
		}
		require.Contains(t, string(f1), "test")
		require.Contains(t, string(f1), "INFO")

		log.Warn("test2")
		f1, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}
		require.Contains(t, string(f1), "test2")
		require.Contains(t, string(f1), "WARN")

		log.Error("test3")
		f1, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}
		require.Contains(t, string(f1), "test3")
		require.Contains(t, string(f1), "ERRO")

		log.Debug("test4", "logger")
		f1, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}
		require.Contains(t, string(f1), "test4")
		require.Contains(t, string(f1), "DEBU")
	})
}
