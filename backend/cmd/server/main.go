package main

import (
	"github.com/SOAT1StackGoLang/Hackaton/internal/transport"
	"github.com/SOAT1StackGoLang/Hackaton/internal/transport/routes"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"

	"log"
)

func main() {
	initializeApp()

	_, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Panicf("failed initializing db: %s\n", err)
	}

	r := mux.NewRouter()
	r = routes.NewHelloRoutes(r, logger.InfoLogger)

	transport.NewHTTPServer(":8080", muxToHttp(r))

}

func muxToHttp(r *mux.Router) http.Handler {
	return r
}
