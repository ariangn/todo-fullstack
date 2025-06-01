// interface-adapter/handler/tag_controller_test.go
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

type mockCreateTagUseCase struct {
	Created *entity.Tag
	Err     error
}

func (m *mockCreateTagUseCase) Execute(ctx context.Context, userID, name string) (*entity.Tag, error) {
	return m.Created, m.Err
}

type mockListTagUseCase struct {
	Tags []*entity.Tag
	Err  error
}

func (m *mockListTagUseCase) Execute(ctx context.Context, userID string) ([]*entity.Tag, error) {
	return m.Tags, m.Err
}

type mockDeleteTagUseCase struct {
	Err error
}

func (m *mockDeleteTagUseCase) Execute(ctx context.Context, id string) error {
	return m.Err
}

func setupTagRouter(tc *TagController, secret string) *chi.Mux {
	r := chi.NewRouter()
	mountProtectedRouter(r, secret)

	r.Post("/tags", tc.Create)
	r.Get("/tags", tc.List)
	r.Delete("/tags/{id}", tc.Delete)
	return r
}

func TestCreateTag_Success(t *testing.T) {
	userID := "user-tag-1"
	now := time.Now().UTC()
	dummyTag := &entity.Tag{
		ID:        uuid.NewString(),
		Name:      "Important",
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tc := &TagController{
		createUC: &mockCreateTagUseCase{Created: dummyTag, Err: nil},
	}

	secret := "tag-secret"
	router := setupTagRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	payload := request.CreateTagDTO{Name: "Important"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp response.TagResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.ID != dummyTag.ID || resp.Name != dummyTag.Name {
		t.Errorf("unexpected create response: %+v", resp)
	}
}

func TestListTag_Success(t *testing.T) {
	userID := "user-tag-2"
	now := time.Now().UTC()
	tag1 := &entity.Tag{
		ID:        uuid.NewString(),
		Name:      "Work",
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	tag2 := &entity.Tag{
		ID:        uuid.NewString(),
		Name:      "Personal",
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tc := &TagController{
		listUC: &mockListTagUseCase{Tags: []*entity.Tag{tag1, tag2}, Err: nil},
	}

	secret := "tag-secret-2"
	router := setupTagRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var respList []response.TagResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &respList); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(respList) != 2 {
		t.Errorf("expected 2 tags, got %d", len(respList))
	}
}

func TestDeleteTag_Success(t *testing.T) {
	userID := "user-tag-3"
	dummyTag := &entity.Tag{
		ID:     uuid.NewString(),
		Name:   "Temp",
		UserID: userID,
	}

	tc := &TagController{
		deleteUC: &mockDeleteTagUseCase{Err: nil},
	}

	secret := "tag-secret-3"
	router := setupTagRouter(tc, secret)
	token := generateAuthToken(t, secret, userID)

	req := httptest.NewRequest(http.MethodDelete, "/tags/"+dummyTag.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
}
