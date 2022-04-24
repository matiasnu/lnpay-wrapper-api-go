package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lnpay-wrapper-api-go/src/api/app/handlers"
	"github.com/lnpay-wrapper-api-go/src/api/config"
	"github.com/lnpay-wrapper-api-go/src/api/utils/logger"
)

var router *gin.Engine

func Start() {
	ConfigureRouter()

	if err := router.Run(config.ConfMap.APIRestServerPort); err != nil {
		logger.Errorf("Error starting router", err)
	}
}

func ConfigureRouter() {
	router = handlers.DefaultRouter()
	mapUrlsToControllers()
	logger.InitLog(logger.ConfigureLog{
		LoggingPath:  config.ConfMap.LoggingPath,
		LoggingFile:  config.ConfMap.LoggingFile,
		LoggingLevel: config.ConfMap.LoggingLevel,
	})
}

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.ServeHTTP(w, req)
}
