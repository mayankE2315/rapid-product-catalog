package product

import (
	"context"
	"testing"

	"github.com/roppenlabs/rapid-product-catalog/internal/config"
	logger "github.com/roppenlabs/rapido-logger-go"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ProductUploadServiceTestSuite struct {
	suite.Suite
	config config.Config
}

func (mps *ProductUploadServiceTestSuite) SetupTest() {
	cfg, err := config.NewConfig()
	if err != nil {
		mps.T().Fatalf("Failed to initialize config: %v", err)
	}
	mps.config = cfg
	logger.Init(mps.config.Get().Log.Level)
}

func TestProductUploadServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductUploadServiceTestSuite))
}

func (mps *ProductUploadServiceTestSuite) TestShouldReturnSuccessWithValidProducts() {
	products := []Product{
		{
			Name:        "Titan Edge 1",
			Category:    "watch",
			Brand:       "titan",
			Price:       12999,
			Description: "Titan Edge Slim Series",
			Images:      []string{"https://cdn.example.com/titan1.png"},
			Inventory:   20,
			Popularity:  4.5,
		},
		{
			Name:        "Apple iPhone 16",
			Category:    "mobile",
			Brand:       "apple",
			Price:       50000,
			Description: "Latest iPhone model",
			Images:      []string{"https://cdn.example.com/iphone16.png"},
			Inventory:   5,
			Popularity:  4.8,
		},
	}

	mockRepo := new(MockRepository)
	// Return result with IDs to simulate database response
	mockResult := &CreateProductsResult{
		ProductIDs: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		Created:    2,
		Updated:    0,
	}
	mockRepo.On("CreateProducts", mock.Anything, products).Return(mockResult, nil)

	testService := NewService(mps.config, mockRepo)
	resp, err := testService.BulkCreateProducts(context.Background(), products)

	assert.Nil(mps.T(), err)
	assert.Equal(mps.T(), 2, resp.Created)
	assert.Equal(mps.T(), 0, resp.Updated)
	assert.Equal(mps.T(), 2, len(resp.ProductIDs))
}
