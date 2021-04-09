package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect to %s. Error: %w", t.address, err)
	}
	fmt.Fprintln(os.Stderr, "...Connected to", t.address)
	return nil
}

func (t *telnetClient) Send() error {
	scanner := bufio.NewScanner(t.in)
	for {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "...EOF")
			break
		}
		str := scanner.Text()
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from channel %v. Error: %w", t.in, err)
		}

		_, err := t.conn.Write([]byte(fmt.Sprintln(str)))
		if err != nil {
			return fmt.Errorf("...Connection was closed by peer")
		}
	}
	return nil
}

func (t *telnetClient) Receive() error {
	scanner := bufio.NewScanner(t.conn)
	for {
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from channel %v. Error: %w", t.conn, err)
		}
		fmt.Fprintln(t.out, text)
	}
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return fmt.Errorf("channel nil")
	}
	return t.conn.Close()
}
