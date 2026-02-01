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
	fmt.Println("Server started with Nuradilet,Bekarys")
	http.ListenAndServe(":9721", nil)

}
