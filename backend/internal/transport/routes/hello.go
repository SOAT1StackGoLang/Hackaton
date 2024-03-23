package routes

import (
	"context"

	"github.com/SOAT1StackGoLang/Hackaton/internal/endpoints"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"net/http"

	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewHelloRoutes(r *mux.Router, logger kitlog.Logger) *mux.Router {
	hello := endpoints.MakeHelloEndpoints()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/hello").Handler(httptransport.NewServer(
		hello.GetHelloEndpoint,
		func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
			return request, nil
		},
		encodeResponse,
		options...,
	))

	return r

}
