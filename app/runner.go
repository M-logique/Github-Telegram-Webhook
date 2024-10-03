package handler

import "fmt"

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