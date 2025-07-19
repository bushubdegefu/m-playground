package services

import (
	"github.com/bushubdegefu/m-playground/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

var AppCacheService *cache.CacheService

func InitServices(client *mongo.Client) {
	var err error
	AppCacheService, err = cache.NewCacheService()
	if err != nil {
		panic("Unable to initialize cache service")
	}

	// Initialize services here
	NewUserService(client)
	NewGroupService(client)
	NewPermissionService(client)
}
