package providers

import (
	"log"
	"realTimeService/configuration"
	"realTimeService/handlers/wsrouter"
	"realTimeService/handlers/wsrouter/handlers"
	"realTimeService/hubs"
	"realTimeService/models"
)

// DependencyInjectionContainer DI Container
type DependencyInjectionContainer struct {
	Hub *hubs.MainHub

	// Router for WebSocket handling
	Router *wsrouter.Router
}

// NewDependencyInjectionContainer Create a new DI container
func NewDependencyInjectionContainer() *DependencyInjectionContainer {
	return &DependencyInjectionContainer{}
}

// InitializeProviders Initialize the singleton variables
func (d *DependencyInjectionContainer) InitializeProviders(cfg *configuration.Config) {
	d.Hub = hubs.NewMainHub()
	d.Router = wsrouter.NewRouter()

	// Register WebSocket message handlers
	d.Router.RegisterHandler(models.FindMatch, handlers.NewFindMatchHandler(d))
	d.Router.RegisterHandler(models.SendMessage, handlers.NewSendHandler(d))
	d.Router.RegisterHandler(models.NextStranger, handlers.NewNextStrangerHandler(d))
	d.Router.RegisterHandler(models.StopChat, handlers.NewStopChatHandler(d))

	log.Println("DependencyInjectionContainer initialized with Hub and MatchingService")
}

func (d *DependencyInjectionContainer) GetHub() *hubs.MainHub {
	return d.Hub
}

func (d *DependencyInjectionContainer) GetRouter() *wsrouter.Router {
	return d.Router
}

func (d *DependencyInjectionContainer) Close() error {
	log.Println("Closing DependencyInjectionContainer")
	return nil
}
