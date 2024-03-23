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
// @Param   user_id   header    string   false    "User ID"
// @Param   entry_at  body    string   true    "Entry At"  SchemaExample({"entry_at": "2022-01-01T00:00:00Z"})
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
// @Router /entries [post]
func decodeCreateEntryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	claims, err := getJWTTokenJSON(r)
	if err != nil {
		log.Println(err)
	}

	username, ok := claims["username"].(string)
	if !ok {
		username = r.Header.Get("user_id")
		log.Println("Bad request: unable to find expected value. User ID from header:", username)
	}

	var requestBody struct {
		EntryAt string `json:"entry_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	request := endpoints.InsertEntryRequest{
		UserID:  username,
		EntryAt: requestBody.EntryAt,
	}

	return request, nil
}
