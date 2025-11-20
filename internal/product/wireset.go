package product

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHandler,
	NewService,
	NewRepository,
)
