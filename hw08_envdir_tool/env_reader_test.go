package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expected := make(Environment)
	expected["BAR"] = EnvValue{"bar", false}
	expected["EMPTY"] = EnvValue{"", false}
	expected["FOO"] = EnvValue{"   foo\nwith new line", false}
	expected["HELLO"] = EnvValue{"\"hello\"", false}
	expected["UNSET"] = EnvValue{"", true}

	f := func(t *testing.T) {
		dir := "./testdata/env"
		mapEnv, err := ReadDir(dir)
		require.Nil(t, err)
		require.Equal(t, expected, mapEnv, "maps don't match")

	}
	t.Run("map comparison", f)

	t.Run("path error", func(t *testing.T) {
		dir := "env"
		_, err := ReadDir(dir)
		require.NotNil(t, err)
	})

	t.Run("filename contains character '='", func(t *testing.T) {
		envTest, err := ioutil.TempFile("./testdata/env", "out=*")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(envTest.Name())

		dir := "./testdata/env"
		mapEnv, err := ReadDir(dir)

		require.Nil(t, err)
		require.Equal(t, expected, mapEnv, "maps don't match")
	})
}
