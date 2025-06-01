// interface-adapter/handler/category_controller_test.go
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

	"github.com/ariangn/todo-go/domain/entity"
	"github.com/ariangn/todo-go/interface-adapter/dto/request"
	"github.com/ariangn/todo-go/interface-adapter/dto/response"
)

type mockCreateCategoryUseCase struct {
	Created *entity.Category
	Err     error
}

func (m *mockCreateCategoryUseCase) Execute(
	ctx context.Context,
	userID, name, color string,
	description *string,
) (*entity.Category, error) {
	return m.Created, m.Err
}

type mockListCategoryUseCase struct {
	Categories []*entity.Category
	Err        error
}

func (m *mockListCategoryUseCase) Execute(ctx context.Context, userID string) ([]*entity.Category, error) {
	return m.Categories, m.Err
}

type mockDeleteCategoryUseCase struct {
	Err error
}

func (m *mockDeleteCategoryUseCase) Execute(ctx context.Context, id string) error {
	return m.Err
}

func setupCategoryRouter(cc *CategoryController, secret string) *chi.Mux {
	r := chi.NewRouter()
	mountProtectedRouter(r, secret)

	r.Post("/categories", cc.Create)
	r.Get("/categories", cc.List)
	r.Delete("/categories/{id}", cc.Delete)
	return r
}

func TestCreateCategory_Success(t *testing.T) {
	userID := "user-cat-1"
	now := time.Now().UTC()
	dummyCat := &entity.Category{
		ID:          uuid.NewString(),
		Name:        "Work",
		Color:       "#ff0000",
		Description: ptrString("Work items"),
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	cc := &CategoryController{
		createUC: &mockCreateCategoryUseCase{Created: dummyCat, Err: nil},
	}

	secret := "cat-secret"
	router := setupCategoryRouter(cc, secret)
	token := generateAuthToken(t, secret, userID)

	payload := request.CreateCategoryDTO{
		Name:        "Work",
		Color:       "#ff0000",
		Description: ptrString("Work items"),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp response.CategoryResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ID != dummyCat.ID || resp.Name != dummyCat.Name {
		t.Errorf("unexpected create response: %+v", resp)
	}
}

func TestListCategory_Success(t *testing.T) {
	userID := "user-cat-2"
	now := time.Now().UTC()
	cat1 := &entity.Category{
		ID:        uuid.NewString(),
		Name:      "Home",
		Color:     "#00ff00",
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	cat2 := &entity.Category{
		ID:        uuid.NewString(),
		Name:      "Errands",
		Color:     "#0000ff",
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	cc := &CategoryController{
		listUC: &mockListCategoryUseCase{Categories: []*entity.Category{cat1, cat2}, Err: nil},
	}

	secret := "cat-secret-2"
	router := setupCategoryRouter(cc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var respList []response.CategoryResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &respList); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(respList) != 2 {
		t.Errorf("expected 2 categories, got %d", len(respList))
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	userID := "user-cat-3"
	dummyCat := &entity.Category{
		ID:     uuid.NewString(),
		Name:   "Temp",
		Color:  "#abcd12",
		UserID: userID,
	}

	cc := &CategoryController{
		deleteUC: &mockDeleteCategoryUseCase{Err: nil},
	}

	secret := "cat-secret-3"
	router := setupCategoryRouter(cc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+dummyCat.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
}
