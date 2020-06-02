package env

import "os"

var (
	APP_ENV    = os.Getenv("APP_ENV")
	APP_PORT   = os.Getenv("APP_PORT")
	DB_NAME    = os.Getenv("DB_NAME")
	DB_URL     = os.Getenv("DB_URL")
	JWT_SECRET = os.Getenv("JWT_SECRET")
)
