package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/VladimirB/gometrics/internal/shared/logger"
	"go.uber.org/zap"
)

func DoRetry(ctx context.Context, fn func() error) error {
	steps := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	var lastErr error

	for i := range steps {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry aborted by context: %w", err)
		}

		if err := fn(); err != nil {
			logger.Log.Warn("retry failed", zap.Error(err))

			lastErr = err

			select {
			case <-time.After(steps[i]):
			case <-ctx.Done():
				return ctx.Err()
			}
			continue
		}

		return nil
	}

	return fmt.Errorf("all retries failed, last error: %w", lastErr)
}
