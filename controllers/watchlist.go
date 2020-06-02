package controllers

import (
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddMovieToUserWatchlist(user, movie *primitive.ObjectID) error {
	if exists, err := Watchlist.Exists(obj{"user": user}); exists && err == nil {
		return Watchlist.PartialUpdateOne(obj{"user": user}, obj{"$addToSet": obj{"movies": movie}}, nil)
	}
	return Watchlist.Insert(&models.Watchlist{
		User:   user,
		Movies: []*primitive.ObjectID{movie},
	})
}

func GetUserWatchlist(user *primitive.ObjectID) (*models.AggWatchlist, error) {
	watchlist := new(models.AggWatchlist)
	iter, err := Watchlist.Aggregate([]obj{
		{
			"$match": obj{
				"user": user,
			},
		},
		{
			"$lookup": obj{
				"from":         "users",
				"localField":   "user",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$lookup": obj{
				"from":         "movies",
				"localField":   "movies",
				"foreignField": "_id",
				"as":           "movies",
			},
		},
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "user.likes",
				"foreignField": "_id",
				"as":           "user.likes",
			},
		},
		{
			"$unwind": "$movies",
		},
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "movies.categories",
				"foreignField": "_id",
				"as":           "movies.categories",
			},
		},
		{
			"$group": obj{
				"_id":       "$_id",
				"createdAt": obj{"$last": "$createdAt"},
				"updatedAt": obj{"$last": "$updatedAt"},
				"user":      obj{"$last": "$user"},
				"movies":    obj{"$addToSet": "$movies"},
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

	err = iter.Decode(watchlist)
	return watchlist, err
}

func RemoveMovieFromUserWatchlist(user, movie *primitive.ObjectID) error {
	return Watchlist.PartialUpdateOne(obj{"user": user, "movies": movie}, obj{"$pull": obj{"movies": movie}}, nil)
}
