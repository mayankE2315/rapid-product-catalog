package product

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.POST("/products/bulk", h.BulkCreateProductsHandler)
}
