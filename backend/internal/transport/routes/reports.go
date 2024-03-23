package routes

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/endpoints"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

func NewReportRoutes(r *mux.Router, svc service.ReportService, logger kitlog.Logger) *mux.Router {
	reports := endpoints.MakeReportsEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

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

	return r
}

func decodeGetReportByReference(_ context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeGetReportRequestByRange(_ context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
}
