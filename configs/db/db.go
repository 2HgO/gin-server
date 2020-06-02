package db

import (
	"log"

	"github.com/2HgO/gin-web-server/configs/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Watchlist *Collection
	Category  *Collection
	Movie     *Collection
	User      *Collection
)

// initialize connections
func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI(env.DB_URL))

	if err != nil {
		log.Fatalln(err)
	}

	if err := client.Connect(nil); err != nil {
		log.Fatalln(err)
	}

	db := client.Database(env.DB_NAME)

	Watchlist = newCollection(db, "watchlists")
	Category = newCollection(db, "categories")
	Movie = newCollection(db, "movies")
	User = newCollection(db, "users")
}
