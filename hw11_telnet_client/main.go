package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.String("timeout", "10s", "timeout connection default 10s")
	host := flag.String("host", "127.0.0.1", "name host or ip address")
	port := flag.String("port", "4242", "port number")
	flag.Parse()

	outErr := os.Stderr
	tim, err := time.ParseDuration(*timeout)
	if err != nil {
		fmt.Fprintln(outErr, err)
		return
	}
	tc := NewTelnetClient(net.JoinHostPort(*host, *port), tim, os.Stdin, os.Stdout)
	err = tc.Connect()
	if err != nil {
		fmt.Fprintln(outErr, err)
		return
	}
	defer tc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		defer stop()
		if err := tc.Receive(); err != nil {
			fmt.Fprintln(outErr, err)
		}
	}()

	go func() {
		defer stop()
		if err := tc.Send(); err != nil {
			fmt.Fprintln(outErr, err)
		}
	}()

	<-ctx.Done()
}
