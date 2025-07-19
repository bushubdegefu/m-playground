package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group Database model info
// @Description App type information
type Group struct {
	ID            primitive.ObjectID   `bson:"_id,omitzero" json:"id,omitzero"`
	Name          string               `bson:"name,omitzero" json:"name,omitzero"`
	PermissionIDs []primitive.ObjectID `bson:"permission_ids,omitzero" json:"permission_ids,omitzero"`
	CreatedAt     time.Time            `bson:"created_at,omitempty"`
	UpdatedAt     time.Time            `bson:"updated_at,omitempty"`
}

// GroupPost model info
// @Description GroupPost type information
type GroupPost struct {
	Name string `bson:"name,omitzero" json:"name,omitzero"`
}

// GroupGet model info
// @Description GroupGet type information
type GroupGet struct {
	ID            primitive.ObjectID   `bson:"_id,omitzero" json:"id,omitzero"`
	Name          string               `bson:"name,omitzero" json:"name,omitzero"`
	PermissionIDs []primitive.ObjectID `bson:"permission_ids,omitzero" json:"permission_ids,omitzero"`
	CreatedAt     time.Time            `bson:"created_at,omitempty"`
	UpdatedAt     time.Time            `bson:"updated_at,omitempty"`
}

// GroupPut model info
// @Description GroupPut type information
type GroupPut struct {
	Name *string `bson:"name,omitzero" json:"name,omitzero"`
}

// GroupPatch model info
// @Description GroupPatch type information
type GroupPatch struct {
	Name *string `bson:"name,omitzero" json:"name,omitzero"`
}
