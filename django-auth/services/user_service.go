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

var HandlerUserService UserService

// UserService wraps MongoDB logic for users
type UserService struct {
	Collection *mongo.Collection
	Client     *mongo.Client
	Database   *mongo.Database
}

// Constructor For Client
func NewUserService(client *mongo.Client) (*UserService, error) {
	database := client.Database("django_auth")
	collection := database.Collection("Users")
	HandlerUserService = UserService{
		Collection: collection,
		Client:     client,
		Database:   database,
	}
	return &HandlerUserService, nil
}

// Utility function for transactions
func (s *UserService) withTransaction(ctx context.Context, fn func(sc mongo.SessionContext) error) error {
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

// Create inserts a new user
func (s *UserService) Create(ctx context.Context, posted_user *models.UserPost) (*models.UserGet, error) {
	var createdUser models.UserGet

	err := s.withTransaction(ctx, func(sc mongo.SessionContext) error {
		hashedPassword := models.HashFunc(posted_user.Password)

		user := models.User{
			ID:          primitive.NewObjectID(),
			Password:    hashedPassword,
			IsSuperuser: posted_user.IsSuperuser,
			Username:    posted_user.Username,
			FirstName:   posted_user.FirstName,
			LastName:    posted_user.LastName,
			Email:       posted_user.Email,
			IsStaff:     posted_user.IsStaff,
			IsActive:    posted_user.IsActive,
			CreatedAt:   time.Now(),
		}

		_, err := s.Collection.InsertOne(ctx, user)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		copier.Copy(createdUser, user)

		return nil
	})

	return &createdUser, err
}

// GetOne fetches a user by ID
func (s *UserService) GetOne(ctx context.Context, id string) (*models.UserGet, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	var user models.UserGet
	err = s.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get returns users with pagination and search
func (s *UserService) Get(ctx context.Context, pagination models.Pagination, searchFields []string, searchTerm []string) ([]models.UserGet, uint, error) {

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

	var users []models.UserGet
	for cursor.Next(ctx) {
		var u models.UserGet
		if err := cursor.Decode(&u); err != nil {
			return nil, uint(totalCount), err
		}
		users = append(users, u)
	}

	return users, uint(totalCount), nil
}

// Update modifies a Users by ID
func (s *UserService) Update(ctx context.Context, patch_user *models.UserPatch, id string) (*models.UserGet, error) {
	// update User
	var updatedUser *models.UserGet

	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.UserGet{}, fmt.Errorf("invalid ID: %w", err)
	}

	err = s.withTransaction(ctx, func(sc mongo.SessionContext) error {
		updateFields := bson.M{}
		if patch_user.Password != nil {
			// setting password string to hash
			hashedPassword := models.HashFunc(*patch_user.Password)
			updateFields["password"] = hashedPassword
		}
		if patch_user.IsSuperuser != nil {
			updateFields["is_superuser"] = *patch_user.IsSuperuser
		}
		if patch_user.Username != nil {
			updateFields["username"] = *patch_user.Username
		}
		if patch_user.FirstName != nil {
			updateFields["first_name"] = *patch_user.FirstName
		}
		if patch_user.LastName != nil {
			updateFields["last_name"] = *patch_user.LastName
		}
		if patch_user.Email != nil {
			updateFields["email"] = *patch_user.Email
		}
		if patch_user.IsStaff != nil {
			updateFields["is_staff"] = *patch_user.IsStaff
		}
		if patch_user.IsActive != nil {
			updateFields["is_active"] = *patch_user.IsActive
		}
		updateFields["updated_at"] = time.Now()

		// filter to use to update value by
		filterUser := bson.M{"_id": user_id}
		updateUser := bson.M{"$set": updateFields}
		// Update the document by ID
		_, err := s.Collection.UpdateOne(ctx, filterUser, updateUser)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}

		return nil
	})

	copier.Copy(&updatedUser, patch_user)
	return updatedUser, err
}

