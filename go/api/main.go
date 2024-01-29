package main

import (
	"Visma/config"
	"Visma/helpers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	error := helpers.RemoveLogFile()
	if error != nil {
		return
	}

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	config.SetupRouter(router)
	err := router.Run("" + helpers.IpAddress() + ":9090")
	if err != nil {
		return
	}

}
