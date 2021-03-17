package main

import (
	"bufio"
	"context"
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
	ctx     context.Context
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
	//ctx := context.Background()
	var err error
	//t.ctx, _ = context.WithTimeout(context.Background(), t.timeout*time.Second)
	//dialer := &net.Dialer{}
	//t.conn, err = dialer.DialContext(t.ctx, "tcp", t.address)
	//dialer.Timeout = 1*time.Millisecond
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	//t.conn, err = net.Dial("tcp", t.address)
	if err != nil {
		return fmt.Errorf("cannot connect to %s. Error: %w", t.address, err)
	}
	_, err = io.WriteString(os.Stderr, "...Connected to "+t.address+"\n")
	if err != nil {
		return fmt.Errorf("cannot write a string. Error: %w", err)
	}
	return nil
}

func (t *telnetClient) Send() error {
	scanner := bufio.NewScanner(t.in)
	for {
		if !scanner.Scan() {
			_, err := io.WriteString(os.Stderr, "...EOF\n")
			if err != nil {
				return fmt.Errorf("cannot write a string. Error: %w", err)
			}
			break
		}
		str := scanner.Text()
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from channel %v. Error: %w", t.in, err)
		}

		_, err := t.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		if err != nil {
			_, err := io.WriteString(os.Stderr, "...Connection was closed by peer\n")
			if err != nil {
				return fmt.Errorf("cannot write a string. Error: %w", err)
			}
			return nil
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
		_, err := io.WriteString(t.out, text+"\n")
		if err != nil {
			return fmt.Errorf("cannot write a string. Error: %w", err)
		}
	}
	return nil
}

func (t *telnetClient) Close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
