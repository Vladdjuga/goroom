package interfaces

import (
	"realTimeService/configuration"
	"realTimeService/handlers/wsrouter"
	"realTimeService/hubs"
)

type Container interface {
	GetHub() *hubs.MainHub
	GetRouter() *wsrouter.Router
	InitializeProviders(cfg *configuration.Config)
	Close() error
}
