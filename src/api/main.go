package main

import (
	"github.com/lnpay-wrapper-api-go/src/api/app"
	"github.com/lnpay-wrapper-api-go/src/api/config"
	_ "github.com/lnpay-wrapper-api-go/src/api/docs"
)

// @title LNPay Wrapper API
// @version 1.0
// @description This is a wrapper golang api for lnpay.
// @termsOfService http://swagger.io/terms/

// @contact.name Matias Nuñez
// @contact.url http://www.swagger.io/support
// @contact.email matiasne45@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	config.Load("config", "parameters", "yml")
	app.Start()
}
