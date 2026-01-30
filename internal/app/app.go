package app

import (
	"fmt"
	"net/http"
)

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Booking PC")
}
func RunServer() {
	http.HandleFunc("/booking", BookingHandler)
	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
