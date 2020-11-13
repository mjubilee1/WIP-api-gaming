package config

import (
	"log"
	"github.com/jackc/pgx/v4"
	"context"
	"github.com/gin-gonic/gin"
	"api-gaming/internal/util"
)

var router = gin.Default()
// POSTGRESURL - Gets the Postgres env url
var POSTGRESURL = util.ViperEnvVariable("POSTGRES_URL")

// DBCONN - Sets the database connection
var DBCONN, err = pgx.Connect(context.Background(), POSTGRESURL)

// InitDB - Used as main connection when application first starts
func InitDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), POSTGRESURL)
	if err != nil {
		log.Println("Error connecting to Sleepless Gamer's Database")
		log.Println(err.Error())
	}

	log.Println("Database connected successfully.", POSTGRESURL)
	return conn, err
}

// DBMiddleware - Passes the database context as middleware
func DBMiddleware(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}