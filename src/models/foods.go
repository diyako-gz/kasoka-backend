package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Foods struct {
	Name   string             `json:"name" form:"name"`
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Price  string             `json:"price" form:"price"`
	Recipe string             `json:"recipe" form:"recipe"`
	Type   string             `json:"type" form:"type"`
	Rate   float64            `json:"rate" form:"rate"`
}

type FoodId struct {
	ID string `json:"_id,omitempty" bson:"_id,omitempty"`
}

type FoodType struct {
	Type string `json:"type" form:"type"`
}

type FoodRate struct {
	Rate float64  `json:"rate" form:"rate"`
}