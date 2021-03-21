package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	timeout := flag.Duration("timeout", 10, "timeout connection default 10s")
	host := flag.String("host", "127.0.0.1", "name host or ip address")
	port := flag.String("port", "4242", "port number")
	flag.Parse()

	outErr := os.Stderr
	tc := NewTelnetClient(net.JoinHostPort(*host, *port), *timeout, os.Stdin, os.Stdout)
	err := tc.Connect()
	if err != nil {
		fmt.Fprintln(outErr, err)
		os.Exit(1)
	}
	defer tc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		defer stop()
		if err := tc.Receive(); err != nil {
			fmt.Fprintln(outErr, err)
			os.Exit(1)
		}
	}()

	go func() {
		defer stop()
		if err := tc.Send(); err != nil {
			fmt.Fprintln(outErr, err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
}
