package main

import (
	"log"
	"realTimeService/configuration"
	"realTimeService/controllers"
	"realTimeService/handlers"
	"realTimeService/interfaces"
	"realTimeService/middlewares"
	"realTimeService/providers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := configuration.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	// Create Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("views/templates/*")

	// Serve static files
	router.Static("/static", "./views/static")

	// Initialize the DI container
	var container interfaces.Container = providers.NewDependencyInjectionContainer()
	container.InitializeProviders(cfg)
	defer func(container interfaces.Container) {
		err := container.Close()
		if err != nil {
			log.Fatalf("Failed to close dependency injection container: %v", err)
			return
		}
	}(container)

	// Initialize controllers
	homeController := controllers.NewHomeController()
	chatController := controllers.NewChatController()
	wsHandler := handlers.NewWsHandler(container)

	// Use middleware
	router.Use(gin.Recovery())

	// HTTP routes
	router.GET("/", homeController.Index)
	router.GET("/chat", chatController.Index)

	// WebSocket endpoint with simplified auth (no JWT required)
	router.GET("/ws", middlewares.SimpleAuthMiddleware(cfg), wsHandler.Handle)

	log.Printf("🚀 Starting anonymous chat server on %s", cfg.HttpPort)
	log.Printf("📍 Home page: http://localhost%s", cfg.HttpPort)
	log.Printf("💬 Chat page: http://localhost%s/chat", cfg.HttpPort)

	err = router.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
