package app

import (
	"fmt"
	"net/http"
)

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Booking PC", r.URL.Path[1:])
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Home", r.URL.Path[1:])
}
func RunServer() {
	http.HandleFunc("/booking", BookingHandler)
	http.HandleFunc("/", HomeHandler)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)

}
