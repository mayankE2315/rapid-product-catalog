package product

import (
	"context"
	"fmt"
	"net/http"

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

func (h *Handler) BulkCreateProductsHandler(ctx *gin.Context) {
	var req BulkCreateProductsRequest

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
	ctx.JSON(http.StatusOK, response)
}

func buildErrorResponse(err *types.StatusError) types.ErrorResponse {
	return types.ErrorResponse{
		Error: types.Error{
			Message:        err.Message,
			Code:           err.Code,
			DisplayMessage: err.Message,
			Status:         "error",
		},
	}
}
