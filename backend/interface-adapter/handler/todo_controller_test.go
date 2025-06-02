// interface-adapter/handler/todo_controller_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ariangn/todo-fullstack/backend/domain/entity"
	"github.com/ariangn/todo-fullstack/backend/domain/valueobject"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/request"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/response"
)

type mockCreateUseCase struct {
	Created *entity.Todo
	Err     error
}

func (m *mockCreateUseCase) Execute(
	ctx context.Context,
	userID, title string,
	body *string,
	dueDate *valueobject.DueDateVO,
	categoryID *string,
	tagIDs []string,
) (*entity.Todo, error) {
	return m.Created, m.Err
}

type mockListUseCase struct {
	Todos []*entity.Todo
	Err   error
}

func (m *mockListUseCase) Execute(ctx context.Context, userID string) ([]*entity.Todo, error) {
	return m.Todos, m.Err
}

type mockFindByIDUseCase struct {
	Todo *entity.Todo
	Err  error
}

func (m *mockFindByIDUseCase) Execute(ctx context.Context, id string) (*entity.Todo, error) {
	return m.Todo, m.Err
}

type mockUpdateUseCase struct {
	Updated *entity.Todo
	Err     error
}

func (m *mockUpdateUseCase) Execute(ctx context.Context, t *entity.Todo) (*entity.Todo, error) {
	return m.Updated, m.Err
}

type mockToggleStatusUseCase struct {
	Updated *entity.Todo
	Err     error
}

func (m *mockToggleStatusUseCase) Execute(ctx context.Context, id string, newStatus entity.Status) (*entity.Todo, error) {
	return m.Updated, m.Err
}

type mockDeleteUseCase struct {
	Err error
}

func (m *mockDeleteUseCase) Execute(ctx context.Context, id string) error {
	return m.Err
}

type mockDuplicateUseCase struct {
	Dup *entity.Todo
	Err error
}

func (m *mockDuplicateUseCase) Execute(ctx context.Context, id string) (*entity.Todo, error) {
	return m.Dup, m.Err
}

func setupTodoRouter(tc *TodoController, secret string) *chi.Mux {
	r := chi.NewRouter()
	mountProtectedRouter(r, secret)

	r.Post("/todos", tc.Create)
	r.Get("/todos", tc.List)
	r.Get("/todos/{id}", tc.GetByID)
	r.Put("/todos/{id}", tc.Update)
	r.Put("/todos/{id}/status", tc.ToggleStatus)
	r.Delete("/todos/{id}", tc.Delete)
	r.Post("/todos/{id}/duplicate", tc.Duplicate)
	return r
}

func TestCreateTodo_Success(t *testing.T) {
	userID := "user-123"
	now := time.Now().UTC()
	dummyTodo := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Test Create",
		Body:      ptrString("Details"),
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tc := &TodoController{
		createUC: &mockCreateUseCase{Created: dummyTodo, Err: nil},
	}

	secret := "secret1"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	createDTO := request.CreateTodoDTO{
		Title: "Test Create",
		Body:  ptrString("Details"),
	}
	payload, _ := json.Marshal(createDTO)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp response.TodoResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ID != dummyTodo.ID || resp.Title != dummyTodo.Title {
		t.Errorf("unexpected create response: %+v", resp)
	}
}

func TestListTodos_Success(t *testing.T) {
	userID := "user-456"
	now := time.Now().UTC()
	todo1 := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "First Todo",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	todo2 := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Second Todo",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tc := &TodoController{
		listUC: &mockListUseCase{Todos: []*entity.Todo{todo1, todo2}, Err: nil},
	}

	secret := "secret2"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var respList []response.TodoResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &respList); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(respList) != 2 {
		t.Errorf("expected 2 todos, got %d", len(respList))
	}
}

func TestGetByID_Success(t *testing.T) {
	userID := "user-789"
	now := time.Now().UTC()
	dummyTodo := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Fetch Todo",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tc := &TodoController{
		findByIDUC: &mockFindByIDUseCase{Todo: dummyTodo, Err: nil},
	}

	secret := "secret3"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodGet, "/todos/"+dummyTodo.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp response.TodoResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ID != dummyTodo.ID || resp.Title != dummyTodo.Title {
		t.Errorf("unexpected get response: %+v", resp)
	}
}

