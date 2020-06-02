package handlers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/2HgO/gin-web-server/configs/env"
	"github.com/2HgO/gin-web-server/controllers"
	"github.com/2HgO/gin-web-server/errors"
)

var CORS = cors.New(
	cors.Config{ // Enable cors for white-listed origins
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		AllowWebSockets:  true,
		AllowFiles:       true,
		AllowMethods:     []string{"*"},
		AllowAllOrigins:  true,
		AllowWildcard:    true,
	},
)

func PanicHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var e error
			switch v := err.(type) {
			case error:
				e = v
			default:
				e = fmt.Errorf("%w: %v", errors.FatalError, v)
			}
			c.Error(e)
			code, message, errType := errors.ErrorHandler(e)
			c.AbortWithStatusJSON(code, gin.H{"success": false, "message": message, "error": errType})
		}
	}()
	c.Next()
}

func SetHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Next()
}

func Handle404(c *gin.Context) {
	c.AbortWithStatusJSON(404, gin.H{
		"success": false,
		"message": "API Endpoint not found",
		"error":   "INVALID_ROUTE_ERROR",
	})
}

func Handle405(c *gin.Context) {
	c.AbortWithStatusJSON(405, gin.H{
		"success": false,
		"message": "Method not allowed",
		"error":   "INVALID_METHOD_ERROR",
	})
}

func UserValidation(c *gin.Context) {
	if len(c.GetHeader("authorization")) < 7 {
		c.AbortWithStatusJSON(400, gin.H{"success": false, "message": "No auth token", "error": errors.ErrAuthentication})
		return
	}
	token := c.GetHeader("authorization")[7:]

	if token == "" {
		c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "Please login to continue.", "error": errors.ErrAuthentication})
		return
	}

	claims, err := ValidateToken(token)
	errors.CheckError(err)

	val, ok := claims["user"]
	if !ok {
		panic(errors.InvalidTokenError)
	}
	id, err := primitive.ObjectIDFromHex(val.(string))
	errors.CheckError(err)

	user, err := controllers.GetUser(&id)
	errors.CheckError(err)

	c.Set("user", user)
	c.Next()
}

// ValidateToken validates a jwt token
func ValidateToken(bearer string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InvalidTokenError
		}
		return []byte(env.JWT_SECRET), nil
	})
	errors.CheckError(err)

	if token.Valid {
		if val, ok := token.Claims.(jwt.MapClaims)["exp"]; ok {
			_time, err := time.Parse(time.RFC3339, val.(string))
			if err != nil {
				return nil, errors.InvalidTokenError
			}
			if time.Now().Before(_time) {
				return token.Claims.(jwt.MapClaims), nil
			} else {
				return nil, errors.ExpiredTokenError
			}
		}
	}

	return nil, errors.InvalidTokenError
}
