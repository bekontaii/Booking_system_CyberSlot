package main

import (
	"fmt"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/app"
)

func main() {
	fmt.Println("Booking_system_CyberSlot server started with Nuradilet,Bekarys")
	app.RunServer()
}
func CreateBooking(b Booking) {
	Save(b)
	AutoExpire()
}
