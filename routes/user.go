package routes

import (
	"net/http"

	"github.com/2HgO/gin-web-server/handlers"
	"github.com/gin-gonic/gin"
)

var userRoutes = []endpoint{
	{
		Path: "",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetUsers,
		},
	},
	{
		Path: "/user/:userID",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetUser,
		},
	},
	{
		Path: "",
		Method: http.MethodPatch,
		Handlers: []gin.HandlerFunc{
			handlers.UpdateUser,
		},
	},
	{
		Path: "",
		Method: http.MethodDelete,
		Handlers: []gin.HandlerFunc{
			handlers.DeleteUser,
		},
	},
	{
		Path: "/search",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.SearchUsers,
		},
	},
}