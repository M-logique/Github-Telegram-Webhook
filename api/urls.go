package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)




func loadRoutes(listenPath string) Routes {
	
	urls := Routes{
		{
			"POST",
			listenPath+"/:botToken/:chatID",
			func(ctx *gin.Context) {
				var uri RequestURI
				blacklistedActions := ctx.Query("blacklisted_actions")
				blacklistedEvents  := ctx.Query("blacklisted_events")

				if err := ctx.ShouldBindUri(&uri); err != nil {
					ctx.JSON(400, gin.H{"msg": err.Error()})
					return
				}
			
				tgStatus := SendToTelegram(
					uri.ChatID, 
					uri.BotToken, 
					ctx.Request, 
					blacklistedActions,
					blacklistedEvents,
				)
				ctx.JSON(http.StatusOK, gin.H{"message": "Send to Telegram: " + tgStatus})
			},
		},
	}
	
	return urls
}