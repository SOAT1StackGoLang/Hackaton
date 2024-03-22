package endpoints

import (
	"context"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service"
	"github.com/SOAT1StackGoLang/Hackaton/internal/service/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"time"
)

type (
	EntriesEndpoint struct {
		CreateEntryEndpoint endpoint.Endpoint
	}
)

func MakeEntriesEndpoint(svc service.EntriesService) EntriesEndpoint {
	return EntriesEndpoint{
		CreateEntryEndpoint: makeCreateEntryEndpoint(svc),
	}

}

func makeCreateEntryEndpoint(svc service.EntriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			entryAt time.Time
		)

		req := request.(InsertEntryRequest)

		uID, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, err
		}

		if req.EntryAt == "" {
			entryAt = time.Now()
		} else {
			entryAt, err = time.Parse(time.RFC3339, req.EntryAt)
			if err != nil {
				return nil, err
			}
		}

		in := &models.Entry{
			UserID:    uID,
			CreatedAt: entryAt,
		}

		eOut, err := svc.CreateEntry(ctx, in)
		if err != nil {
			return nil, err
		}

		out := InsertEntryResponse{
			ID:      eOut.ID.String(),
			UserID:  eOut.UserID.String(),
			EntryAt: eOut.CreatedAt.Format(time.RFC3339),
		}

		return out, nil
	}
}
