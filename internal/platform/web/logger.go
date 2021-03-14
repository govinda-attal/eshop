package web

import (
	"context"

	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/google/uuid"
	"github.com/govinda-attal/eshop/internal/platform"
	"github.com/labstack/echo/v4"
)

func CtxMware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				rq   = c.Request()
				rqID = rq.Header.Get(platform.TraceID.String())
			)
			if rqID == "" {
				rqID = uuid.New().String()
			}
			log := logctx.Default.WithField(platform.TraceID.String(), rqID)
			ctx := logctx.New(rq.Context(), log)
			rq = rq.WithContext(context.WithValue(ctx, platform.TraceID, rqID))
			c.SetRequest(rq)
			return next(c)
		}
	}
}
