package main

import (
	"log"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}
