package main

import (
	"log"
	"api-gaming/internal/config"
	"context"
	routes "api-gaming/internal/handlers"
)

func init() {
	config.InitRedis()
	config.InitMux()
	//config.InitGoogle()
}

func main() {	
	connDB, err := config.InitDB()     
	if err != nil {
		log.Fatal(err)
		return
	}

	defer connDB.Close(context.Background())
	
	// Our server will live in the handlers package
	routes.Run()
}