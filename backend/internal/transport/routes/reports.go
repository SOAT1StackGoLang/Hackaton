package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/SOAT1StackGoLang/Hackaton/internal/endpoints"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewReportRoutes(r *mux.Router, svc service.ReportService, logger kitlog.Logger) *mux.Router {
	reports := endpoints.MakeReportsEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).
		Path("/api/reports").
		Queries("start", "{start:[0-9]{4}-[0-9]{2}-[0-9]{2}}", "end", "{end:[0-9]{4}-[0-9]{2}-[0-9]{2}}").
		Handler(
			httptransport.
				NewServer(
					reports.GetReportByRange,
					decodeGetReportRequestByRange,
					encodeResponse,
					options...,
				),
		)

	r.Methods(http.MethodGet).
		Path("/api/reports/daily").
		Queries("reference", "{reference:[0-9]{4}-[0-9]{2}-[0-9]{2}}").
		Handler(
			httptransport.
				NewServer(
					reports.GetReportByReferenceDate,
					decodeGetReportByReference,
					encodeResponse,
					options...,
				),
		)

	r.Methods(http.MethodGet).
		Path("/api/reports/csv")

	return r
}

// decodeGetReportByReference godoc
// @Summary Decode JWT token and get report by reference
// @Tags Reports
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user_id     header   string     false  "User ID" default(testing)
// @Param   reference   query    string     true  "Reference Date" Format("2006-01-02") default(2024-03-24)
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Bad Routing"
// @Router /api/reports/daily [get]
func decodeGetReportByReference(_ context.Context, r *http.Request) (interface{}, error) {
	claims, err := getJWTTokenJSON(r)
	if err != nil {
		log.Println(err)
	}
	username, ok := claims["username"].(string)
	if !ok {
		username = r.Header.Get("user_id")
		log.Println("Bad request: unable to find expected value. User ID from header:", username)
	}

	request := endpoints.TimekeepingReportByReferenceRequest{}

	request.UserID = username
	request.ReferenceDate, ok = mux.Vars(r)["reference"]
	if !ok {
		return nil, ErrBadRouting
	}

	return request, nil

}

// decodeGetReportRequestByRange godoc
// @Summary Decode JWT token and get report by range
// @Tags Reports
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user_id     header   string     true  "User ID" default(testing)
// @Param   start       query    string     true  "Start Date" Format("2006-01-02") default(2024-03-01)
// @Param   end         query    string     true  "End Date" Format("2006-01-02") default(2024-03-31)
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Bad Routing"
// @Router /api/reports [get]
func decodeGetReportRequestByRange(_ context.Context, r *http.Request) (interface{}, error) {
	claims, err := getJWTTokenJSON(r)
	if err != nil {
		log.Println(err)
	}
	username, ok := claims["username"].(string)
	if !ok {
		username = r.Header.Get("user_id")
		log.Println("Bad request: unable to find expected value. User ID from header:", username)
	}

	req := endpoints.TimekeepingReportRequest{
		UserID: username,
	}

	req.Start, ok = mux.Vars(r)["start"]
	if !ok {
		return nil, ErrBadRouting
	}
	req.End, ok = mux.Vars(r)["end"]
	if !ok {
		return nil, ErrBadRouting
	}

	return req, nil
}
