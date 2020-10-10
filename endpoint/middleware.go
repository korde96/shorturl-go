package endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/endpoint"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("req", fmt.Sprint(request), "resp", fmt.Sprint(response), "transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}
