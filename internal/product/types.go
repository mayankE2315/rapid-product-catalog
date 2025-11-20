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
	Price       float64            `json:"price" binding:"required,min=0" bson:"price"`
	Description string             `json:"description" binding:"required" bson:"description"`
	Images      []string           `json:"images" binding:"required" bson:"images"`
	Inventory   int                `json:"inventory" binding:"required,min=0" bson:"availableQty"`
	Popularity  float64            `json:"popularity" binding:"required" bson:"popularity"`
}

type BulkCreateProductsResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Created  int       `json:"created"`
	Products []Product `json:"products,omitempty"`
}
