package handler

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"

    "github.com/ariangn/todo-go/application/category"
    "github.com/ariangn/todo-go/interface-adapter/dto/request"
    "github.com/ariangn/todo-go/interface-adapter/dto/response"
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
    userID, _ := GetUserIDFromContext(r.Context())
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
    userID, _ := GetUserIDFromContext(r.Context())
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
    // similar to Create/Update in Todo
    w.WriteHeader(http.StatusNotImplemented)
}

func (cc *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := cc.deleteUC.Execute(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
