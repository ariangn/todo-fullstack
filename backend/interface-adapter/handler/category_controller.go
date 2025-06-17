package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ariangn/todo-fullstack/backend/application/category"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/request"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/response"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
)

type CategoryController struct {
	createUC category.CreateUseCase
	listUC   category.ListUseCase
	updateUC category.UpdateUseCase
	deleteUC category.DeleteUseCase
}

func NewCategoryController(
	cUC category.CreateUseCase,
	lUC category.ListUseCase,
	uUC category.UpdateUseCase,
	dUC category.DeleteUseCase,
) *CategoryController {
	return &CategoryController{cUC, lUC, uUC, dUC}
}

func (cc *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	var dto request.CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	categoryEntity, err := cc.createUC.Execute(
		r.Context(),
		userID,
		dto.Name,
		dto.Color,
		dto.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respDTO := response.CategoryResponseDTO{
		ID:          categoryEntity.ID,
		Name:        categoryEntity.Name,
		Color:       categoryEntity.Color,
		Description: categoryEntity.Description,
		UserID:      categoryEntity.UserID,
		CreatedAt:   categoryEntity.CreatedAt,
		UpdatedAt:   categoryEntity.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respDTO)
}

func (cc *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	categories, err := cc.listUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var respList []response.CategoryResponseDTO
	for _, c := range categories {
		respList = append(respList, response.CategoryResponseDTO{
			ID:          c.ID,
			Name:        c.Name,
			Color:       c.Color,
			Description: c.Description,
			UserID:      c.UserID,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respList)
}

func (cc *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("DEBUG: could not read body: %v", err)
	}
	// restore r.Body for the rest of the handler
	r.Body = io.NopCloser(bytes.NewBuffer(raw))

	// 1) auth
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	// 2) path param
	id := chi.URLParam(r, "id")

	// 3) ensure body wasn’t empty
	if len(raw) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "request body was empty"})
		return
	}

	// 4) decode into DTO
	var dto request.UpdateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	// 5) execute update use-case
	updatedEntity, err := cc.updateUC.Execute(
		r.Context(),
		userID,
		id,
		getString(dto.Name),
		getString(dto.Color),
		getString(dto.Description),
	)
	if err != nil {
		log.Printf("DEBUG: updateUC.Execute returned error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 6) success response
	respDTO := response.CategoryResponseDTO{
		ID:          updatedEntity.ID,
		Name:        updatedEntity.Name,
		Color:       updatedEntity.Color,
		Description: updatedEntity.Description,
		UserID:      updatedEntity.UserID,
		CreatedAt:   updatedEntity.CreatedAt,
		UpdatedAt:   updatedEntity.UpdatedAt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respDTO)
}

func (cc *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := cc.deleteUC.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// getString safely dereferences a *string, returning an empty string if nil.
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
