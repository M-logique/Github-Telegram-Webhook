package handler

import "github.com/gin-gonic/gin"



type Route struct {
	Method string
	Path   string
	Handler func(*gin.Context)
}

type Routes []Route


type GitHubWebhook struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	Repository struct {
		Name            string `json:"name"`
		FullName        string `json:"full_name"`
		HTMLURL         string `json:"html_url"`
		StargazersCount int    `json:"stargazers_count"`
		ForksCount      int    `json:"forks_count"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	Action string `json:"action"`
	After  string `json:"after"`
}

type RequestURI struct {
	BotToken   string `uri:"botToken" binding:"required"`
	ChatID     string `uri:"chatID" binding:"required"`
}

