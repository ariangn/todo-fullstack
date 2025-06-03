package di

import (
	"github.com/ariangn/todo-fullstack/backend/application/category"
	"github.com/ariangn/todo-fullstack/backend/application/tag"
	"github.com/ariangn/todo-fullstack/backend/application/todo"
	"github.com/ariangn/todo-fullstack/backend/application/user"
	"github.com/ariangn/todo-fullstack/backend/infrastructure/auth"
	"github.com/ariangn/todo-fullstack/backend/infrastructure/database"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/handler"
)

type Container struct {
	AuthClient         auth.AuthClientInterface
	UserController     *handler.UserController
	TodoController     *handler.TodoController
	CategoryController *handler.CategoryController
	TagController      *handler.TagController
}

func InitializeContainer() (*Container, error) {
	// ─── (1) Auth Client ───────────────────────────────────────────────────────
	// NewAuthClient no longer takes any arguments
	authClient := auth.NewAuthClient()

	// ─── (2) Supabase / DB Client ─────────────────────────────────────────────
	// NewSupabaseClient now returns ( *SupabaseClient, error )
	supabaseClient, err := database.NewSupabaseClient()
	if err != nil {
		return nil, err
	}

	// ─── (3) Repositories ─────────────────────────────────────────────────────
	userRepo := database.NewUserRepository(supabaseClient)
	todoRepo := database.NewTodoRepository(supabaseClient)
	categoryRepo := database.NewCategoryRepository(supabaseClient)
	tagRepo := database.NewTagRepository(supabaseClient)

	// ─── (4) User Use‐Cases ────────────────────────────────────────────────────
	registerUC := user.NewRegisterUseCase(userRepo)
	// LoginUseCase expects (UserRepository, AuthClientInterface)
	loginUC := user.NewLoginUseCase(userRepo, authClient)
	// FindByIDUseCase expects (UserRepository)
	findByIDUC := user.NewFindByIDUseCase(userRepo)

	// ─── (5) Todo Use‐Cases ────────────────────────────────────────────────────
	// Note: NewCreateUseCase requires (TodoRepository, CategoryRepository, TagRepository)
	createTodoUC := todo.NewCreateUseCase(todoRepo, categoryRepo, tagRepo)
	listTodoUC := todo.NewListUseCase(todoRepo)
	findTodoByIDUC := todo.NewFindByIDUseCase(todoRepo)
	updateTodoUC := todo.NewUpdateUseCase(todoRepo)
	toggleStatusUC := todo.NewToggleStatusUseCase(todoRepo)
	deleteTodoUC := todo.NewDeleteUseCase(todoRepo)
	duplicateTodoUC := todo.NewDuplicateUseCase(todoRepo)

	// ─── (6) Category Use‐Cases ────────────────────────────────────────────────
	createCategoryUC := category.NewCreateUseCase(categoryRepo)
	listCategoryUC := category.NewListUseCase(categoryRepo)
	updateCategoryUC := category.NewUpdateUseCase(categoryRepo)
	deleteCategoryUC := category.NewDeleteUseCase(categoryRepo)

	// ─── (7) Tag Use‐Cases ─────────────────────────────────────────────────────
	createTagUC := tag.NewCreateUseCase(tagRepo)
	listTagUC := tag.NewListUseCase(tagRepo)
	updateTagUC := tag.NewUpdateUseCase(tagRepo)
	deleteTagUC := tag.NewDeleteUseCase(tagRepo)

	// ─── (8) Controllers ───────────────────────────────────────────────────────
	userController := handler.NewUserController(registerUC, loginUC, findByIDUC)

	// NewTodoController signature is:
	//   NewTodoController(
	//     CreateUseCase,
	//     ListUseCase,
	//     FindByIDUseCase,
	//     UpdateUseCase,
	//     ToggleStatusUseCase,
	//     DeleteUseCase,
	//     DuplicateUseCase,
	//   )
	todoController := handler.NewTodoController(
		createTodoUC,
		listTodoUC,
		findTodoByIDUC,
		updateTodoUC,
		toggleStatusUC,
		deleteTodoUC,
		duplicateTodoUC,
	)

	categoryController := handler.NewCategoryController(
		createCategoryUC,
		listCategoryUC,
		updateCategoryUC,
		deleteCategoryUC,
	)

	tagController := handler.NewTagController(
		createTagUC,
		listTagUC,
		updateTagUC,
		deleteTagUC,
	)

	return &Container{
		AuthClient:         authClient,
		UserController:     userController,
		TodoController:     todoController,
		CategoryController: categoryController,
		TagController:      tagController,
	}, nil
}
