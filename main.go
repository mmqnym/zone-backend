package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"

	"zoneBackend/internal/api/futures"
	"zoneBackend/internal/api/general"
	"zoneBackend/internal/api/spot"
	"zoneBackend/internal/middleware/auth"
	"zoneBackend/internal/middleware/rateLimit"
)

var useRateLimit bool = true

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go waitShutdown(c)

	server := setupServer()

	log.Println("Server is running on port 8888")
	server.Run(":8888")
}

func setupServer() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	server := gin.Default()
	server.UseH2C = true
	server.Use(auth.Resolve)

	if useRateLimit {
		server.Use(rateLimit.Resolve)
	}

	registRoute(server)

	return server
}

func registRoute(server *gin.Engine) {
	server.GET("/exchangeInfo", general.GetExchangeInfo)

	spotGroup := server.Group("/spot")

	spotGroup.GET("/balance", spot.GetBalance)
	spotGroup.GET("/transfer/records", spot.GetTransactionRecords)

	futureGroup := server.Group("/futures")

	futureGroup.GET("/balance", futures.GetBalance)
}

func waitShutdown(sig chan os.Signal) {
	<-sig

	log.Println("Server is shutting down...")
	os.Exit(0)
}
