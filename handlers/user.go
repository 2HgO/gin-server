package handlers

import (
	"net/http"

	"github.com/2HgO/gin-web-server/controllers"
	"github.com/2HgO/gin-web-server/errors"
	"github.com/2HgO/gin-web-server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c *gin.Context) {
	user := new(models.User)
	errors.CheckError(c.ShouldBindJSON(user))

	errors.CheckError(controllers.CreateUser(user))

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User created successfully", "data": user})
}

func UpdateUser(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)

	update := new(models.User)
	update.SetID(*user.GetID())
	errors.CheckError(c.ShouldBindJSON(update))

	updatedUser, err := controllers.UpdateUser(user.GetID(), update)
	errors.CheckError(err)

	updatedUser.Password = ""
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User updated successfully", "data": updatedUser})
}

func DeleteUser(c *gin.Context) {
	user := c.MustGet("user").(*models.AggUser)

	_, err := controllers.DeleteUser(user.GetID())
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User deleted successfully"})
}

func GetUsers(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	users, count, err := controllers.GetUsers(pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Users retrieved successfully", "count": count, "data": users})
}

func GetUser(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	userID, err := primitive.ObjectIDFromHex(c.Param("userID"))
	errors.CheckError(err)

	user, err := controllers.GetUser(&userID)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User retrieved successfully", "data": user})
}

func SearchUsers(c *gin.Context) {
	_ = c.MustGet("user").(*models.AggUser)

	pagination := new(QueryParams)
	errors.CheckError(c.ShouldBindQuery(pagination))

	users, count, err := controllers.SearchUsers(pagination.Query, pagination.Page, pagination.Limit)
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Users retrieved successfully", "count": count, "data": users})
}
