package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		fromTEst := "testdata/input.txt"
		toTEst, err := ioutil.TempFile("", "out.*.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTEst.Name()) // очистка

		offsetTest := int64(7000)
		limitTest := int64(0)
		err = Copy(fromTEst, toTEst.Name(), offsetTest, limitTest)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})

	t.Run("unsupported file", func(t *testing.T) {
		fromTest := "testdata/input.txt"
		toTest := ""
		err := Copy(fromTest, toTest, 0, 0)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))

		fromTest = ""
		toTEst, err := ioutil.TempFile("", "out.*.txt")

		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTEst.Name()) // очистка

		offsetTest := int64(0)
		limitTest := int64(0)
		err = Copy(fromTest, toTEst.Name(), offsetTest, limitTest)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))

		fromTest = "/dev/urandom"
		err = Copy(fromTest, toTEst.Name(), offsetTest, limitTest)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))

	})

	t.Run("full copied successfully", func(t *testing.T) {
		fromTest := "testdata/input.txt"
		toTest, err := ioutil.TempFile("", "out.*.txt")

		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTest.Name()) // очистка

		offsetTest := int64(0)
		limitTest := int64(0)
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		info, _ := toTest.Stat()

		require.Nil(t, err)
		require.True(t, info.Size() == 6617)
	})

	t.Run("offset = 100, 1000 byte copied successfully", func(t *testing.T) {
		fromTest := "testdata/input.txt"
		toTest, err := ioutil.TempFile("", "out.*.txt")

		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTest.Name()) // очистка

		offsetTest := int64(100)
		limitTest := int64(1000)
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.Nil(t, err)
		f1, err := ioutil.ReadFile(toTest.Name())
		if err != nil {
			return
		}

		f2, err := ioutil.ReadFile("testdata/out_offset100_limit1000.txt")
		if err != nil {
			return
		}

		require.True(t, bytes.Equal(f1, f2))
	})

	t.Run("end of file reached", func(t *testing.T) {
		fromTest := "testdata/input.txt"
		toTest, err := ioutil.TempFile("", "out.*.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTest.Name()) // очистка

		offsetTest := int64(6000)
		limitTest := int64(1000)
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.Nil(t, err)
		info, _ := toTest.Stat()

		require.Nil(t, err)
		require.True(t, info.Size() == 617)
	})
}
