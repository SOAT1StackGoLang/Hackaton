package endpoints

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	"github.com/go-kit/kit/endpoint"
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

func makeGetReportByReferenceDateEndpoint(_ service.ReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return TimekeepingReportResponse{
			UserID:        "mock",
			ReferenceDate: "2024-01-01",
			Open:          "false",
			WorkedHours:   "8 hora(s) e 25 minuto(s)",
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
								CreatedAt: "2024-03-23T09:00:00-03:00",
							},
							EndingEntry: &Entry{
								ID:        "batata-entry2",
								CreatedAt: "2024-03-23T13:00:00-03:00",
							},
						},
						{
							WorkedTime: "4 hora(s) e 25 minuto(s)",
							StartingEntry: Entry{
								ID:        "batata-entry3",
								CreatedAt: "2024-03-23T14:00:00-03:00",
							},
							EndingEntry: &Entry{
								ID:        "batata-entry-3",
								CreatedAt: "2024-03-23T18:25:00-03:00",
							},
						},
					},
				},
			},
		}, nil
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
