package main

import (
	"log"

	"github.com/2HgO/gin-web-server/configs/env"
)

func init() {
	switch "" {
	case env.JWT_SECRET:
		log.Fatalln("JWT secret not set")
	case env.DB_URL:
		log.Fatalln("DB URL not set")
	case env.DB_NAME:
		log.Fatalln("DB Name not set")
	default:
		log.Println("All dependency variables available")
	}
}
