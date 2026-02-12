package booking

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/middleware"
	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
)

type Handler struct {
	service  *Service
	userRepo user.Repository
}

func NewHandler(service *Service, userRepo user.Repository) *Handler {
	return &Handler{service: service, userRepo: userRepo}
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
	case http.MethodPut, http.MethodPatch:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
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

	username, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || username == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	u, err := h.userRepo.GetByUsername(username)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	req.UserID = u.ID

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

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	booking, err := h.service.UpdateBooking(id, req)
	if err != nil {
		switch err {
		case ErrInvalidTimeRange, ErrInvalidStatus:
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		case ErrBookingConflict:
			writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
		case ErrBookingNotFound:
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to update booking"})
		}
		return
	}

	writeJSON(w, http.StatusOK, booking)
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
