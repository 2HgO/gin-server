package models

import (
	"github.com/2HgO/gin-web-server/configs/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Watchlist struct {
	db.DocBase `json:",inline" bson:",inline"`
	User       *primitive.ObjectID   `json:"user" bson:"user"`
	Movies     []*primitive.ObjectID `json:"movies" bson:"movies"`
}

type AggWatchlist struct {
	db.DocBase `json:",inline" bson:",inline"`
	User       *AggUser    `json:"user" bson:"user"`
	Movies     []*AggMovie `json:"movies" bson:"movies"`
}
