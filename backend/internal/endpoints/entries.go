package endpoints

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type (
	EntriesEndpoint struct {
		InsertTimekeepingEndpoint endpoint.Endpoint
	}
)

func MakeTimekeepingEndpoint(svc service.TimekeepingService) EntriesEndpoint {
	return EntriesEndpoint{
		InsertTimekeepingEndpoint: makeCreateEntryEndpoint(svc),
	}

}

func makeCreateEntryEndpoint(svc service.TimekeepingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			entryAt time.Time
		)

		req := request.(InsertTimekeepingEntryRequest)

		if req.UserID == "" {
			return nil, errors.New("missing user id")
		}

		if req.EntryAt == "" {
			entryAt = time.Now()
		} else {
			entryAt, err = time.Parse(time.RFC3339, req.EntryAt)
			if err != nil {
				return nil, err
			}
		}

		resp, err := svc.InsertEntry(ctx, req.UserID, entryAt)
		if err != nil {
			return nil, err
		}
		
		return timeKeepingResponseFromModels(resp), nil
	}
}
