package main

import (
	"net/http"

	"github.com/SCKelemen/Cassius/log"
	"github.com/SCKelemen/Cassius/config"
	"github.com/SCKelemen/Cassius/data"
	"github.com/SCKelemen/Cassius/mail"
	"github.com/SCKelemen/Cassius/web"
)

func main() {
	config, err := config.LoadConfigFromFile("cassius.config")
	if err != nil {
		panic(err)
	}

	logger, err := log.NewLogger() 
	if err != nil {
		panic(err)
	}

	pool, err := data.NewPool(config, logger)
	if err != nil {
		panic(err)
	}

	mailer, err := mail.NewMailer(config, logger)
	if err != nil {
		panic(err)
	}

	apiHandler := web.NewAPIHandler(config, pool, mailer, logger)
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", apiHandler))

	err = http.ListenAndServe(config.ListenLocation, nil)
	if err != nil {
		panic(err)
	}
}
