package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("test EOF", func(t *testing.T) {
		reader, writer, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		stdout := os.Stdout
		stderr := os.Stderr
		defer func() {
			os.Stdout = stdout
			os.Stderr = stderr
		}()

		os.Stdout = writer
		os.Stderr = writer

		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			var outData bytes.Buffer
			r, w := io.Pipe()

			timeout, err := time.ParseDuration("5s")
			require.NoError(t, err)

			tc := NewTelnetClient(l.Addr().String(), timeout, r, &outData)
			require.NoError(t, tc.Connect())
			defer func() { require.NoError(t, tc.Close()) }()

			err = w.Close()
			require.NoError(t, err)
			err = tc.Send()
			require.NoError(t, err)

			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			require.NoError(t, err)
			require.Contains(t, string(buf[:n]), "...EOF")
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()
		}()

		wg.Wait()
	})

	t.Run("test connection failed", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer
		address := "google.com:80"
		timeout := time.Duration(10)

		client := NewTelnetClient(address, timeout, in, out)
		err := client.Connect()
		require.Equal(t, "cannot connect to google.com:80. Error: dial tcp: i/o timeout", err.Error())
	})

	t.Run("missing address", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer
		address := ""
		timeout, err := time.ParseDuration("1s")
		require.Nil(t, err)

		client := NewTelnetClient(address, timeout, in, out)
		err = client.Connect()
		// require.Error(t, client.Connect())
		require.Equal(t, "cannot connect to . Error: dial tcp: missing address", err.Error())
	})
}
