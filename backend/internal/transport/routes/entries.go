package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SOAT1StackGoLang/Hackaton/internal/endpoints"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
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

// CreateEntry godoc
//
// @Summary Create a new entry
// @Tags Entries
// @Security	ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user_id   header    string   true    "User ID"
// @Param   entry_at  body    string   true    "Entry At"  SchemaExample({"entry_at": "2022-01-01T00:00:00Z"})
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
// @Router /entries [post]
func decodeCreateEntryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// TODO extract userID
	var request endpoints.InsertEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}
