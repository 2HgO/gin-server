package handlers

import (
	"net/http"

	"github.com/2HgO/gin-web-server/controllers"
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddMovieToUserWatchlist(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)

	movieID, err := primitive.ObjectIDFromHex(c.Param("movieID"))
	errors.CheckError(err)

	errors.CheckError(controllers.AddMovieToUserWatchlist(user.GetID(), &movieID))

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "movie added successfully"})
}

func GetWatchlist(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)

	watchlist, err := controllers.GetUserWatchlist(user.GetID())
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Watchlist retrieved successfully", "data": watchlist})
}

func GetUserWatchlist(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	userID, err := primitive.ObjectIDFromHex(c.Param("userID"))
	errors.CheckError(err)

	watchlist, err := controllers.GetUserWatchlist(&userID)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Watchlist retrieved successfully", "data": watchlist})
}

func RemoveMovieFromUserWatchlist(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)

	movieID, err := primitive.ObjectIDFromHex(c.Param("movieID"))
	errors.CheckError(err)

	errors.CheckError(controllers.RemoveMovieFromUserWatchlist(user.GetID(), &movieID))

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "movie removed successfully"})
}
