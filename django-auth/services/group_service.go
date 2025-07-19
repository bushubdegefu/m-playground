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

var HandlerGroupService GroupService

// UserService wraps MongoDB logic for users
type GroupService struct {
	Collection *mongo.Collection
	Client     *mongo.Client
	Database   *mongo.Database
}

// Constructor For Client
func NewGroupService(client *mongo.Client) (*GroupService, error) {
	database := client.Database("django_auth")
	collection := database.Collection("Groups")
	HandlerGroupService = GroupService{
		Collection: collection,
		Client:     client,
		Database:   database,
	}
	return &HandlerGroupService, nil
}

// Utility function for transactions
func (s *GroupService) withTransaction(ctx context.Context, fn func(sc mongo.SessionContext) error) error {
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

// Create inserts a new group
func (s *GroupService) Create(ctx context.Context, posted_group *models.GroupPost) (*models.GroupGet, error) {
	var createdGroup models.GroupGet

	err := s.withTransaction(ctx, func(sc mongo.SessionContext) error {

		group := models.Group{
			ID:        primitive.NewObjectID(),
			Name:      posted_group.Name,
			CreatedAt: time.Now(),
		}

		_, err := s.Collection.InsertOne(ctx, group)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		copier.Copy(createdGroup, group)

		return nil
	})

	return &createdGroup, err
}

// GetOne fetches a group by ID
func (s *GroupService) GetOne(ctx context.Context, id string) (*models.GroupGet, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	var group models.GroupGet
	err = s.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// Get returns groups with pagination and search
func (s *GroupService) Get(ctx context.Context, pagination models.Pagination, searchFields []string, searchTerm []string) ([]models.GroupGet, uint, error) {

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

	var groups []models.GroupGet
	for cursor.Next(ctx) {
		var u models.GroupGet
		if err := cursor.Decode(&u); err != nil {
			return nil, uint(totalCount), err
		}
		groups = append(groups, u)
	}

	return groups, uint(totalCount), nil
}

// Update modifies a Groups by ID
func (s *GroupService) Update(ctx context.Context, patch_group *models.GroupPatch, id string) (*models.GroupGet, error) {
	// update User
	var updatedGroup *models.GroupGet

	group_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.GroupGet{}, fmt.Errorf("invalid ID: %w", err)
	}

	err = s.withTransaction(ctx, func(sc mongo.SessionContext) error {
		updateFields := bson.M{}
		if patch_group.Name != nil {
			updateFields["name"] = *patch_group.Name
		}
		updateFields["updated_at"] = time.Now()

		// filter to use to update value by
		filterGroup := bson.M{"_id": group_id}
		updateGroup := bson.M{"$set": updateFields}
		// Update the document by ID
		_, err := s.Collection.UpdateOne(ctx, filterGroup, updateGroup)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		return nil
	})

	copier.Copy(&updatedGroup, patch_group)
	return updatedGroup, err
}

// Delete removes a group by ID
func (s *GroupService) Delete(ctx context.Context, id string) error {

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
// ##########  Relationship  Services to Permission
// ##########################################################

func (s *GroupService) AddGroupToPermission(ctx context.Context, groupID, permissionID string) error {
	group_id, err := primitive.ObjectIDFromHex(groupID)
	permission_id, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(ctx, bson.M{"_id": group_id}, bson.M{
		"$addToSet": bson.M{"permission_ids": permission_id}, // Prevents duplicates
	})
	return err
}

func (s *GroupService) RemoveGroupFromPermission(ctx context.Context, groupID, permissionID string) error {
	group_id, err := primitive.ObjectIDFromHex(groupID)
	permission_id, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(ctx, bson.M{"_id": group_id}, bson.M{
		"$pull": bson.M{"permission_ids": permission_id},
	})
	return err
}

func (s *GroupService) GetGroupPermissions(ctx context.Context, groupID string, pagination models.Pagination) ([]models.Permission, uint, error) {
	group_id, _ := primitive.ObjectIDFromHex(groupID)
	var group models.Group
	if err := s.Collection.FindOne(ctx, bson.M{"_id": group_id}).Decode(&group); err != nil {
		return nil, 0, err
	}

	permissionCollection := s.Database.Collection("Permissions")
	filter := bson.M{"_id": bson.M{"$in": group.PermissionIDs}}
	opts := options.Find().
		SetSkip(int64(pagination.Page * pagination.Size)).
		SetLimit(int64(pagination.Size))

	total, _ := permissionCollection.CountDocuments(ctx, filter)

	cursor, err := permissionCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var permissions []models.Permission
	for cursor.Next(ctx) {
		var g models.Permission
		cursor.Decode(&g)
		permissions = append(permissions, g)
	}

	return permissions, uint(total), nil
}

// #########################
// No Pagination Services###
// #########################

func (s *GroupService) GetAllPermissionsForGroup(ctx context.Context, groupID string) ([]models.Permission, error) {
	group_id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %w", err)
	}

	var group models.Group
	if err := s.Collection.FindOne(ctx, bson.M{"_id": group_id}).Decode(&group); err != nil {
		return nil, fmt.Errorf("failed to fetch group: %w", err)
	}

	permissionCollection := s.Database.Collection("Permissions")
	filter := bson.M{"_id": bson.M{"$in": group.PermissionIDs}}

	cursor, err := permissionCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
	}
	defer cursor.Close(ctx)

	var permissions []models.Permission
	for cursor.Next(ctx) {
		var g models.Permission
		if err := cursor.Decode(&g); err != nil {
			return nil, fmt.Errorf("failed to decode permission: %w", err)
		}
		permissions = append(permissions, g)
	}

	return permissions, nil
}

func (s *GroupService) GetAllPermissionsgroupDoesNotHave(ctx context.Context, groupID string) ([]models.Permission, error) {
	group_id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %w", err)
	}

	var group models.Group
	if err := s.Collection.FindOne(ctx, bson.M{"_id": group_id}).Decode(&group); err != nil {
		return nil, fmt.Errorf("failed to fetch group: %w", err)
	}

	permissionCollection := s.Database.Collection("Permissions")

	filter := bson.M{}
	if len(group.PermissionIDs) > 0 {
		filter["_id"] = bson.M{"$nin": group.PermissionIDs}
	}

	cursor, err := permissionCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
	}
	defer cursor.Close(ctx)

	var permissions []models.Permission
	for cursor.Next(ctx) {
		var g models.Permission
		if err := cursor.Decode(&g); err != nil {
			return nil, fmt.Errorf("failed to decode permission: %w", err)
		}
		permissions = append(permissions, g)
	}

	return permissions, nil
}

// ##########################################################
// ##########  Custom Services Add Here   ###################
// ##########################################################
