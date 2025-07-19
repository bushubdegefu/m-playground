package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bushubdegefu/m-playground/django-auth/models"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var HandlerPermissionService PermissionService

// UserService wraps MongoDB logic for users
type PermissionService struct {
	Collection *mongo.Collection
	Client     *mongo.Client
	Database   *mongo.Database
}

// Constructor For Client
func NewPermissionService(client *mongo.Client) (*PermissionService, error) {
	database := client.Database("django_auth")
	collection := database.Collection("Permissions")
	HandlerPermissionService = PermissionService{
		Collection: collection,
		Client:     client,
		Database:   database,
	}
	return &HandlerPermissionService, nil
}

// Utility function for transactions
func (s *PermissionService) withTransaction(ctx context.Context, fn func(sc mongo.SessionContext) error) error {
	session, err := s.Client.StartSession()
	if err != nil {
		return fmt.Errorf("start session failed: %w", err)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		if err := fn(sc); err != nil {
			session.AbortTransaction(sc)
			return err
		}
		return session.CommitTransaction(sc)
	})
}

// Create inserts a new permission
func (s *PermissionService) Create(ctx context.Context, posted_permission *models.PermissionPost) (*models.PermissionGet, error) {
	var createdPermission models.PermissionGet

	err := s.withTransaction(ctx, func(sc mongo.SessionContext) error {

		permission := models.Permission{
			ID:        primitive.NewObjectID(),
			Name:      posted_permission.Name,
			CreatedAt: time.Now(),
		}

		_, err := s.Collection.InsertOne(ctx, permission)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		copier.Copy(createdPermission, permission)

		return nil
	})

	return &createdPermission, err
}

// GetOne fetches a permission by ID
func (s *PermissionService) GetOne(ctx context.Context, id string) (*models.PermissionGet, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	var permission models.PermissionGet
	err = s.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&permission)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// Get returns permissions with pagination and search
func (s *PermissionService) Get(ctx context.Context, pagination models.Pagination, searchFields []string, searchTerm []string) ([]models.PermissionGet, uint, error) {

	// Build search query if any
	filter := bson.M{}
	if len(searchTerm) > 0 && len(searchFields) > 0 && len(searchFields) >= len(searchTerm) {
		var orConditions []bson.M
		for index, term := range searchTerm {
			orConditions = append(orConditions, bson.M{
				searchFields[index]: bson.M{"$regex": term, "$options": "i"},
			})
		}
		filter["$or"] = orConditions
	}

	//pagination logic
	skip := int64(pagination.Page * pagination.Size)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(pagination.Size))

	// Count all documents (no filter)
	totalCount, _ := s.Collection.CountDocuments(ctx, filter)

	cursor, err := s.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, uint(totalCount), err
	}
	defer cursor.Close(ctx)

	var permissions []models.PermissionGet
	for cursor.Next(ctx) {
		var u models.PermissionGet
		if err := cursor.Decode(&u); err != nil {
			return nil, uint(totalCount), err
		}
		permissions = append(permissions, u)
	}

	return permissions, uint(totalCount), nil
}

// Update modifies a Permissions by ID
func (s *PermissionService) Update(ctx context.Context, patch_permission *models.PermissionPatch, id string) (*models.PermissionGet, error) {
	// update User
	var updatedPermission *models.PermissionGet

	permission_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.PermissionGet{}, fmt.Errorf("invalid ID: %w", err)
	}

	err = s.withTransaction(ctx, func(sc mongo.SessionContext) error {
		updateFields := bson.M{}
		if patch_permission.Name != nil {
			updateFields["name"] = *patch_permission.Name
		}
		updateFields["updated_at"] = time.Now()

		// filter to use to update value by
		filterPermission := bson.M{"_id": permission_id}
		updatePermission := bson.M{"$set": updateFields}
		// Update the document by ID
		_, err := s.Collection.UpdateOne(ctx, filterPermission, updatePermission)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		return nil
	})

	copier.Copy(&updatedPermission, patch_permission)
	return updatedPermission, err
}

// Delete removes a permission by ID
func (s *PermissionService) Delete(ctx context.Context, id string) error {

	err := s.withTransaction(ctx, func(sc mongo.SessionContext) error {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid ID: %w", err)
		}

		result, err := s.Collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			return err
		}

		if result.DeletedCount == 0 {
			return errors.New("no document deleted")
		}

		return nil
	})

	return err
}

// ##########################################################
// ##########  Custom Services Add Here   ###################
// ##########################################################
