package product

import (
	"context"
	"fmt"

	"github.com/roppenlabs/rapid-product-catalog/internal/types"
	"github.com/roppenlabs/rapid-product-catalog/internal/utils"
	logger "github.com/roppenlabs/rapido-logger-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateProducts(ctx context.Context, products []Product) ([]Product, error)
}

type repositoryImpl struct {
	collection *mongo.Collection
}

func NewRepository(db *utils.DBInstance) Repository {
	if db == nil || db.TestDB == nil {
		panic("database cannot be nil")
	}
	return &repositoryImpl{
		collection: db.TestDB.Collection("rapidProducts"),
	}
}

func (r *repositoryImpl) CreateProducts(ctx context.Context, products []Product) ([]Product, error) {
	updatedProducts := make([]Product, 0, len(products))
	upsertOptions := options.Update().SetUpsert(true)

	for _, product := range products {
		// Filter by name and category to find existing product
		filter := bson.M{
			"name":     product.Name,
			"category": product.Category,
		}

		// Update document - update all fields except ID (preserve existing ID)
		update := bson.M{
			"$set": bson.M{
				"name":         product.Name,
				"category":     product.Category,
				"brand":        product.Brand,
				"price":        product.Price,
				"description":  product.Description,
				"images":       product.Images,
				"availableQty": product.Inventory,
				"popularity":   product.Popularity,
			},
		}

		_, err := r.collection.UpdateOne(ctx, filter, update, upsertOptions)
		if err != nil {
			logger.Error(logger.Format{
				Message: fmt.Sprintf("Error upserting product %s", product.Name),
				Data: map[string]string{
					"productName": product.Name,
					"error":       err.Error(),
				},
			})
			return nil, types.NewInternalServerError()
		}

		// Fetch the updated/inserted product from database to ensure response reflects DB state
		var updatedProduct Product
		err = r.collection.FindOne(ctx, filter).Decode(&updatedProduct)
		if err != nil {
			logger.Error(logger.Format{
				Message: fmt.Sprintf("Error finding product %s after upsert", product.Name),
				Data: map[string]string{
					"productName": product.Name,
					"error":       err.Error(),
				},
			})
			return nil, types.NewInternalServerError()
		}

		updatedProducts = append(updatedProducts, updatedProduct)
	}

	return updatedProducts, nil
}
