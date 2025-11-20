package product

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) BulkCreateProducts(ctx context.Context, products []Product) (BulkCreateProductsResponse, error) {
	ret := s.Mock.Called(ctx, products)
	return ret.Get(0).(BulkCreateProductsResponse), ret.Error(1)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateProducts(ctx context.Context, products []Product) ([]Product, error) {
	ret := m.Mock.Called(ctx, products)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]Product), ret.Error(1)
}
