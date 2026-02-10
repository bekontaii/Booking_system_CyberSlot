package booking

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleBookings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodGet:
		h.handleList(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	booking, err := h.service.CreateBooking(req)
	if err != nil {
		switch err {
		case ErrInvalidTimeRange:
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		case ErrBookingConflict:
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to create booking"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	bookings := h.service.GetAllBookings()
	writeJSON(w, http.StatusOK, bookings)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
