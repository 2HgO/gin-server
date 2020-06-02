package controllers

import (
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type categoryOut struct {
	Categories []*models.Category `bson:"categories"`
	Count      int64              `bson:"count"`
}

func CreateCategory(category *models.Category) error { return Category.Insert(category) }

func DeleteCategory(id *primitive.ObjectID) (*models.Category, error) {
	category := new(models.Category)
	err := Category.DeleteOne(obj{"_id": id}, category)
	return category, err
}

func GetCategories(page, limit uint) ([]*models.Category, int64, error) {
	iter, err := Category.Aggregate([]obj{
		{
			"$facet": obj{
				"categories": []obj{
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
	res := &categoryOut{
		Count:      0,
		Categories: make([]*models.Category, 0),
	}
	if !iter.Next(nil) {
		return res.Categories, res.Count, nil
	}
	err = iter.Decode(res)
	return res.Categories, res.Count, err
}

func GetCategory(id *primitive.ObjectID) (*models.Category, error) {
	category := new(models.Category)
	err := Category.FindByID(id, category)
	return category, err
}

func SearchCategories(query string, page, limit uint) ([]*models.Category, int64, error) {
	return nil, 0, errors.NotImplementedError
}
