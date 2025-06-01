//go:build wireinject
// +build wireinject

package di

import (
    "github.com/google/wire"
)

func InitializeEverything() (*handler.Server, error) {
    wire.Build(ProviderSet)
    return &handler.Server{}, nil
}
