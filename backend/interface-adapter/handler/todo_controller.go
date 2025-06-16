// interface-adapter/handler/todo_controller.go
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ariangn/todo-fullstack/backend/application/todo"
	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/valueobject"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/request"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/response"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
)

type TodoController struct {
	createUC     todo.CreateUseCase
	listUC       todo.ListUseCase
	findByIDUC   todo.FindByIDUseCase
	updateUC     todo.UpdateUseCase
	toggleStatus todo.ToggleStatusUseCase
	deleteUC     todo.DeleteUseCase
	duplicateUC  todo.DuplicateUseCase
}

func NewTodoController(
	cUC todo.CreateUseCase,
	lUC todo.ListUseCase,
	fUC todo.FindByIDUseCase,
	uUC todo.UpdateUseCase,
	tUC todo.ToggleStatusUseCase,
	dUC todo.DeleteUseCase,
	dupUC todo.DuplicateUseCase,
) *TodoController {
	return &TodoController{
		createUC:     cUC,
		listUC:       lUC,
		findByIDUC:   fUC,
		updateUC:     uUC,
		toggleStatus: tUC,
		deleteUC:     dUC,
		duplicateUC:  dupUC,
	}
}

func (tc *TodoController) Create(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	fmt.Println("RAW TODO BODY:", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // rewind for decoding

	var dto request.CreateTodoDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Println("❌ Failed to decode todo DTO:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	fmt.Printf("CONTROLLER DECODED DTO: %+v\n", dto)

	// Convert dueDate to valueobject.DueDateVO if provided
	var dueDateVO *valueobject.DueDateVO
	if dto.DueDate != nil {
		dvo, err := valueobject.NewDueDateVO(*dto.DueDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dueDateVO = &dvo
	}

	todoEntity, err := tc.createUC.Execute(
		r.Context(),
		userID,
		dto.Title,
		dto.Body,
		entity.Status(dto.Status),
		dueDateVO,
		dto.CategoryID,
		dto.TagIDs,
	)
	if err != nil {
		log.Printf("❌ Usecase CreateTodo error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respDTO := response.TodoResponseDTO{
		ID:          todoEntity.ID,
		Title:       todoEntity.Title,
		Body:        todoEntity.Body,
		Status:      string(todoEntity.Status),
		DueDate:     todoEntity.DueDate,
		CompletedAt: todoEntity.CompletedAt,
		UserID:      todoEntity.UserID,
		CategoryID:  todoEntity.CategoryID,
		TagIDs:      todoEntity.TagIDs,
		CreatedAt:   todoEntity.CreatedAt,
		UpdatedAt:   todoEntity.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respDTO)
}

func (tc *TodoController) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	todos, err := tc.listUC.Execute(r.Context(), userID)
	if err != nil {
		log.Printf("error in listUC.Execute: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var respList []response.TodoResponseDTO
	for _, t := range todos {
		respList = append(respList, response.TodoResponseDTO{
			ID:          t.ID,
			Title:       t.Title,
			Body:        t.Body,
			Status:      string(t.Status),
			DueDate:     t.DueDate,
			CompletedAt: t.CompletedAt,
			UserID:      t.UserID,
			CategoryID:  t.CategoryID,
			TagIDs:      t.TagIDs,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respList)
}

func (tc *TodoController) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	todoEntity, err := tc.findByIDUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todoEntity == nil || todoEntity.UserID != userID {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	respDTO := response.TodoResponseDTO{
		ID:          todoEntity.ID,
		Title:       todoEntity.Title,
		Body:        todoEntity.Body,
		Status:      string(todoEntity.Status),
		DueDate:     todoEntity.DueDate,
		CompletedAt: todoEntity.CompletedAt,
		UserID:      todoEntity.UserID,
		CategoryID:  todoEntity.CategoryID,
		TagIDs:      todoEntity.TagIDs,
		CreatedAt:   todoEntity.CreatedAt,
		UpdatedAt:   todoEntity.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respDTO)
}

func (tc *TodoController) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	existing, err := tc.findByIDUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil || existing.UserID != userID {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	var dto request.UpdateTodoDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// Apply updates to the existing entity
	if dto.Title != nil {
		existing.Title = *dto.Title
	}
	if dto.Body != nil {
		existing.Body = dto.Body
	}
	if dto.Status != nil {
		existing.Status = entity.Status(*dto.Status)
		if existing.Status == entity.StatusCompleted {
			now := time.Now().UTC()
			existing.CompletedAt = &now
		} else {
			existing.CompletedAt = nil
		}
	}
	if dto.DueDate != nil {
		// dto.DueDate is *time.Time; use its value directly
		dvo, err := valueobject.NewDueDateVO(*dto.DueDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t := dvo.Time()
		existing.DueDate = &t
	}
	if dto.CategoryID != nil {
		existing.CategoryID = dto.CategoryID
	}
	if dto.TagIDs != nil {
		existing.TagIDs = *dto.TagIDs
	}
	existing.UpdatedAt = time.Now().UTC()

	updated, err := tc.updateUC.Execute(r.Context(), existing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respDTO := response.TodoResponseDTO{
		ID:          updated.ID,
		Title:       updated.Title,
		Body:        updated.Body,
		Status:      string(updated.Status),
		DueDate:     updated.DueDate,
		CompletedAt: updated.CompletedAt,
		UserID:      updated.UserID,
		CategoryID:  updated.CategoryID,
		TagIDs:      updated.TagIDs,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respDTO)
}

func (tc *TodoController) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	existing, err := tc.findByIDUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil || existing.UserID != userID {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Failed to decode body: %v", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	newStatus := entity.Status(body.Status)
	updated, err := tc.toggleStatus.Execute(r.Context(), id, newStatus)
	if err != nil {
		log.Printf("Failed to toggle status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respDTO := response.TodoResponseDTO{
		ID:          updated.ID,
		Title:       updated.Title,
		Body:        updated.Body,
		Status:      string(updated.Status),
		DueDate:     updated.DueDate,
		CompletedAt: updated.CompletedAt,
		UserID:      updated.UserID,
		CategoryID:  updated.CategoryID,
		TagIDs:      updated.TagIDs,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respDTO)
}

func (tc *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	existing, err := tc.findByIDUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil || existing.UserID != userID {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	if err := tc.deleteUC.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (tc *TodoController) Duplicate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	existing, err := tc.findByIDUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil || existing.UserID != userID {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	dup, err := tc.duplicateUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respDTO := response.TodoResponseDTO{
		ID:          dup.ID,
		Title:       dup.Title,
		Body:        dup.Body,
		Status:      string(dup.Status),
		DueDate:     dup.DueDate,
		CompletedAt: dup.CompletedAt,
		UserID:      dup.UserID,
		CategoryID:  dup.CategoryID,
		TagIDs:      dup.TagIDs,
		CreatedAt:   dup.CreatedAt,
		UpdatedAt:   dup.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respDTO)
}
