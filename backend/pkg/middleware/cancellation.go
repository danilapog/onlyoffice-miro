package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/labstack/echo/v4"
)

type CancellationMiddleware struct {
	logger service.Logger
}

func NewCancellationMiddleware(logger service.Logger) *CancellationMiddleware {
	return &CancellationMiddleware{
		logger: logger,
	}
}

func (m *CancellationMiddleware) HandleRequestCancellation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		c.SetRequest(req.WithContext(ctx))

		done := make(chan error)

		go func() {
			done <- next(c)
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				m.logger.Info(ctx, "request canceled by client",
					service.Fields{
						"path":   req.URL.Path,
						"method": req.Method,
						"error":  ctx.Err().Error(),
					})
				return echo.NewHTTPError(http.StatusRequestTimeout, "Request canceled")
			} else if ctx.Err() == context.DeadlineExceeded {
				m.logger.Info(ctx, "request timeout",
					service.Fields{
						"path":   req.URL.Path,
						"method": req.Method,
						"error":  ctx.Err().Error(),
					})
				return echo.NewHTTPError(http.StatusRequestTimeout, "Request timeout")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		case err := <-done:
			return err
		}
	}
}
