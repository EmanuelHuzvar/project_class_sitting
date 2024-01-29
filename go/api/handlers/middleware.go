package handlers

import "github.com/gin-gonic/gin"

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, token, language ,teamname")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("language", "csharp")
		c.Header("language", "java")
		c.Header("language", "python")
		c.Header("teamname", "*")
		//c.Writer.Header().Set("teamname", "My-Value")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
