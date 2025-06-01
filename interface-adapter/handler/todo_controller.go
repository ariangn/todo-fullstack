package handler

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/ariangn/todo-go/application/todo"
    "github.com/ariangn/todo-go/domain/entity"
    "github.com/ariangn/todo-go/domain/value-object"
    "github.com/ariangn/todo-go/interface-adapter/dto/request"
    "github.com/ariangn/todo-go/interface-adapter/dto/response"
)

type TodoController struct {
    createUC     todo.CreateUseCase
    listUC       todo.ListUseCase
    updateUC     todo.UpdateUseCase
    toggleStatus todo.ToggleStatusUseCase
    deleteUC     todo.DeleteUseCase
    duplicateUC  todo.DuplicateUseCase
}

func NewTodoController(
    cUC todo.CreateUseCase,
    lUC todo.ListUseCase,
    uUC todo.UpdateUseCase,
    tUC todo.ToggleStatusUseCase,
    dUC todo.DeleteUseCase,
    dupUC todo.DuplicateUseCase,
) *TodoController {
    return &TodoController{
        createUC:     cUC,
        listUC:       lUC,
        updateUC:     uUC,
        toggleStatus: tUC,
        deleteUC:     dUC,
        duplicateUC:  dupUC,
    }
}

func (tc *TodoController) Create(w http.ResponseWriter, r *http.Request) {
    userID, ok := GetUserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var dto request.CreateTodoDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        http.Error(w, "invalid request payload", http.StatusBadRequest)
        return
    }

    // convert dueDate to valueobject.DueDateVO if provided
    var dueDateVO *value-object.DueDateVO
    if dto.DueDate != nil {
        dvo, err := value-object.NewDueDateVO(*dto.DueDate)
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
        dueDateVO,
        dto.CategoryID,
        dto.TagIDs,
    )
    if err != nil {
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
    userID, _ := GetUserIDFromContext(r.Context())
    todos, err := tc.listUC.Execute(r.Context(), userID)
    if err != nil {
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
    id := chi.URLParam(r, "id")
    todoEntity, err := tc.listUC.Execute(r.Context(), /*write a FindByIDUseCase or just call repo directly*/)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // TODO: implement full GetByID here, call a FindByIDUseCase
    w.WriteHeader(http.StatusNotImplemented)
}

func (tc *TodoController) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var dto request.UpdateTodoDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        http.Error(w, "invalid request payload", http.StatusBadRequest)
        return
    }
    // fetch existing, apply updates, then call tc.updateUC.Execute(...)
    w.WriteHeader(http.StatusNotImplemented)
}

func (tc *TodoController) ToggleStatus(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var body struct {
        Status string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request payload", http.StatusBadRequest)
        return
    }
    newStatus := entity.Status(body.Status)
    updated, err := tc.toggleStatus.Execute(r.Context(), id, newStatus)
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

func (tc *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := tc.deleteUC.Execute(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (tc *TodoController) Duplicate(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
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
