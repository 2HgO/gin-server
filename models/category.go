package models

import "github.com/2HgO/gin-web-server/configs/db"

type Category struct {
	db.DocBase `json:",inline" bson:",inline"`
	Name       string `json:"name" bson:"name" binding:"required"`
	Icon       string `json:"icon" bson:"icon" binding:"required,url"`
}
