package models

import (
	"github.com/2HgO/gin-web-server/configs/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	db.DocBase `json:",inline" bson:",inline"`
	Name       string                `json:"name" bson:"name" binding:"required"`
	Categories []*primitive.ObjectID `json:"categories" bson:"categories" binding:"required,dive,category"`
	Release    *Date                 `json:"release,omitempty" bson:"release,omitempty" binding:"required"`
}

type AggMovie struct {
	db.DocBase `json:",inline" bson:",inline"`
	Name       string      `json:"name" bson:"name"`
	Categories []*Category `json:"categories" bson:"categories"`
	Release    *Date       `json:"release,omitempty" bson:"release,omitempty"`
}
