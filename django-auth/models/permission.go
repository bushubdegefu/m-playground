package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Permission Database model info
// @Description App type information
type Permission struct {
	ID       primitive.ObjectID `bson:"_id,omitzero" json:"id,omitzero"`
	Name     string             `bson:"name,omitzero" json:"name,omitzero"`
	Codename time.Time          `bson:"codename,omitzero" json:"codename,omitzero"`

	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

// PermissionPost model info
// @Description PermissionPost type information
type PermissionPost struct {
	Name string `bson:"name,omitzero" json:"name,omitzero"`
}

// PermissionGet model info
// @Description PermissionGet type information
type PermissionGet struct {
	ID primitive.ObjectID `bson:"_id,omitzero" json:"id,omitzero"`

	Codename time.Time `bson:"codename,omitzero" json:"codename,omitzero"`

	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

// PermissionPut model info
// @Description PermissionPut type information
type PermissionPut struct {
	Name *string `bson:"name,omitzero" json:"name,omitzero"`
}

// PermissionPatch model info
// @Description PermissionPatch type information
type PermissionPatch struct {
	Name *string `bson:"name,omitzero" json:"name,omitzero"`
}
