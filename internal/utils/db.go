package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/roppenlabs/rapid-product-catalog/internal/config"
	logger "github.com/roppenlabs/rapido-logger-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DBInstance contains the MongoDB database instance
type DBInstance struct {
	TestDB *mongo.Database
}

func NewDBInstance(conf config.Config) (*DBInstance, error) {
	mongoConf := conf.Get().Datastores.TestDB

	testDB, err := initDB(mongoConf)
	if err != nil {
		logger.Error(logger.Format{
			Message: "Failed to initialize mongo instance",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
		return nil, err
	}

	return &DBInstance{
		TestDB: testDB,
	}, nil
}

func initDB(conf config.MongoDB) (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=%s&readPreference=secondaryPreferred&replicaSet=%s&appName=%s",
		conf.User, conf.Password, conf.Hosts, conf.Database, conf.AuthSource, conf.ReplicaSet, conf.AppName)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(uint64(conf.Options.MaxPoolSize)).
		SetMinPoolSize(uint64(conf.Options.MinPoolSize)).
		SetMaxConnIdleTime(time.Duration(conf.Options.IdleTimeout) * time.Second)

	maxTimeOut := time.Duration(conf.Options.ConnectionTimeout)
	if maxTimeOut == 0 {
		maxTimeOut = 5 // default to 5 seconds
	}
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeOut*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	logger.Info(logger.Format{Message: "DB connected"})

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	logger.Info(logger.Format{Message: "DB pinged"})
	return client.Database(conf.Database), nil
}

// Close closes the MongoDB connection
func (db *DBInstance) Close(ctx context.Context) {
	logger.Info(logger.Format{
		Message: "Closing mongo connection",
	})

	if err := db.TestDB.Client().Disconnect(ctx); err != nil {
		logger.Error(logger.Format{
			Message: "Failed to close testdb connection",
			Data: map[string]string{
				"error": err.Error(),
			},
		})
	}
}
