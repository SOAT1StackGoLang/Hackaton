package main

import (
	"flag"
	"github.com/SOAT1StackGoLang/Hackaton/pkg/helpers"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"os"
)

// initializeApp initializes the application by loading the configuration, connecting to the datastore,
// and subscribing to the Redis channel for receiving messages.
// It returns a pointer to the RedisStore and an error if any.

var (
	binding       string
	connString    string
	paymentURI    string
	productionURI string
)

func initializeApp() {
	flag.StringVar(&binding, "httpbind", ":8000", "address/port to bind listen socket")
	flag.Parse()
	//err := godotenv.Load()
	//if err != nil {
	//	logger.InfoLogger.Log("load err", err.Error())
	//}
	helpers.ReadPgxConnEnvs()
	paymentURI = os.Getenv("PAYMENT_URI")
	productionURI = os.Getenv("PRODUCTION_URI")
	connString = helpers.GetConnectionParams()

	logger.InitializeLogger()

	logger.Info("Connecting to datastore...")

	return
}
