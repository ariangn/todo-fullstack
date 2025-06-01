//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/ariangn/todo-go/application/category"
	"github.com/ariangn/todo-go/application/tag"
	"github.com/ariangn/todo-go/application/todo"
	"github.com/ariangn/todo-go/application/user"
	"github.com/ariangn/todo-go/infrastructure/auth"
	"github.com/ariangn/todo-go/infrastructure/database"
	"github.com/ariangn/todo-go/interface-adapter/handler"
)

// Container holds all the components
type Container struct {
	// AuthClient allows ValidateToken / GenerateToken
	AuthClient auth.AuthClientInterface

	// Supabase client used by all repositories
	SupabaseClient *database.SupabaseClient

	// HTTP controllers (handlers)
	UserController     *handler.UserController
	TodoController     *handler.TodoController
	CategoryController *handler.CategoryController
	TagController      *handler.TagController
}

// ProvideSupabaseClient constructs a SupabaseClient
func ProvideSupabaseClient() (*database.SupabaseClient, error) {
	return database.NewSupabaseClient()
}

var RepositorySet = wire.NewSet(
	database.NewUserRepository,
	database.NewTodoRepository,
	database.NewCategoryRepository,
	database.NewTagRepository,
)

var AuthSet = wire.NewSet(
	auth.NewAuthClient, // reads JWT_SECRET from env
)

var UserUseCaseSet = wire.NewSet(
	user.NewRegisterUseCase,
	user.NewLoginUseCase,
)

var TodoUseCaseSet = wire.NewSet(
	todo.NewCreateUseCase,
	todo.NewListUseCase,
	todo.NewUpdateUseCase,
	todo.NewToggleStatusUseCase,
	todo.NewDeleteUseCase,
	todo.NewDuplicateUseCase,
)

var CategoryUseCaseSet = wire.NewSet(
	category.NewCreateUseCase,
	category.NewListUseCase,
	category.NewUpdateUseCase,
	category.NewDeleteUseCase,
)

var TagUseCaseSet = wire.NewSet(
	tag.NewCreateUseCase,
	tag.NewListUseCase,
	tag.NewUpdateUseCase,
	tag.NewDeleteUseCase,
)

var HandlerSet = wire.NewSet(
	handler.NewUserController,
	handler.NewTodoController,
	handler.NewCategoryController,
	handler.NewTagController,
)

var ProviderSet = wire.NewSet(
	ProvideSupabaseClient,
	RepositorySet,
	AuthSet,
	UserUseCaseSet,
	TodoUseCaseSet,
	CategoryUseCaseSet,
	TagUseCaseSet,
	HandlerSet,

	// Tell Wire how to build a Container from available providers
	wire.Struct(
		new(Container),
		"AuthClient",
		"SupabaseClient",
		"UserController",
		"TodoController",
		"CategoryController",
		"TagController",
	),
)

// InitializeContainer is the Wire‐generated entry point (see wire_gen.go)
func InitializeContainer() (*Container, error) {
	wire.Build(ProviderSet)
	return &Container{}, nil
}
