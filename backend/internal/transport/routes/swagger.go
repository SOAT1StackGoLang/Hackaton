package routes

import (
	"net/http"

	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"

	_ "github.com/SOAT1StackGoLang/Hackaton/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewSwaggerRoutes(r *mux.Router, logger kitlog.Logger) *mux.Router {

	r.Methods(http.MethodGet).PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL pointing to the API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
	// redirect / to /swagger/index.html
	r.Methods(http.MethodGet).Path("/").Handler(http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))

	return r
}
