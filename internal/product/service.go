package product

import (
	"context"
	"fmt"

	"github.com/roppenlabs/rapid-product-catalog/internal/config"
)

type Service interface {
	BulkCreateProducts(ctx context.Context, products []Product) (BulkCreateProductsResponse, error)
}

type serviceImpl struct {
	cfg        config.Config
	repository Repository
}

func NewService(cfg config.Config, repo Repository) Service {
	service := &serviceImpl{
		cfg:        cfg,
		repository: repo,
	}
	return service
}

func (s *serviceImpl) BulkCreateProducts(ctx context.Context, products []Product) (BulkCreateProductsResponse, error) {
	updatedProducts, err := s.repository.CreateProducts(ctx, products)
	if err != nil {
		return BulkCreateProductsResponse{}, err
	}

	response := BulkCreateProductsResponse{
		Success:  true,
		Message:  fmt.Sprintf("Successfully processed %d products (created or updated)", len(updatedProducts)),
		Created:  len(updatedProducts),
		Products: updatedProducts,
	}
	return response, nil
}
