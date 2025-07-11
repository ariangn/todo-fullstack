package handler

import (
	"encoding/json"
	"net/http"

	"bytes"
	"io"

	"github.com/ariangn/todo-fullstack/backend/application/user"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/request"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/dto/response"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
)

type UserController struct {
	registerUC user.RegisterUseCase
	loginUC    user.LoginUseCase
	findByIDUC user.FindByIDUseCase // ← new use‐case for fetching by ID
}

func NewUserController(
	rUC user.RegisterUseCase,
	lUC user.LoginUseCase,
	fbUC user.FindByIDUseCase,
) *UserController {
	return &UserController{
		registerUC: rUC,
		loginUC:    lUC,
		findByIDUC: fbUC,
	}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Re-use body
	var dto request.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	userEntity, err := uc.registerUC.Execute(
		r.Context(),
		dto.Email,
		dto.Password,
		dto.Name,
		dto.Timezone,
		dto.AvatarURL,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respDTO := response.UserResponseDTO{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		AvatarURL: userEntity.AvatarURL,
		Timezone:  userEntity.Timezone,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respDTO)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var dto request.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := uc.loginUC.Execute(r.Context(), dto.Email, dto.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   60 * 60 * 24, // 1 day
	})

	// Optionally return user info or a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Expire the cookie by setting MaxAge to -1
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // HTTPS only in prod
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logout successful"})
}

// Me returns the currently authenticated user's info (requires AuthMiddleware).
func (uc *UserController) Me(w http.ResponseWriter, r *http.Request) {
	// Extract userID from context (populated by AuthMiddleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Use the FindByID use-case to fetch user details
	userEntity, err := uc.findByIDUC.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	respDTO := response.UserResponseDTO{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		AvatarURL: userEntity.AvatarURL,
		Timezone:  userEntity.Timezone,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respDTO)
}
