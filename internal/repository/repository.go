package repository

import (
	"context"
	"time"
)

func newDBContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 3*time.Second)
}
