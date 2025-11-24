package product

import (
	"context"
	"fmt"

	"github.com/roppenlabs/rapid-product-catalog/internal/config"
	"github.com/roppenlabs/rapid-product-catalog/internal/types"
	logger "github.com/roppenlabs/rapido-logger-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	BulkCreateProducts(ctx context.Context, products []Product) (CreateProductsResponse, error)
	SearchProducts(ctx context.Context, params SearchParams) (SearchProductsResponse, error)
	GetProductByID(ctx context.Context, productID primitive.ObjectID) (*Product, error)
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

func (s *serviceImpl) BulkCreateProducts(ctx context.Context, products []Product) (CreateProductsResponse, error) {
	result, err := s.repository.CreateProducts(ctx, products)
	if err != nil {
		return CreateProductsResponse{}, err
	}

	totalProcessed := result.Created + result.Updated
	response := CreateProductsResponse{
		Success:    true,
		Message:    fmt.Sprintf("Successfully processed %d products (%d created, %d updated)", totalProcessed, result.Created, result.Updated),
		Created:    result.Created,
		Updated:    result.Updated,
		ProductIDs: result.ProductIDs,
	}
	return response, nil
}

func (s *serviceImpl) SearchProducts(ctx context.Context, params SearchParams) (SearchProductsResponse, error) {

	logger.Info(logger.Format{Message: "Searching products", Data: map[string]string{"params": fmt.Sprintf("%+v", params)}})

	products, err := s.repository.SearchProducts(ctx, params.Categories, params.Brands, params.MinPrice, params.MaxPrice, params.SearchText, params.Limit)
	if err != nil {
		return SearchProductsResponse{}, err
	}

	if len(products) == 0 {
		return SearchProductsResponse{}, types.NewNotFoundError("No products found matching the search criteria")
	}

	response := SearchProductsResponse{
		Success:  true,
		Message:  fmt.Sprintf("Found %d products", len(products)),
		Count:    len(products),
		Products: products,
	}

	return response, nil
}

func (s *serviceImpl) GetProductByID(ctx context.Context, productID primitive.ObjectID) (*Product, error) {
	logger.Info(logger.Format{Message: "Fetching product by ID", Data: map[string]string{"productID": productID.Hex()}})

	product, err := s.repository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}
