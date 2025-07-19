package services

import "go.mongodb.org/mongo-driver/mongo"

func InitServices(client *mongo.Client) {
	// Initialize services here
	NewUserService(client)
	NewUserService(client)
	NewUserService(client)
	NewGroupService(client)
	NewGroupService(client)
	NewGroupService(client)
	NewPermissionService(client)
	NewPermissionService(client)
	NewPermissionService(client)
}
