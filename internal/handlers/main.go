package handlers

import (
	"time"
	cors "github.com/rs/cors/wrapper/gin"
	"api-gaming/internal/config"
	"github.com/gin-gonic/gin"
	"context"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
)

var router = gin.Default()
var authenticator auth.Authenticator
var cache store.Cache

// Run will start the server
func Run() {
	getHandlers()
	router.Run(":9990")
}

/* Setup GoGuardian - A simple clean, and idomatic way 
* to create a powerful modern API and web authentication.
* Sole purpoose is to authenticate requests, which it does
* through an extensible set of authentication methods known 
* as strategies.
*/
func setupGoGuardian() {   
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*5)
	tokenStrategy := bearer.New(verifyToken, cache)
	authenticator.EnableStrategy(bearer.CachedStrategyKey,    tokenStrategy)
}

/*func handlerMiddleware() gin.HandlerFunc {

}*/

// Get handlers will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getHandlers() {
	corsConfig := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge: int(12 * time.Hour),
	})
	
	// Set up application middlewares
	router.Use(corsConfig)
	router.Use(config.DBMiddleware(config.DBCONN))
	setupGoGuardian()
	addUserHandlers()
	addVideoHandler()
}