func TestUpdateTodo_Success(t *testing.T) {
	userID := "user-abc"
	now := time.Now().UTC()
	existing := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Old Title",
		Body:      ptrString("Old Body"),
		Status:    entity.StatusTodo,
		UserID:    userID,
		DueDate:   ptrTime(now.Add(24 * time.Hour)),
		CreatedAt: now,
		UpdatedAt: now,
	}
	updated := &entity.Todo{
		ID:        existing.ID,
		Title:     "New Title",
		Body:      ptrString("New Body"),
		Status:    entity.StatusInProgress,
		UserID:    userID,
		DueDate:   ptrTime(now.Add(48 * time.Hour)),
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now.Add(1 * time.Hour),
	}

	tc := &TodoController{
		findByIDUC: &mockFindByIDUseCase{Todo: existing, Err: nil},
		updateUC:   &mockUpdateUseCase{Updated: updated, Err: nil},
	}

	secret := "secret4"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	updateDTO := request.UpdateTodoDTO{
		Title:   ptrString("New Title"),
		Body:    ptrString("New Body"),
		Status:  ptrString(string(entity.StatusInProgress)),
		DueDate: ptrTime(now.Add(48 * time.Hour)),
	}
	payload, _ := json.Marshal(updateDTO)

	req := httptest.NewRequest(http.MethodPut, "/todos/"+existing.ID, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp response.TodoResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Title != updated.Title || (resp.Body == nil && updated.Body != nil) {
		t.Errorf("unexpected update response: %+v", resp)
	}
}

func TestToggleStatus_Success(t *testing.T) {
	userID := "user-def"
	now := time.Now().UTC()
	existing := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Toggle Me",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	toggled := &entity.Todo{
		ID:          existing.ID,
		Title:       existing.Title,
		Status:      entity.StatusCompleted,
		CompletedAt: ptrTime(now.Add(2 * time.Hour)),
		UserID:      userID,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   now.Add(2 * time.Hour),
	}

	tc := &TodoController{
		findByIDUC:   &mockFindByIDUseCase{Todo: existing, Err: nil},
		toggleStatus: &mockToggleStatusUseCase{Updated: toggled, Err: nil},
	}

	secret := "secret5"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	body := map[string]string{"status": string(entity.StatusCompleted)}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/todos/"+existing.ID+"/status", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp response.TodoResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Status != string(entity.StatusCompleted) || resp.CompletedAt == nil {
		t.Errorf("unexpected toggle response: %+v", resp)
	}
}

func TestDeleteTodo_Success(t *testing.T) {
	userID := "user-ghi"
	now := time.Now().UTC()
	existing := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Delete Me",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tc := &TodoController{
		findByIDUC: &mockFindByIDUseCase{Todo: existing, Err: nil},
		deleteUC:   &mockDeleteUseCase{Err: nil},
	}

	secret := "secret6"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodDelete, "/todos/"+existing.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
}

func TestDuplicateTodo_Success(t *testing.T) {
	userID := "user-jkl"
	now := time.Now().UTC()
	original := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Original",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	duplicate := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "Original (Copy)",
		Status:    entity.StatusTodo,
		UserID:    userID,
		CreatedAt: now.Add(1 * time.Hour),
		UpdatedAt: now.Add(1 * time.Hour),
	}

	tc := &TodoController{
		findByIDUC:  &mockFindByIDUseCase{Todo: original, Err: nil},
		duplicateUC: &mockDuplicateUseCase{Dup: duplicate, Err: nil},
	}

	secret := "secret7"
	router := setupTodoRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodPost, "/todos/"+original.ID+"/duplicate", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp response.TodoResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ID != duplicate.ID || resp.Title != duplicate.Title {
		t.Errorf("unexpected duplicate response: %+v", resp)
	}
}
