package club

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleClubs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createClub(w, r)
	case http.MethodGet:
		h.listClubs(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleClubByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/clubs/")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid club id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getClub(w, r, id)
	case http.MethodPut:
		h.updateClub(w, r, id)
	case http.MethodDelete:
		h.deleteClub(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createClub(w http.ResponseWriter, r *http.Request) {
	var input Club
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	club, err := h.service.CreateClub(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, club)
}

func (h *Handler) listClubs(w http.ResponseWriter, r *http.Request) {
	clubs, err := h.service.GetClubs()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to fetch clubs"})
		return
	}

	writeJSON(w, http.StatusOK, clubs)
}

func (h *Handler) getClub(w http.ResponseWriter, r *http.Request, id int) {
	club, err := h.service.GetClubByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, club)
}

func (h *Handler) updateClub(w http.ResponseWriter, r *http.Request, id int) {
	var input Club
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	club, err := h.service.UpdateClub(id, input)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, club)
}

func (h *Handler) deleteClub(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.DeleteClub(id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch err {
	case ErrClubNotFound:
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case ErrInvalidClubInput:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func parseID(path string, prefix string) (int, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, strconv.ErrSyntax
	}

	idPart := strings.TrimPrefix(path, prefix)
	if idPart == "" || strings.Contains(idPart, "/") {
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(idPart)
}
