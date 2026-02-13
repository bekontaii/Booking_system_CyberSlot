package booking

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/middleware"
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

func (h *Handler) HandleBookingByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "booking not found"})
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleClubBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	if authUser.ClubID == nil {
		writeJSON(w, http.StatusOK, h.service.GetAllBookings())
		return
	}
	writeJSON(w, http.StatusOK, h.service.GetBookingsByClub(*authUser.ClubID))
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser.UserID <= 0 {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	req.UserID = authUser.UserID

	booking, err := h.service.CreateBooking(req)
	if err != nil {
		switch err {
		case ErrInvalidTimeRange:
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		case ErrClubInactive:
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

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.DeleteBooking(id); err != nil {
		switch err {
		case ErrBookingNotFound:
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to delete booking"})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
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

func parseID(path string) (int, bool) {
	trimmed := strings.TrimPrefix(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return 0, false
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
