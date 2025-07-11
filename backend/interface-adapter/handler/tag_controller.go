package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ariangn/todo-fullstack/backend/application/tag"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/request"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/response"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
)

type TagController struct {
	createUC tag.CreateUseCase
	listUC   tag.ListUseCase
	updateUC tag.UpdateUseCase
	deleteUC tag.DeleteUseCase
}

func NewTagController(
	cUC tag.CreateUseCase,
	lUC tag.ListUseCase,
	uUC tag.UpdateUseCase,
	dUC tag.DeleteUseCase,
) *TagController {
	return &TagController{cUC, lUC, uUC, dUC}
}

func (tc *TagController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}

	// read and log body safely
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// rewind for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var dto request.CreateTagDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	tagEntity, err := tc.createUC.Execute(r.Context(), userID, dto.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respDTO := response.TagResponseDTO{
		ID:        tagEntity.ID,
		Name:      tagEntity.Name,
		UserID:    tagEntity.UserID,
		CreatedAt: tagEntity.CreatedAt,
		UpdatedAt: tagEntity.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respDTO)
}

func (tc *TagController) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	tags, err := tc.listUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var respList []response.TagResponseDTO
	for _, t := range tags {
		respList = append(respList, response.TagResponseDTO{
			ID:        t.ID,
			Name:      t.Name,
			UserID:    t.UserID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respList)
}

func (tc *TagController) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (tc *TagController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := tc.deleteUC.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
