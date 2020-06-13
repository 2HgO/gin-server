package routes

import (
	"net/http"

	"github.com/2HgO/gin-web-server/handlers"
	"github.com/gin-gonic/gin"
)

var watchlistRoutes = []endpoint{
	{
		Path:   "",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetWatchlist,
		},
	},
	{
		Path:   "/user/:userID",
		Method: http.MethodGet,
		Handlers: []gin.HandlerFunc{
			handlers.GetUserWatchlist,
		},
	},
	{
		Path:   "/movie/:movieID",
		Method: http.MethodPut,
		Handlers: []gin.HandlerFunc{
			handlers.AddMovieToUserWatchlist,
		},
	},
	{
		Path:   "/movie/:movieID",
		Method: http.MethodDelete,
		Handlers: []gin.HandlerFunc{
			handlers.RemoveMovieFromUserWatchlist,
		},
	},
}
