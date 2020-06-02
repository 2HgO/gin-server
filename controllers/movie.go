package controllers

import (
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type movieOut struct {
	Movies []*models.AggMovie `bson:"movies"`
	Count  int64              `bson:"count"`
}

func CreateMovie(movie *models.Movie) error {
	return Movie.Insert(movie)
}

func DeleteMovie(id *primitive.ObjectID) (*models.Movie, error) {
	movie := new(models.Movie)
	err := Movie.DeleteOne(obj{"_id": id}, movie)
	return movie, err
}

func GetMovies(page, limit uint) ([]*models.AggMovie, int64, error) {
	iter, err := Movie.Aggregate([]obj{
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "categories",
				"foreignField": "_id",
				"as":           "categories",
			},
		},
		{
			"$facet": obj{
				"movies": []obj{
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

	res := &movieOut{
		Count:  0,
		Movies: make([]*models.AggMovie, 0),
	}
	if !iter.Next(nil) {
		return res.Movies, res.Count, nil
	}

	err = iter.Decode(res)
	return res.Movies, res.Count, err
}

func GetMovie(id *primitive.ObjectID) (*models.AggMovie, error) {
	movie := new(models.AggMovie)
	iter, err := Movie.Aggregate([]obj{
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "categories",
				"foreignField": "_id",
				"as":           "categories",
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

	err = iter.Decode(movie)
	return movie, err
}

func GetMoviesByCategory(category *primitive.ObjectID, page, limit uint) ([]*models.AggMovie, int64, error) {
	iter, err := Movie.Aggregate([]obj{
		{
			"$match": obj{
				"categories": category,
			},
		},
		{
			"$lookup": obj{
				"from":         "categories",
				"localField":   "categories",
				"foreignField": "_id",
				"as":           "categories",
			},
		},
		{
			"$facet": obj{
				"movies": []obj{
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

	res := &movieOut{
		Count:  0,
		Movies: make([]*models.AggMovie, 0),
	}
	if !iter.Next(nil) {
		return res.Movies, res.Count, nil
	}

	err = iter.Decode(res)
	return res.Movies, res.Count, err
}

func SearchMovies(query string, page, limit uint) ([]*models.AggMovie, int64, error) {
	return nil, 0, errors.NotImplementedError
}
