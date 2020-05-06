package main

import (
	"github.com/ichtrojan/thoth"
	"github.com/vergly/payment/routes"
	"log"
	"net/http"
)

func main() {
	api := routes.Api()

	logger, err := thoth.Init("log")

	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":9999", api); err != nil {
		logger.Log(err)
	}
}
