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
		GetReportByReferenceDate     endpoint.Endpoint
		GetReportByRange             endpoint.Endpoint
		GetReportByRangeAndUserIDCSV endpoint.Endpoint
	}
)

func MakeReportsEndpoint(svc service.ReportService) ReportsEndpoint {
	return ReportsEndpoint{
		GetReportByReferenceDate: makeGetReportByReferenceDateEndpoint(svc),
		GetReportByRange:         makeGetReportByRangeEndpoint(svc),
	}
}

func makeGetReportByReferenceDateAndUserIDCSVEndpoint(rS service.ReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(TimekeepingCSVReportRequest)

		location, err := time.LoadLocation("America/Sao_Paulo")
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed loading location", err.Error()))
			return
		}

		start, err := time.ParseInLocation("2006-01-02", req.Start, location)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed parsing time", err.Error()))
			return
		}
		end, err := time.ParseInLocation("2006-01-02", req.End, location)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed parsing time", err.Error()))
			return
		}

		byteOut, err := rS.GetReportCSVByRangeAndUserID(ctx, req.UserID, start, end)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed getting report", err.Error()))
			return nil, err
		}

		return byteOut, nil
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

func makeGetReportByRangeEndpoint(svc service.ReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(TimekeepingReportRequest)
		location, err := time.LoadLocation("America/Sao_Paulo")
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed loading location", err.Error()))
			return
		}

		start, err := time.ParseInLocation("2006-01-02", req.Start, location)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed parsing time", err.Error()))
			return
		}
		end, err := time.ParseInLocation("2006-01-02", req.End, location)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %s", "failed parsing time", err.Error()))
			return
		}

		servOut, err := svc.GetReportByRangeAndUserID(ctx, req.UserID, start, end)

		return timeKeepingReportResponseFromModels(servOut), nil
	}
}
