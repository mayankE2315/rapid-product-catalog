package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/roppenlabs/rapid-product-catalog/internal/config"
	"github.com/roppenlabs/rapid-product-catalog/internal/health"
	"github.com/roppenlabs/rapid-product-catalog/internal/product"
	logger "github.com/roppenlabs/rapido-logger-go"
)

type Handlers struct {
	HealthHandler  *health.Handler
	ProductHandler *product.Handler
}

func (s *Server) InitRoutes(h Handlers, c config.Config) {
	router := s.routerGroups.rootRouter

	// Health routes
	router.GET("/sanity", h.HealthHandler.CheckSanity)
	router.GET("/health", h.HealthHandler.CheckHealth)

	// Product routes
	router.POST("/products/bulk", h.ProductHandler.CreateProductsHandler)
	router.POST("/products/search", h.ProductHandler.SearchProductsHandler)
	router.GET("/products/:productId", h.ProductHandler.GetProductByIDHandler)

	// Register pprof handlers
	if c.Get().ProfilingEnabled {
		logger.Info(logger.Format{
			Message: "ALERT! Profiling enabled. Please be aware of the performance impact it could have",
		})
		pprof.Register(router)
	}
}
