package routes

import (
	"net/http"

	"github.com/2HgO/gin-web-server/handlers"
	"github.com/gin-gonic/gin"
)

var categoryRoutes = []endpoint{
	{
		Path:   "",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetCategories,
		},
	},
	{
		Path:   "/category/:categoryID",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetCategory,
		},
	},
	{
		Path:   "/search",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.SearchCategories,
		},
	},
	{
		Path:   "",
		Method: http.MethodPut,
		Handlers: []gin.HandlerFunc{
			handlers.CreateCategory,
		},
	},
	{
		Path:   "/category/:categoryID",
		Method: http.MethodDelete,
		Handlers: []gin.HandlerFunc{
			handlers.DeleteCategory,
		},
	},
}
