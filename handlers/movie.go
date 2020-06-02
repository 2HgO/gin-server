package handlers

import (
	"net/http"

	"github.com/2HgO/gin-web-server/controllers"
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateMovie(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)
	if user.Role != models.ADMIN {
		panic(errors.PermissionError)
	}

	movie := new(models.Movie)
	errors.CheckError(c.ShouldBindJSON(movie))

	errors.CheckError(controllers.CreateMovie(movie))

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Movie created successfully", "data": movie})
}

func DeleteMovie(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)
	if user.Role != models.ADMIN {
		panic(errors.PermissionError)
	}

	movieID, err := primitive.ObjectIDFromHex(c.Param("movieID"))
	errors.CheckError(err)

	movie, err := controllers.DeleteMovie(&movieID)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movie deleted successfully", "data": movie})
}

func GetMovies(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	movies, count, err := controllers.GetMovies(pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movies retrieved successfully", "count": count, "data": movies})
}

func GetMovie(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	movieID, err := primitive.ObjectIDFromHex(c.Param("movieID"))
	errors.CheckError(err)

	movie, err := controllers.GetMovie(&movieID)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movie retrieved successfully", "data": movie})
}

func GetMoviesByCategory(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	categoryID, err := primitive.ObjectIDFromHex(c.Param("categoryID"))
	errors.CheckError(err)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	movies, count, err := controllers.GetMoviesByCategory(&categoryID, pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movies retrieved successfully", "count": count, "data": movies})
}

func SearchMovies(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	movies, count, err := controllers.SearchMovies(pagination.Query, pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movies retrieved successfully", "count": count, "data": movies})
}
