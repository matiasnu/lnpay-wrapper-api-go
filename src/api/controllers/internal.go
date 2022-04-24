package controllers

import "github.com/gin-gonic/gin"

// Ping is the handler of test app
// @Summary Ping
// @Description test if the router works correctly
// @Tags ping
// @Produce  json
// @Success 200
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.String(200, "pong")
}
