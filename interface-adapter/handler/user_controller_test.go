// interface-adapter/handler/user_controller_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ariangn/todo-go/domain/entity"
	"github.com/ariangn/todo-go/interface-adapter/dto/request"
	"github.com/ariangn/todo-go/interface-adapter/dto/response"
)

type mockRegisterUseCase struct {
	Created *entity.User
	Err     error
}

func (m *mockRegisterUseCase) Execute(
	ctx context.Context,
	email, password string,
	name *string,
	timezone string,
	avatarURL *string,
) (*entity.User, error) {
	return m.Created, m.Err
}

type mockLoginUseCase struct {
	Token string
	Err   error
}

func (m *mockLoginUseCase) Execute(
	ctx context.Context,
	email, password string,
) (string, error) {
	return m.Token, m.Err
}

func TestRegister_Success(t *testing.T) {
	now := time.Now().UTC()
	dummyUser := &entity.User{
		ID:        uuid.NewString(),
		Email:     "bob@example.com",
		Password:  "hashedpassword",
		Name:      ptrString("Bob"),
		AvatarURL: ptrString("https://example.com/avatar.png"),
		Timezone:  "America/Chicago",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctrl := &UserController{
		registerUC: &mockRegisterUseCase{Created: dummyUser, Err: nil},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", ctrl.Register)

	payload := request.CreateUserDTO{
		Email:     "bob@example.com",
		Password:  "plaintext",
		Name:      ptrString("Bob"),
		Timezone:  "America/Chicago",
		AvatarURL: ptrString("https://example.com/avatar.png"),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp response.UserResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.Email != dummyUser.Email || resp.Name == nil || *resp.Name != *dummyUser.Name {
		t.Errorf("unexpected register response: %+v", resp)
	}
}

func TestLogin_Success(t *testing.T) {
	mockToken := "mock-jwt-token"
	ctrl := &UserController{
		loginUC: &mockLoginUseCase{Token: mockToken, Err: nil},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/login", ctrl.Login)

	payload := request.LoginUserDTO{
		Email:    "bob@example.com",
		Password: "plaintext",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp["token"] != mockToken {
		t.Errorf("unexpected login response: %+v", resp)
	}
}
