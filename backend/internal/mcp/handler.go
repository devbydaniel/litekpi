package mcp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// Handler handles HTTP requests for MCP API keys and MCP protocol.
type Handler struct {
	service *Service
}

// NewHandler creates a new MCP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateKey handles creating a new MCP API key.
//
//	@Summary		Create MCP API key
//	@Description	Create a new MCP API key for the organization (shown only once). Requires admin role.
//	@Tags			mcp
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateKeyRequest	true	"Key data"
//	@Success		201		{object}	CreateKeyResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/mcp/keys [post]
func (h *Handler) CreateKey(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req CreateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	response, err := h.service.CreateKey(r.Context(), user.OrganizationID, user.ID, req)
	if err != nil {
		if errors.Is(err, ErrKeyNameEmpty) {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "API key name is required"})
			return
		}
		log.Printf("create MCP API key error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create MCP API key"})
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// ListKeys handles listing MCP API keys for the organization.
//
//	@Summary		List MCP API keys
//	@Description	Get all MCP API keys for the authenticated user's organization. Requires admin role.
//	@Tags			mcp
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListKeysResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/mcp/keys [get]
func (h *Handler) ListKeys(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	keys, err := h.service.ListKeys(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list MCP API keys error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to list MCP API keys"})
		return
	}

	respondJSON(w, http.StatusOK, ListKeysResponse{Keys: keys})
}

// DeleteKey handles deleting an MCP API key.
//
//	@Summary		Delete MCP API key
//	@Description	Delete an MCP API key by ID. Requires admin role.
//	@Tags			mcp
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Key ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/mcp/keys/{id} [delete]
func (h *Handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	keyID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid key ID"})
		return
	}

	err = h.service.DeleteKey(r.Context(), user.OrganizationID, keyID)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			respondJSON(w, http.StatusNotFound, ErrorResponse{Error: "MCP API key not found"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondJSON(w, http.StatusForbidden, ErrorResponse{Error: "unauthorized"})
			return
		}
		log.Printf("delete MCP API key error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to delete MCP API key"})
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "MCP API key deleted"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
