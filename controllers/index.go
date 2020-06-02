package controllers

import "github.com/2HgO/gin-web-server/configs/db"

type obj = map[string]interface{}

var (
	Watchlist = db.Watchlist
	Category  = db.Category
	Movie     = db.Movie
	User      = db.User
)
