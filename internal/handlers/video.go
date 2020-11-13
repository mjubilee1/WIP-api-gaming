package handlers

import (
	"log"
	"net/http"
	"github.com/jackc/pgx/v4"
	"github.com/gin-gonic/gin"
	"context"
	"strconv"
	"api-gaming/internal/config"
	"time"
)

// Representation of stream data
type Stream struct {
	ID int64 `json:"id"`
	StreamID string `json:"streamID"`
	CreatedAt time.Time `json:"createdAt"`
}

// listen for the mux webhooks events
func create(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(*pgx.Conn)


	stream := config.CreateLiveStream()
	i, err := strconv.ParseInt(stream.Data.CreatedAt, 10, 64) // convert unit timestamp from string to int64 then create the time stamp with time.Unix
	
	if err != nil {
		panic(err)
	}

	tm := time.Unix(i, 0)
	_, err = conn.Exec(context.Background(), "INSERT INTO streaming (stream_id, created_at) VALUES ($1, $2)", stream.Data.Id, tm)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK,stream)
}

func get(c *gin.Context) {
	streamLookup := Stream{}
	db, _ := c.Get("db")
	conn := db.(*pgx.Conn)
	row := conn.QueryRow(context.Background(), "SELECT stream_id FROM streaming ORDER BY created_at DESC LIMIT 1")
	streamID := row.Scan(&streamLookup.StreamID)

	if streamID == pgx.ErrNoRows {
		log.Println("Error finding streaming id")
	}
	
	getStream := config.GetLiveStream(streamLookup.StreamID)
	c.JSON(http.StatusOK, getStream)
}

func addVideoHandler() {
	router.GET("/video/create-stream", create)
	router.GET("/video/current-stream", get)
}