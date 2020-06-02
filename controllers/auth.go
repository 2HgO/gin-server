package controllers

import (
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (*models.AggUser, error) {
	user := new(models.AggUser)
	iter, err := User.Aggregate([]obj{
		{
			"$match": obj{
				"email": email,
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

	if err := iter.Decode(user); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
