package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Defining duration
	const d = time.Second

	// current time system
	currentTime := time.Now()
	fmt.Println("current time:", currentTime.Round(d))

	// точное время
	exactTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		log.Fatalf("ntp error: %s", err)
	}
	fmt.Println("exact time:", exactTime.Round(d))
}
