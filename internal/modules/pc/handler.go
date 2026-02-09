package pc

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

func (h *Handler) HandlePCs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createPC(w, r)
	case http.MethodGet:
		h.listPCs(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandlePCByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/pcs/")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid pc id"})
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.updatePC(w, r, id)
	case http.MethodDelete:
		h.deletePC(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandlePCsByClub(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	clubID, err := parseClubIDFromPath(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid club id"})
		return
	}

	pcs, err := h.service.GetPCsByClub(clubID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pcs)
}

func (h *Handler) createPC(w http.ResponseWriter, r *http.Request) {
	var input PC
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	pc, err := h.service.CreatePC(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, pc)
}

func (h *Handler) listPCs(w http.ResponseWriter, r *http.Request) {
	pcs, err := h.service.GetAllPCs()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to fetch pcs"})
		return
	}

	writeJSON(w, http.StatusOK, pcs)
}

func (h *Handler) updatePC(w http.ResponseWriter, r *http.Request, id int) {
	var input PC
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	pc, err := h.service.UpdatePC(id, input)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pc)
}

func (h *Handler) deletePC(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.DeletePC(id); err != nil {
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
	case ErrPCNotFound:
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case ErrInvalidPCInput:
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

func parseClubIDFromPath(path string) (int, error) {
	const prefix = "/clubs/"
	const suffix = "/pcs"

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return 0, strconv.ErrSyntax
	}

	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.TrimSuffix(trimmed, suffix)
	if trimmed == "" || strings.Contains(trimmed, "/") {
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(trimmed)
}
