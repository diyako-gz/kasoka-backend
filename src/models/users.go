package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	LastName  string             `json:"lastname" bson:"lastname"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Phone     string             `json:"phone" bson:"phone"`
	Role      string             `json:"role" bson:"role"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated-at" bson:"updated-at"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)



type CreateUserHeader struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Username string `json:"username" form:"username"`
	Name     string `json:"name" form:"name"`
	LastName string `json:"lastname" form:"lastname"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	UserName string `json:"username"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type LogInWhitPass struct {
	UserName string             `json:"username" form:"username"`
	Password string             `json:"password" form:"password"`
	ID       primitive.ObjectID `json:"id" form:"id"`
}

type LogInWhitId struct {
	ID primitive.ObjectID `json:"id" form:"id"`
}

type UpdatePass struct {
	Password string `json:"password" form:"password"`
}

type UpdateUserName struct {
	UserName string `json:"username" form:"username"`
}

type Bmi struct {
	Weight float64 `json:"weight" form:"weight"`
	Height float64 `json:"height" form:"height"`
}
