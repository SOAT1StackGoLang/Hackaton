package persistence

import "context"

type TimekeepingRepository interface {
	// CreateTimekeeping creates a new timekeeping record
	CreateTimekeeping(ctx context.Context, userID string) error
	// GetTimekeepingFromDate returns a timekeeping record for a given date
	ListTimekeepingFromDateAndUserID(ctx context.Context, userID string, date string) error
}
