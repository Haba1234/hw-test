package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	timeout := flag.String("timeout", "10s", "timeout connection")
	host := flag.String("host", "127.0.0.1", "name host or ip address")
	port := flag.String("port", "4242", "port number")
	flag.Parse()

	tim, err := time.ParseDuration(*timeout)
	if err != nil {
		return
	}
	tc := NewTelnetClient(net.JoinHostPort(*host, *port), tim, os.Stdin, os.Stdout)
	err = tc.Connect()
	if err != nil {
		_, err = io.WriteString(os.Stderr, fmt.Sprintf("Error: %v\n", err))
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tc.Receive()
		if err != nil {
			_, err = io.WriteString(os.Stderr, fmt.Sprintf("Error: %v\n", err))
			if err != nil {
				log.Fatalln(err)
			}
		}
		tc.Close()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tc.Send()
		if err != nil {
			_, err = io.WriteString(os.Stderr, fmt.Sprintf("Error: %v\n", err))
			if err != nil {
				log.Fatalln(err)
			}
		}
		tc.Close()
	}()
	wg.Wait()

}
