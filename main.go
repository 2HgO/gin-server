package main

import (
	"log"

	"github.com/2HgO/gin-web-server/configs/env"
	"github.com/2HgO/gin-web-server/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	routes.LoadRoutes(router)

	log.Fatalln(router.Run(":" + env.APP_PORT))
}
