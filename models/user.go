package models

import (
	"strings"

	"github.com/2HgO/gin-web-server/configs/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	db.DocBase `json:",inline" bson:",inline" structs:"-"`
	FirstName  string                `json:"firstName" bson:"firstName" binding:"required_without=ID" structs:"firstName,omitempty"`
	LastName   string                `json:"lastName" bson:"lastName" binding:"required_without=ID" structs:"lastName,omitempty"`
	Email      string                `json:"email" bson:"email" binding:"required_without=ID,omitempty,email" structs:"-"`
	Password   string                `json:"password,omitempty" bson:"password" binding:"required_without=ID,omitempty,min=6" structs:"-"`
	DOB        *Date                 `json:"dob,omitempty" bson:"dob,omitmepty" structs:"dob,omitempty,omitnested"`
	Role       Role                  `json:"role" bson:"role" structs:"-"`
	Likes      []*primitive.ObjectID `json:"likes" bson:"likes" binding:"omitempty,dive,category" structs:"likes,omitempty,omitnested"`
}

type AggUser struct {
	db.DocBase `json:",inline" bson:",inline"`
	FirstName  string      `json:"firstName" bson:"firstName"`
	LastName   string      `json:"lastName" bson:"lastName"`
	Email      string      `json:"email" bson:"email"`
	Password   string      `json:"-" bson:"password"`
	Role       Role        `json:"role" bson:"role"`
	DOB        *Date       `json:"dob,omitempty" bson:"dob,omitmepty"`
	Likes      []*Category `json:"likes" bson:"likes"`
}

func (u *User) PreSave(*db.Collection) error {
	u.Email = strings.ToLower(u.Email)
	return nil
}
