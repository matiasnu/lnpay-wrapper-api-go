package app

import "github.com/lnpay-wrapper-api-go/src/api/controllers"

func mapUrlsToControllers() {
	router.GET("/ping", controllers.Ping)
}
