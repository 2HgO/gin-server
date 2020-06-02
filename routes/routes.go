package routes

import (
	"github.com/2HgO/gin-web-server/handlers"
	"github.com/gin-gonic/gin"
)

type endpoint struct {
	Path     string
	Method   string
	Handlers []gin.HandlerFunc
}

func LoadRoutes(r *gin.Engine) {
	r.HandleMethodNotAllowed = true
	r.MaxMultipartMemory = 20 << 20
	r.Use(gin.Logger())
	r.Use(handlers.PanicHandler)
	r.Use(handlers.SetHeaders)
	r.Use(handlers.CORS)
	r.NoRoute(handlers.Handle404)
	r.NoMethod(handlers.Handle405)

	loadRoutes(r.Group("auth"), authRoutes...)

	r.Use(handlers.UserValidation)

	loadRoutes(r.Group("users"), userRoutes...)
	loadRoutes(r.Group("categories"), categoryRoutes...)
	loadRoutes(r.Group("watchlists"), watchlistRoutes...)
	loadRoutes(r.Group("movies"), movieRoutes...)
}

func loadRoutes(r *gin.RouterGroup, endpoints ...endpoint) {
	for _, endpoint := range endpoints {
		r.Handle(endpoint.Method, endpoint.Path, endpoint.Handlers...)
	}
}
