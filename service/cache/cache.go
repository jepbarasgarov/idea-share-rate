package cache

import (
	"context"
	"time"
)

type Service interface {
	TokenSetWithExpiry(ctx context.Context, userID string, token string, expiry time.Duration) (err error)
	GetUserIDByToken(ctx context.Context, token string) (userID *string, err error)
	DeletePreviousRefreshToken(ctx context.Context, token string) (err error)

	WorkerRestrictWithExpiry(ctx context.Context, fullname string, expiry time.Duration) (err error)
	HasRestrctionForWorker(ctx context.Context, fullname string) (hasBlock bool, err error)
}
