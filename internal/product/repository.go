package product

import (
	"context"

	"github.com/roppenlabs/rapid-product-catalog/internal/types"
	"github.com/roppenlabs/rapid-product-catalog/internal/utils"
	logger "github.com/roppenlabs/rapido-logger-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateProducts(ctx context.Context, products []Product) (*CreateProductsResult, error)
	SearchProducts(ctx context.Context, categories []string, brands []string, minPrice, maxPrice *float64, searchText string, limit int) ([]Product, error)
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

func (r *repositoryImpl) CreateProducts(ctx context.Context, products []Product) (*CreateProductsResult, error) {
	if len(products) == 0 {
		return &CreateProductsResult{
			Created:    0,
			Updated:    0,
			ProductIDs: []primitive.ObjectID{},
		}, nil
	}

	// Build bulk write operations
	models := make([]mongo.WriteModel, 0, len(products))
	productFilters := make([]bson.M, 0, len(products))

	for _, product := range products {
		// Filter by name and category to find existing product
		filter := bson.M{
			"name":     product.Name,
			"category": product.Category,
		}
		productFilters = append(productFilters, filter)

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

		// Create UpdateOneModel with upsert option
		updateModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		models = append(models, updateModel)
	}

	// Execute bulk write operation
	bulkWriteOptions := options.BulkWrite().SetOrdered(false)
	bulkResult, err := r.collection.BulkWrite(ctx, models, bulkWriteOptions)
	if err != nil {
		logger.Error(logger.Format{
			Message: "Error executing bulk write for products",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
		return nil, types.NewInternalServerError()
	}

	// Fetch all products that were inserted/updated to get their IDs
	var updatedProducts []Product
	if len(productFilters) > 0 {
		// Build a query to fetch all products using $or with all filters
		findFilter := bson.M{
			"$or": productFilters,
		}

		cursor, err := r.collection.Find(ctx, findFilter)
		if err != nil {
			logger.Error(logger.Format{
				Message: "Error fetching products after bulk write",
				Data: map[string]string{
					"error": err.Error(),
				},
			})
			return nil, types.NewInternalServerError()
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &updatedProducts); err != nil {
			logger.Error(logger.Format{
				Message: "Error decoding products after bulk write",
				Data: map[string]string{
					"error": err.Error(),
				},
			})
			return nil, types.NewInternalServerError()
		}
	}

	// Extract product IDs
	productIDs := make([]primitive.ObjectID, 0, len(updatedProducts))
	for _, product := range updatedProducts {
		if !product.ID.IsZero() {
			productIDs = append(productIDs, product.ID)
		}
	}

	// Calculate created and updated counts
	// For upsert operations:
	// - UpsertedCount: documents that were inserted (new)
	// - MatchedCount: documents that matched the filter (existing)
	// - ModifiedCount: documents that were actually modified (subset of matched)
	created := int(bulkResult.UpsertedCount)
	// All matched documents are considered updated (even if values didn't change)
	updated := int(bulkResult.MatchedCount)

	return &CreateProductsResult{
		Created:    created,
		Updated:    updated,
		ProductIDs: productIDs,
	}, nil
}

func (r *repositoryImpl) SearchProducts(ctx context.Context, categories []string, brands []string, minPrice, maxPrice *float64, searchText string, limit int) ([]Product, error) {
	filter := bson.M{}

	// Category filter - support multiple categories
	if len(categories) > 0 {
		filter["category"] = bson.M{"$in": categories}
	}

	// Brand filter - support multiple brands
	if len(brands) > 0 {
		filter["brand"] = bson.M{"$in": brands}
	}

	// Price range filter
	if minPrice != nil || maxPrice != nil {
		priceFilter := bson.M{}
		if minPrice != nil {
			priceFilter["$gte"] = *minPrice
		}
		if maxPrice != nil {
			priceFilter["$lte"] = *maxPrice
		}
		filter["price"] = priceFilter
	}

	// Text search - search in name and description
	if searchText != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": searchText, "$options": "i"}},
			{"description": bson.M{"$regex": searchText, "$options": "i"}},
		}
	}

	// Sort by popularity (descending) and limit results
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "popularity", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		logger.Error(logger.Format{
			Message: "Error searching products",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
		return nil, types.NewInternalServerError()
	}
	defer cursor.Close(ctx)

	var products []Product
	if err := cursor.All(ctx, &products); err != nil {
		logger.Error(logger.Format{
			Message: "Error decoding products",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
		return nil, types.NewInternalServerError()
	}

	return products, nil
}
