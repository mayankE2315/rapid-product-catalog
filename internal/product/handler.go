package product

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/roppenlabs/rapid-product-catalog/internal/types"
	logger "github.com/roppenlabs/rapido-logger-go"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateProductsHandler(ctx *gin.Context) {
	var req BulkCreateProductsRequest

	logger.Info(logger.Format{Message: "Request received for bulk create products", Data: map[string]string{"request": fmt.Sprintf("%+v", req)}})

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(logger.Format{Message: fmt.Sprintf("Invalid request body: %v", err)})
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(types.NewValidationError(fmt.Sprintf("Invalid request: %v", err))))
		return
	}

	if len(req.Products) == 0 {
		logger.Error(logger.Format{Message: "Products array is empty"})
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(types.NewValidationError("Products array cannot be empty")))
		return
	}

	if validationErr := validateProducts(req.Products); validationErr != nil {
		logger.Error(logger.Format{Message: validationErr.Error()})
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(validationErr))
		return
	}

	response, err := h.service.BulkCreateProducts(context.Background(), req.Products)

	if err != nil {
		statusError, ok := err.(*types.StatusError)
		if !ok {
			serverError := types.NewInternalServerError()
			ctx.JSON(http.StatusInternalServerError, buildErrorResponse(serverError))
			return
		}
		ctx.JSON(statusError.HTTPCode, buildErrorResponse(statusError))
		return
	}
	logger.Info(logger.Format{Message: "Response for bulk create products", Data: map[string]string{"response": fmt.Sprintf("%+v", response)}})
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) SearchProductsHandler(ctx *gin.Context) {
	var req SearchProductsRequest

	logger.Info(logger.Format{Message: "Request received for search products", Data: map[string]string{"request": fmt.Sprintf("%+v", req)}})

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error(logger.Format{Message: fmt.Sprintf("Invalid request body: %v", err)})
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(types.NewValidationError(fmt.Sprintf("Invalid request: %v", err))))
		return
	}

	// Normalize request to internal format
	searchParams := normalizeSearchRequest(req)

	response, err := h.service.SearchProducts(context.Background(), searchParams)

	if err != nil {
		statusError, ok := err.(*types.StatusError)
		if !ok {
			serverError := types.NewInternalServerError()
			ctx.JSON(http.StatusInternalServerError, buildErrorResponse(serverError))
			return
		}
		ctx.JSON(statusError.HTTPCode, buildErrorResponse(statusError))
		return
	}

	logger.Info(logger.Format{Message: "Response for search products", Data: map[string]string{"response": fmt.Sprintf("%+v", response)}})
	ctx.JSON(http.StatusOK, response)
}

func validateProducts(products []Product) *types.StatusError {
	for i, product := range products {
		if strings.TrimSpace(product.Name) == "" {
			return types.NewValidationError(fmt.Sprintf("Product at index %d: name cannot be empty", i))
		}
		if strings.TrimSpace(product.Category) == "" {
			return types.NewValidationError(fmt.Sprintf("Product at index %d: category cannot be empty", i))
		}
		if strings.TrimSpace(product.Brand) == "" {
			return types.NewValidationError(fmt.Sprintf("Product at index %d: brand cannot be empty", i))
		}
		if product.Price <= 0 {
			return types.NewValidationError(fmt.Sprintf("Product at index %d: price must be greater than 0", i))
		}
	}
	return nil
}

func normalizeSearchRequest(req SearchProductsRequest) SearchParams {
	params := SearchParams{
		Categories: []string{},
		Brands:     []string{},
		SearchText: req.Search,
		Limit:      15,
	}

	// Convert category to []string
	if req.Category != nil {
		switch v := req.Category.(type) {
		case string:
			if v != "" {
				params.Categories = []string{v}
			}
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str != "" {
					params.Categories = append(params.Categories, str)
				}
			}
		case []string:
			params.Categories = v
		}
	}

	// Convert brand to []string
	if req.Brand != nil {
		switch v := req.Brand.(type) {
		case string:
			if v != "" {
				params.Brands = []string{v}
			}
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str != "" {
					params.Brands = append(params.Brands, str)
				}
			}
		case []string:
			params.Brands = v
		}
	}

	// Extract price range
	if req.PriceRange != nil {
		if req.PriceRange.Min > 0 {
			params.MinPrice = &req.PriceRange.Min
		}
		if req.PriceRange.Max > 0 {
			params.MaxPrice = &req.PriceRange.Max
		}
	}

	return params
}

func buildErrorResponse(err *types.StatusError) types.ErrorResponse {
	return types.ErrorResponse{
		Error: types.Error{
			Message: err.Message,
			Code:    err.Code,
			Status:  "error",
		},
	}
}
