package booking

import (
	"encoding/json"
	"net/http"
)

func HandleBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var b Booking
		json.NewDecoder(r.Body).Decode(&b)
		CreateBooking(b)
		w.WriteHeader(http.StatusCreated)
		return
	}

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(GetAll())
	}
}
