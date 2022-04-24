package handlers

/**
* @author mnunez
 */

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lnpay-wrapper-api-go/src/api/utils/apierrors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var production bool = os.Getenv("GO_ENVIRONMENT") == "production"

func DefaultRouter() *gin.Engine {
	return CustomRouter(RouterConfig{})
}

func CustomRouter(conf RouterConfig) *gin.Engine {
	router := gin.New()

	if !conf.DisableSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	if !production {
		router.Use(gin.Logger())
	}

	router.NoRoute(noRouteHandler)
	return router
}

type RouterConfig struct {
	DisableSwagger bool
}

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, apierrors.NewNotFoundApiError(fmt.Sprintf("Resource %s not found.", c.Request.URL.Path)))
}
