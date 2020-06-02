package handlers

import (
	"net/http"
	"time"

	"github.com/2HgO/gin-web-server/configs/env"
	"github.com/2HgO/gin-web-server/controllers"
	"github.com/2HgO/gin-web-server/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	login := &struct{
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"min=6"`
	}{}

	errors.CheckError(c.ShouldBindJSON(login))

	user, err := controllers.Login(login.Email, login.Password)
	errors.CheckError(err)

	claims := jwt.MapClaims{"user": user.GetID(), "exp": time.Now().Add(time.Hour*24)}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(env.JWT_SECRET))
	errors.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "login successful", "token": token, "data": user})
}

func ForgotPassword(c *gin.Context) {
	errors.CheckError(errors.PermissionError)
}

func ResetPassword(c *gin.Context) {
	errors.CheckError(errors.PermissionError)
}