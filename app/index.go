package handler

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfigFromEnv()
	
	if err != nil {
		panic(err)
	}
	
	listenPath := config["path"]

	var listenEndpoint string

	if listenPath, ok := listenPath.(string); ok {
		listenEndpoint = listenPath
	} else {
		panic("path should be an string")
	}


	routes := loadRoutes(listenEndpoint)
	router := NewRouter(routes)
	router.ServeHTTP(w, r)

	
}