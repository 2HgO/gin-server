package routes

import (
	"net/http"

	"github.com/2HgO/gin-web-server/handlers"
	"github.com/gin-gonic/gin"
)

var movieRoutes = []endpoint{
	{
		Path:   "/",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetMovies,
		},
	},
	{
		Path:   "/movie/:movieID",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetMovie,
		},
	},
	{
		Path:   "/category/:categoryID",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetMoviesByCategory,
		},
	},
	{
		Path:   "/search",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.SearchMovies,
		},
	},
	{
		Path:   "/",
		Method: http.MethodPut,
		Handlers: []gin.HandlerFunc{
			handlers.CreateMovie,
		},
	},
	{
		Path:   "/",
		Method: http.MethodDelete,
		Handlers: []gin.HandlerFunc{
			handlers.DeleteMovie,
		},
	},
}
