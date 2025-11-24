package product

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) BulkCreateProducts(ctx context.Context, products []Product) (CreateProductsResponse, error) {
	ret := s.Mock.Called(ctx, products)
	return ret.Get(0).(CreateProductsResponse), ret.Error(1)
}

func (s *MockService) SearchProducts(ctx context.Context, params SearchParams) (SearchProductsResponse, error) {
	ret := s.Mock.Called(ctx, params)
	return ret.Get(0).(SearchProductsResponse), ret.Error(1)
}

func (s *MockService) GetProductByID(ctx context.Context, productID primitive.ObjectID) (*Product, error) {
	ret := s.Mock.Called(ctx, productID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*Product), ret.Error(1)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateProducts(ctx context.Context, products []Product) (*CreateProductsResult, error) {
	ret := m.Mock.Called(ctx, products)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*CreateProductsResult), ret.Error(1)
}

func (m *MockRepository) SearchProducts(ctx context.Context, categories []string, brands []string, minPrice, maxPrice *float64, searchText string, limit int) ([]Product, error) {
	ret := m.Mock.Called(ctx, categories, brands, minPrice, maxPrice, searchText, limit)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]Product), ret.Error(1)
}

func (m *MockRepository) GetProductByID(ctx context.Context, productID primitive.ObjectID) (*Product, error) {
	ret := m.Mock.Called(ctx, productID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*Product), ret.Error(1)
}
