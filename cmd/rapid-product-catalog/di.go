//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/roppenlabs/rapid-product-catalog/internal/config"
	"github.com/roppenlabs/rapid-product-catalog/internal/health"
	"github.com/roppenlabs/rapid-product-catalog/internal/product"
	"github.com/roppenlabs/rapid-product-catalog/internal/server"
	"github.com/roppenlabs/rapid-product-catalog/internal/utils"
)

type ServerDependencies struct {
	config   config.Config
	server   *server.Server
	handlers server.Handlers
}

func InitDependencies() (ServerDependencies, error) {
	wire.Build(
		wire.Struct(new(ServerDependencies), "*"),
		wire.Struct(new(server.Handlers), "*"),
		server.WireSet,
		product.WireSet,
		health.WireSet,
		utils.WireSet,
		config.GetConfig,
	)

	return ServerDependencies{}, nil
}
