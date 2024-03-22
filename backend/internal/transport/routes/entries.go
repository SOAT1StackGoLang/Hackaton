package routes

import (
	"context"
	"encoding/json"
	"github.com/SOAT1StackGoLang/Hackaton/internal/endpoints"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

func NewEntriesRoutes(r *mux.Router, svc service.EntriesService, logger kitlog.Logger) *mux.Router {
	entries := endpoints.MakeEntriesEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodPost).Path("/entries").Handler(httptransport.NewServer(
		entries.CreateEntryEndpoint,
		decodeCreateEntryRequest,
		encodeResponse,
		options...,
	))

	return r

}

func decodeCreateEntryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// TODO extract userID
	var request endpoints.InsertEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}
