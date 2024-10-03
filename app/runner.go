package handler

import "fmt"

// Run initializes the application by loading configuration values from the environment,
// setting up the router with the specified routes, and starting the HTTP server.
//
// The function sets up the routes using the provided listenEndpoint and starts the server
// to listen for incoming connections.
func Run() {
	config, err := loadConfigFromEnv()
	
	if err != nil {
		panic(err)
	}
	
	listenPort := config["port"]
	listenAddr := config["addr"]
	listenPath := config["path"]

	var listenEndpoint string

	if listenPath, ok := listenPath.(string); ok {
		listenEndpoint = listenPath
	} else {
		panic("path should be an string")
	}


	routes := loadRoutes(listenEndpoint)
	router := NewRouter(routes)

	router.Run(fmt.Sprintf("%s:%v", listenAddr, listenPort))
}