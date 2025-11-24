package product

import "go.mongodb.org/mongo-driver/bson/primitive"

type BulkCreateProductsRequest struct {
	Products []Product `json:"products" binding:"required"`
}

type Product struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Category    string             `json:"category" binding:"required" bson:"category"`
	Brand       string             `json:"brand" binding:"required" bson:"brand"`
	Price       float64            `json:"price" binding:"required,gt=0" bson:"price"`
	Description string             `json:"description" binding:"required" bson:"description"`
	Images      []string           `json:"images" binding:"required" bson:"images"`
	Inventory   int                `json:"inventory" binding:"required,min=0" bson:"availableQty"`
	Popularity  float64            `json:"popularity" binding:"required" bson:"popularity"`
}

type CreateProductsResponse struct {
	Success    bool                 `json:"success"`
	Message    string               `json:"message"`
	Created    int                  `json:"created"`
	Updated    int                  `json:"updated"`
	ProductIDs []primitive.ObjectID `json:"productIds,omitempty"`
}

type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type SearchProductsRequest struct {
	Category   interface{} `json:"category"` // Can be string or []string
	Brand      interface{} `json:"brand"`    // Can be string or []string
	PriceRange *PriceRange `json:"priceRange"`
	Search     string      `json:"search"`
}

// SearchParams is the normalized internal representation used by the service
type SearchParams struct {
	Categories []string
	Brands     []string
	MinPrice   *float64
	MaxPrice   *float64
	SearchText string
	Limit      int
}

type SearchProductsResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Count    int       `json:"count"`
	Products []Product `json:"products"`
}

type CreateProductsResult struct {
	Created    int
	Updated    int
	ProductIDs []primitive.ObjectID
}