// Delete removes a user by ID
func (s *UserService) Delete(ctx context.Context, id string) error {

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

func (s *UserService) AddUserToPermission(ctx context.Context, userID, permissionID string) error {
	user_id, err := primitive.ObjectIDFromHex(userID)
	permission_id, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(ctx, bson.M{"_id": user_id}, bson.M{
		"$addToSet": bson.M{"permission_ids": permission_id}, // Prevents duplicates
	})
	return err
}

func (s *UserService) RemoveUserFromPermission(ctx context.Context, userID, permissionID string) error {
	user_id, err := primitive.ObjectIDFromHex(userID)
	permission_id, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(ctx, bson.M{"_id": user_id}, bson.M{
		"$pull": bson.M{"permission_ids": permission_id},
	})
	return err
}

func (s *UserService) GetUserPermissions(ctx context.Context, userID string, pagination models.Pagination) ([]models.Permission, uint, error) {
	user_id, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	if err := s.Collection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user); err != nil {
		return nil, 0, err
	}

	permissionCollection := s.Database.Collection("Permissions")
	filter := bson.M{"_id": bson.M{"$in": user.PermissionIDs}}
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

func (s *UserService) GetAllPermissionsForUser(ctx context.Context, userID string) ([]models.Permission, error) {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	if err := s.Collection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	permissionCollection := s.Database.Collection("Permissions")
	filter := bson.M{"_id": bson.M{"$in": user.PermissionIDs}}

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

func (s *UserService) GetAllPermissionsuserDoesNotHave(ctx context.Context, userID string) ([]models.Permission, error) {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	if err := s.Collection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	permissionCollection := s.Database.Collection("Permissions")

	filter := bson.M{}
	if len(user.PermissionIDs) > 0 {
		filter["_id"] = bson.M{"$nin": user.PermissionIDs}
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
// ##########  Relationship  Services to Group
// ##########################################################

func (s *UserService) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	user_id, err := primitive.ObjectIDFromHex(userID)
	group_id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(ctx, bson.M{"_id": user_id}, bson.M{
		"$addToSet": bson.M{"group_ids": group_id}, // Prevents duplicates
	})
	return err
}

func (s *UserService) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	user_id, err := primitive.ObjectIDFromHex(userID)
	group_id, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	_, err = s.Collection.UpdateOne(ctx, bson.M{"_id": user_id}, bson.M{
		"$pull": bson.M{"group_ids": group_id},
	})
	return err
}

func (s *UserService) GetUserGroups(ctx context.Context, userID string, pagination models.Pagination) ([]models.Group, uint, error) {
	user_id, _ := primitive.ObjectIDFromHex(userID)
	var user models.User
	if err := s.Collection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user); err != nil {
		return nil, 0, err
	}

	groupCollection := s.Database.Collection("Groups")
	filter := bson.M{"_id": bson.M{"$in": user.GroupIDs}}
	opts := options.Find().
		SetSkip(int64(pagination.Page * pagination.Size)).
		SetLimit(int64(pagination.Size))

	total, _ := groupCollection.CountDocuments(ctx, filter)

	cursor, err := groupCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var groups []models.Group
	for cursor.Next(ctx) {
		var g models.Group
		cursor.Decode(&g)
		groups = append(groups, g)
	}

	return groups, uint(total), nil
}

// #########################
// No Pagination Services###
// #########################

func (s *UserService) GetAllGroupsForUser(ctx context.Context, userID string) ([]models.Group, error) {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	if err := s.Collection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	groupCollection := s.Database.Collection("Groups")
	filter := bson.M{"_id": bson.M{"$in": user.GroupIDs}}

	cursor, err := groupCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer cursor.Close(ctx)

	var groups []models.Group
	for cursor.Next(ctx) {
		var g models.Group
		if err := cursor.Decode(&g); err != nil {
			return nil, fmt.Errorf("failed to decode group: %w", err)
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func (s *UserService) GetAllGroupsuserDoesNotHave(ctx context.Context, userID string) ([]models.Group, error) {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	if err := s.Collection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	groupCollection := s.Database.Collection("Groups")

	filter := bson.M{}
	if len(user.GroupIDs) > 0 {
		filter["_id"] = bson.M{"$nin": user.GroupIDs}
	}

	cursor, err := groupCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer cursor.Close(ctx)

	var groups []models.Group
	for cursor.Next(ctx) {
		var g models.Group
		if err := cursor.Decode(&g); err != nil {
			return nil, fmt.Errorf("failed to decode group: %w", err)
		}
		groups = append(groups, g)
	}

	return groups, nil
}

// ##########################################################
// ##########  Custom Services Add Here   ###################
// ##########################################################
