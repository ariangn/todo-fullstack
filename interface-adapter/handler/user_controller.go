package handler

import (
    "encoding/json"
    "net/http"

    "github.com/ariangn/todo-go/application/user"
    "github.com/ariangn/todo-go/interface-adapter/dto/request"
    "github.com/ariangn/todo-go/interface-adapter/dto/response"
)

type UserController struct {
    registerUC user.RegisterUseCase
    loginUC    user.LoginUseCase
}

func NewUserController(rUC user.RegisterUseCase, lUC user.LoginUseCase) *UserController {
    return &UserController{registerUC: rUC, loginUC: lUC}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
