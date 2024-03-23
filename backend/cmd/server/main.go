package main

import (
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/persistence"
	"github.com/SOAT1StackGoLang/Hackaton/internal/transport"
	"github.com/SOAT1StackGoLang/Hackaton/internal/transport/routes"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"time"

	"log"
)

func main() {
	initializeApp()

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Panicf("failed initializing db: %s\n", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("failed getting db connection: %s\n", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Apply(&gorm.Config{
		ConnPool: sqlDB,
	}); err != nil {
		log.Panicf("failed applying db config: %s\n", err)
	}

	r := mux.NewRouter()

	r = routes.NewHelloRoutes(r, logger.InfoLogger)
	r = routes.NewSwaggerRoutes(r, logger.InfoLogger)

	entriesRepo := persistence.NewTimekeepingRepository(db, logger.InfoLogger)
	entriesSvc := service.NewTimekeepingService(entriesRepo)

	r = routes.NewEntriesRoutes(r, entriesSvc, logger.InfoLogger)

	transport.NewHTTPServer(":8080", muxToHttp(r))

}

func muxToHttp(r *mux.Router) http.Handler {
	return r
}
