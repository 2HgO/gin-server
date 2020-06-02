package controllers

import (
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userOut struct {
	Count int64             `bson:"count"`
	Users []*models.AggUser `bson:"users"`
}

func CreateUser(user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return User.Insert(user)
}

func UpdateUser(id *primitive.ObjectID, user *models.User) (*models.User, error) {
	update := structs.Map(user)
	user = new(models.User)

	err := User.PartialUpdateOne(obj{"_id": id}, update, user)
	return user, err
}

func DeleteUser(id *primitive.ObjectID) (*models.User, error) {
	user := new(models.User)
	err := User.DeleteOne(obj{"_id": id}, user)
	return user, err
}

func GetUsers(page, limit uint) ([]*models.AggUser, int64, error) {
	iter, err := User.Aggregate([]obj{
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "likes",
				"foreignField": "_id",
				"as":           "likes",
			},
		},
		{
			"$facet": obj{
				"users": []obj{
					{
						"$skip": (page - 1) * limit,
					},
					{
						"$limit": limit,
					},
				},
				"count": []obj{
					{
						"$count": "total",
					},
				},
			},
		},
		{
			"$unwind": "$count",
		},
		{
			"$addFields": obj{
				"count": "$count.total",
			},
		},
	})

	if err != nil {
		return nil, 0, err
	}

	res := &userOut{
		Count: 0,
		Users: make([]*models.AggUser, 0),
	}
	if !iter.Next(nil) {
		return res.Users, res.Count, nil
	}

	err = iter.Decode(res)
	return res.Users, res.Count, err
}

func GetUser(id *primitive.ObjectID) (*models.AggUser, error) {
	user := new(models.AggUser)
	iter, err := User.Aggregate([]obj{
		{
			"$match": obj{
				"_id": id,
			},
		},
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "likes",
				"foreignField": "_id",
				"as":           "likes",
			},
		},
		{
			"$limit": 1,
		},
	})
	if err != nil {
		return nil, err
	}
	if !iter.Next(nil) {
		return nil, errors.NotFoundError
	}

	err = iter.Decode(user)

	return user, err
}

func SearchUsers(query string, page, limit uint) ([]*models.AggUser, int64, error) {
	return nil, 0, errors.NotImplementedError
}
