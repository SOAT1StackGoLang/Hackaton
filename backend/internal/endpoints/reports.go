package endpoints

import (
	"context"
	"fmt"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	logger "github.com/SOAT1StackGoLang/Hackaton/pkg/middleware"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	ReportsEndpoint struct {
		GetReportByReferenceDate endpoint.Endpoint
		GetReportByRange         endpoint.Endpoint
	}
)

func MakeReportsEndpoint(svc service.ReportService) ReportsEndpoint {
	return ReportsEndpoint{
		GetReportByReferenceDate: makeGetReportByReferenceDateEndpoint(svc),
		GetReportByRange:         makeGetReportByRangeEndpoint(svc),
	}
}

func makeGetReportByReferenceDateEndpoint(rS service.ReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(TimekeepingReportByReferenceRequest)

		location, err := time.LoadLocation("America/Sao_Paulo")
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed loading location", err.Error()))
			return
		}

		parsedTime, err := time.ParseInLocation("2006-01-02", req.ReferenceDate, location)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed parsing time", err.Error()))
			return
		}

		servOut, err := rS.GetReportByReferenceDateAndUserID(ctx, req.UserID, parsedTime)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed getting report", err.Error()))
			return nil, err
		}

		return timeKeepingResponseFromModels(servOut), nil
	}
}

func makeGetReportByRangeEndpoint(_ service.ReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return TimekeepingReportResponse{
			UserID:      "mock",
			Start:       "2024-01-01",
			End:         "2024-01-02",
			Open:        "false",
			WorkedHours: "16 hora(s) e 50 minuto(s)",
			Report: []TimekeepingResponse{
				{
					ID:            "batata1",
					UserID:        "mock",
					ReferenceDate: "2024-01-01",
					UpdatedAt:     "2024-01-01",
					WorkedTime:    "8 hora(s) e 25 minuto(s)",
					Open:          false,
					Details: []Detail{
						{
							WorkedTime: "4 hora(s)",
							StartingEntry: Entry{
								ID:        "batata-entry1",
								CreatedAt: "2024-01-01T09:00:00-03:00",
							},
							EndingEntry: &Entry{
								ID:        "batata-entry2",
								CreatedAt: "2024-01-01T13:00:00-03:00",
							},
						},
						{
							WorkedTime: "4 hora(s) e 25 minuto(s)",
							StartingEntry: Entry{
								ID:        "batata-entry3",
								CreatedAt: "2024-01-01T14:00:00-03:00",
							},
							EndingEntry: &Entry{
								ID:        "batata-entry-3",
								CreatedAt: "2024-01-01T18:25:00-03:00",
							},
						},
					},
				},
				{
					ID:            "batata2",
					UserID:        "mock",
					ReferenceDate: "2024-01-02",
					UpdatedAt:     "2024-01-02",
					WorkedTime:    "8 hora(s) e 25 minuto(s)",
					Open:          false,
					Details: []Detail{
						{
							WorkedTime: "4 hora(s)",
							StartingEntry: Entry{
								ID:        "batata-entry4",
								CreatedAt: "2024-01-02T09:00:00-03:00",
							},
							EndingEntry: &Entry{
								ID:        "batata-entry5",
								CreatedAt: "2024-01-02T13:00:00-03:00",
							},
						},
						{
							WorkedTime: "4 hora(s) e 25 minuto(s)",
							StartingEntry: Entry{
								ID:        "batata-entry6",
								CreatedAt: "2024-01-02T14:00:00-03:00",
							},
							EndingEntry: &Entry{
								ID:        "batata-entry-7",
								CreatedAt: "2024-01-02T18:25:00-03:00",
							},
						},
					},
				},
			},
		}, nil
	}
}
