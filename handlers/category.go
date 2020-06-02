package handlers

import (
	"net/http"

	"github.com/2HgO/gin-web-server/controllers"
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)
	if user.Role != models.ADMIN {
		panic(errors.PermissionError)
	}

	category := new(models.Category)
	errors.CheckError(c.ShouldBindJSON(category))

	errors.CheckError(controllers.CreateCategory(category))

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "category created succesfully", "data": category})
}

func DeleteCategory(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)
	if user.Role != models.ADMIN {
		panic(errors.PermissionError)
	}

	categoryID, err := primitive.ObjectIDFromHex(c.Param("categoryID"))
	errors.CheckError(err)

	_, err = controllers.DeleteCategory(&categoryID)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category deleted succesfully"})
}

func GetCategory(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	categoryID, err := primitive.ObjectIDFromHex(c.Param("categoryID"))
	errors.CheckError(err)

	category, err := controllers.GetCategory(&categoryID)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category retrieved succesfully", "data": category})

}

func GetCategories(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	categories, count, err := controllers.GetCategories(pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category retrieved succesfully", "count": count, "data": categories})
}

func SearchCategories(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	categories, count, err := controllers.SearchCategories(pagination.Query, pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "category retrieved succesfully", "count": count, "data": categories})
}
