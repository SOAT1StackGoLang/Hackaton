package main

import (
	"flag"
	"github.com/joho/godotenv"

	"github.com/SOAT1StackGoLang/Hackaton/pkg/helpers"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
)

// initializeApp initializes the application by loading the configuration, connecting to the datastore,
// and subscribing to the Redis channel for receiving messages.
// It returns a pointer to the RedisStore and an error if any.

var (
	connString string
)

func initializeApp() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		logger.InfoLogger.Log("load err", err.Error())
	}

	helpers.ReadPgxConnEnvs()
	connString = helpers.GetConnectionParams()

	logger.InitializeLogger()

	return
}
