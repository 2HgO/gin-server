package routes

import (
	"net/http"

	"github.com/2HgO/gin-web-server/handlers"
	"github.com/gin-gonic/gin"
)

var authRoutes = []endpoint{
	{
		Path: "/login",
		Method: http.MethodPost,
		Handlers: []gin.HandlerFunc{
			handlers.Login,
		},
	},
	{
		Path: "/forgot-password",
		Method: http.MethodPost,
		Handlers: []gin.HandlerFunc{
			handlers.ForgotPassword,
		},
	},
	{
		Path: "/reset-password",
		Method: http.MethodPost,
		Handlers: []gin.HandlerFunc{
			handlers.ResetPassword,
		},
	},
	{
		Path: "/sign-up",
		Method: http.MethodPut,
		Handlers: []gin.HandlerFunc{
			handlers.CreateUser,
		},
	},
}