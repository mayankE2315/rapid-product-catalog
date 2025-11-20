package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/roppenlabs/rapid-product-catalog/internal/types"
	logger "github.com/roppenlabs/rapido-logger-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/roppenlabs/rapid-product-catalog/internal/testutils"
	"github.com/stretchr/testify/suite"
)

type ProductUploadHandlerTestSuite struct {
	suite.Suite
	service  *MockService
	validate *validator.Validate
	server   *testutils.TestServer
	handler  *Handler
}

func (mph *ProductUploadHandlerTestSuite) SetupTest() {
	mph.service = new(MockService)
	mph.server = testutils.NewServer()
	mph.handler = NewHandler(mph.service)
	mph.handler.InitRoutes(mph.server.Router())
	logger.Init("debug")
}

func (mph *ProductUploadHandlerTestSuite) TestShouldReturnSuccessWithValidInput() {
	products := []Product{
		{
			Name:        "Titan Edge 1",
			Category:    "watch",
			Brand:       "titan",
			Price:       12999,
			Description: "Titan Edge Slim Series",
			Images:      []string{"https://cdn.example.com/titan1.png"},
			Inventory:   20,
		},
	}
	requestBody := BulkCreateProductsRequest{Products: products}
	expectedResponse := BulkCreateProductsResponse{
		Success:  true,
		Message:  "Successfully created 1 products",
		Created:  1,
		Products: products,
	}

	mph.service.On("BulkCreateProducts", mock.Anything, products).Return(expectedResponse, nil)

	requestBodyBytes, _ := json.Marshal(requestBody)
	mph.server.PerformRequest("/products/bulk-create", "post", requestBodyBytes)
	var actualResponse BulkCreateProductsResponse
	json.NewDecoder(mph.server.Recorder().Body).Decode(&actualResponse)

	assert.Equal(mph.T(), http.StatusOK, mph.server.Recorder().Code)
	assert.Equal(mph.T(), expectedResponse.Success, actualResponse.Success)
	assert.Equal(mph.T(), expectedResponse.Created, actualResponse.Created)
}

func (mph *ProductUploadHandlerTestSuite) TestShouldReturnErrorWhenProductsArrayIsEmpty() {
	requestBody := BulkCreateProductsRequest{Products: []Product{}}
	expectedResponse := types.ErrorResponse{
		Error: types.Error{
			Message:        "Products array cannot be empty",
			Code:           "error_processing_request",
			DisplayMessage: "Products array cannot be empty",
			Status:         "error",
		},
	}

	requestBodyBytes, _ := json.Marshal(requestBody)
	mph.server.PerformRequest("/products/bulk-create", "post", requestBodyBytes)
	var actualResponse types.ErrorResponse
	json.NewDecoder(mph.server.Recorder().Body).Decode(&actualResponse)

	assert.Equal(mph.T(), http.StatusBadRequest, mph.server.Recorder().Code)
	assert.Equal(mph.T(), expectedResponse.Error.Message, actualResponse.Error.Message)
	mph.service.AssertNotCalled(mph.T(), "BulkCreateProducts", mock.Anything, mock.Anything)
}

func (mph *ProductUploadHandlerTestSuite) TestShouldReturnErrorWhenExpectedErrorIsThrown() {
	products := []Product{
		{
			Name:        "Titan Edge 1",
			Category:    "watch",
			Brand:       "titan",
			Price:       12999,
			Description: "Titan Edge Slim Series",
			Images:      []string{"https://cdn.example.com/titan1.png"},
			Inventory:   20,
		},
	}
	requestBody := BulkCreateProductsRequest{Products: products}
	expectedResponse := types.ErrorResponse{
		Error: types.Error{
			Message:        "Validation failed",
			Code:           "error_processing_request",
			DisplayMessage: "Validation failed",
			Status:         "error",
		},
	}

	mph.service.On("BulkCreateProducts", mock.Anything, products).Return(BulkCreateProductsResponse{}, types.NewValidationError("Validation failed"))

	requestBodyBytes, _ := json.Marshal(requestBody)
	mph.server.PerformRequest("/products/bulk-create", "post", requestBodyBytes)
	var actualResponse types.ErrorResponse
	json.NewDecoder(mph.server.Recorder().Body).Decode(&actualResponse)

	assert.Equal(mph.T(), http.StatusBadRequest, mph.server.Recorder().Code)
	assert.Equal(mph.T(), expectedResponse.Error.Message, actualResponse.Error.Message)
}

func (mph *ProductUploadHandlerTestSuite) TestShouldReturnServerErrorWhenRandomErrorIsThrown() {
	products := []Product{
		{
			Name:        "Titan Edge 1",
			Category:    "watch",
			Brand:       "titan",
			Price:       12999,
			Description: "Titan Edge Slim Series",
			Images:      []string{"https://cdn.example.com/titan1.png"},
			Inventory:   20,
		},
	}
	requestBody := BulkCreateProductsRequest{Products: products}
	expectedResponse := types.ErrorResponse{
		Error: types.Error{
			Message:        "An unknown error has occurred",
			Code:           "internal_server_error",
			DisplayMessage: "An unknown error has occurred",
			Status:         "error",
		},
	}

	mph.service.On("BulkCreateProducts", mock.Anything, products).Return(BulkCreateProductsResponse{}, errors.New("random error"))

	requestBodyBytes, _ := json.Marshal(requestBody)
	mph.server.PerformRequest("/products/bulk-create", "post", requestBodyBytes)
	var actualResponse types.ErrorResponse
	json.NewDecoder(mph.server.Recorder().Body).Decode(&actualResponse)

	assert.Equal(mph.T(), http.StatusInternalServerError, mph.server.Recorder().Code)
	assert.Equal(mph.T(), expectedResponse.Error.Message, actualResponse.Error.Message)
}

func TestProductUploadHandlerTest(t *testing.T) {
	suite.Run(t, new(ProductUploadHandlerTestSuite))
}
