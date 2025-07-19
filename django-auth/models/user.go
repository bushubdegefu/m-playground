package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User Database model info
// @Description App type information
type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitzero" json:"id,omitzero"`
	Password      string               `bson:"password,omitzero" json:"password,omitzero"`
	LastLogin     time.Time            `bson:"last_login,omitzero"  json:"last_login,omitzero"`
	IsSuperuser   bool                 `bson:"is_superuser,omitzero" json:"is_superuser"`
	Username      string               `bson:"username,omitzero" json:"username,omitzero"`
	FirstName     string               `bson:"first_name,omitzero" json:"first_name"`
	LastName      string               `bson:"last_name,omitzero" json:"last_name"`
	Email         string               `bson:"email,omitzero" json:"email,omitzero"`
	IsStaff       bool                 `bson:"is_staff,omitzero" json:"is_staff"`
	IsActive      bool                 `bson:"is_active,omitzero" json:"is_active"`
	GroupIDs      []primitive.ObjectID `bson:"group_ids,omitempty"    json:"group_ids,omitzero"`
	PermissionIDs []primitive.ObjectID `bson:"permission_ids,omitempty" json:"permission_ids,omitempty"`

	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

// UserPost model info
// @Description UserPost type information
type UserPost struct {
	Password string `bson:"password,omitzero" json:"password,omitzero"`

	IsSuperuser bool   `bson:"is_superuser,omitzero" json:"is_superuser"`
	Username    string `bson:"username,omitzero" json:"username,omitzero"`
	FirstName   string `bson:"first_name,omitzero" json:"first_name"`
	LastName    string `bson:"last_name,omitzero" json:"last_name"`
	Email       string `bson:"email,omitzero" json:"email,omitzero"`
	IsStaff     bool   `bson:"is_staff,omitzero" json:"is_staff"`
	IsActive    bool   `bson:"is_active,omitzero" json:"is_active"`
}

// UserGet model info
// @Description UserGet type information
type UserGet struct {
	ID primitive.ObjectID `bson:"_id,omitzero" json:"id,omitzero"`

	LastLogin   time.Time `bson:"last_login,omitzero"  json:"last_login,omitzero"`
	IsSuperuser bool      `bson:"is_superuser,omitzero" json:"is_superuser"`
	Username    string    `bson:"username,omitzero" json:"username,omitzero"`
	FirstName   string    `bson:"first_name,omitzero" json:"first_name"`
	LastName    string    `bson:"last_name,omitzero" json:"last_name"`
	Email       string    `bson:"email,omitzero" json:"email,omitzero"`
	IsStaff     bool      `bson:"is_staff,omitzero" json:"is_staff"`
	IsActive    bool      `bson:"is_active,omitzero" json:"is_active"`

	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

// UserPut model info
// @Description UserPut type information
type UserPut struct {
	Password *string `bson:"password,omitzero" json:"password,omitzero"`

	IsSuperuser *bool   `bson:"is_superuser,omitzero" json:"is_superuser"`
	Username    *string `bson:"username,omitzero" json:"username,omitzero"`

	Email    *string               `bson:"email,omitzero" json:"email,omitzero"`
	IsStaff  *bool                 `bson:"is_staff,omitzero" json:"is_staff"`
	IsActive *bool                 `bson:"is_active,omitzero" json:"is_active"`
	GroupIDs *[]primitive.ObjectID `bson:"group_ids,omitempty"    json:"group_ids,omitzero"`
}

// UserPatch model info
// @Description UserPatch type information
type UserPatch struct {
	Password *string `bson:"password,omitzero" json:"password,omitzero"`

	IsSuperuser *bool   `bson:"is_superuser,omitzero" json:"is_superuser"`
	Username    *string `bson:"username,omitzero" json:"username,omitzero"`
	FirstName   *string `bson:"first_name,omitzero" json:"first_name"`
	LastName    *string `bson:"last_name,omitzero" json:"last_name"`
	Email       *string `bson:"email,omitzero" json:"email,omitzero"`
	IsStaff     *bool   `bson:"is_staff,omitzero" json:"is_staff"`
	IsActive    *bool   `bson:"is_active,omitzero" json:"is_active"`
}
