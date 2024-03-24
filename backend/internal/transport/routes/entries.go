package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/SOAT1StackGoLang/Hackaton/internal/endpoints"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewTimekeepingRoutes(r *mux.Router, svc service.TimekeepingService, logger kitlog.Logger) *mux.Router {
	entries := endpoints.MakeTimekeepingEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodPost).Path("/api/clock-in").Handler(httptransport.NewServer(
		entries.InsertTimekeepingEndpoint,
		decodeInsertEntryRequest,
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
// @Param   user_id   header    string   false    "User ID" default(testing)
// @Param   entry_at  body    string   true    "Entry At"  SchemaExample({"entry_at": "2024-03-24T00:52:24Z"})
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
// @Router /api/clock-in [post]
func decodeInsertEntryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	claims, err := getJWTTokenJSON(r)
	if err != nil {
		log.Println(err)
	}
	username, ok := claims["username"].(string)
	if !ok {
		username = r.Header.Get("user_id")
		log.Println("Bad request: unable to find expected value. User ID from header:", username)
	}

	var request endpoints.InsertTimekeepingEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.UserID != "" {
		return nil, http.ErrBodyNotAllowed
	}

	request.UserID = username

	return request, nil
}
