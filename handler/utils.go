package handler

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// loadConfigFromEnv loads configuration values from the .env file
//
//
// If there is an error loading the .env file or converting the port
// to an integer, it returns an error.
//
// Returns:
// - A map[string]interface{} containing:
//     - "port": The port number as an integer.
//     - "addr": The server's bind address as a string.
// - An error if the .env file cannot be loaded or the port is invalid.
//
func loadConfigFromEnv() (map[string]interface{}, error) {

    godotenv.Load()

    pathStr := os.Getenv("LISTEN_PATH")
    if pathStr == "" {
        pathStr = "/" // Defaults to / if LISTEN_PATH is not set 
    }

    portStr := os.Getenv("LISTEN_PORT")
    if portStr == "" {
        portStr = "8080"  // Defaults to 8080 if LISTEN_PORT is not set
    }

    addr := os.Getenv("LISTEN_ADDRESS")
    if addr == "" {
        addr = "127.0.0.1"  // Defaults to 127.0.0.1 if LISTEN_ADDRESS is not set
    }

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid port value: %v", err)
    }

    config := map[string]interface{}{
        "port": port,
        "addr": addr,
        "path": pathStr,
    }

    return config, nil
}


// NewRouter initializes a new Gin router and registers the provided routes.
//
// Params:
// - routes: A slice of Route structs, where each struct contains the method (string),
//           path (string), and handler function (gin.HandlerFunc).
//
// Example usage:
//
//  routes := Routes{
//         {Method: "GET", Path: "/users", Handler: getUsersHandler},
//         {Method: "POST", Path: "/users", Handler: createUserHandler},
//  }
//  router := NewRouter(routes)
//
// Returns:
//  - A *gin.Engine instance with the specified routes registered.
//
// The `Method` field in each Route can be any valid HTTP method (e.g., "GET", "POST", "PUT", "PATCH", "DELETE", etc.),
// and the function will dynamically bind the appropriate handler to that method.
func NewRouter(routes Routes) *gin.Engine {
    router := gin.Default()
    for _, route := range routes {
        // Handle dynamically maps the route's method to the handler.
        router.Handle(route.Method, route.Path, route.Handler)
    }

	router.Static("/gen", "./templates/generate_webhook")

    return router
}


func FormatGitHubEvent(ctx *gin.Context) string {
	var webhook GitHubWebhook

	if err := ctx.ShouldBindJSON(&webhook); err != nil {
		return "Error binding JSON payload"
	}

	var message string
	switch webhook.RefType {
	case "branch":
		if webhook.Action == "created" {
			message = "Branch created: " + webhook.Ref
		} else if webhook.Action == "deleted" {
			message = "Branch deleted: " + webhook.Ref
		}
	case "tag":
		message = "Tag: " + webhook.Ref + " has been created"
	default:
		message = "Action: " + webhook.Action + " on repository: " + webhook.Repository.Name + " by user: " + webhook.Sender.Login
	}

	return message
}

func isBlacklisted(e, blacklisted string) bool {
	blacklist := strings.Split(blacklisted, ",")
	for _, s := range blacklist {
		if strings.TrimSpace(s) == e {
			return true
		}
	}
	return false
}