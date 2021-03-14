package web

import (
	"fmt"

	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	apperr "github.com/govinda-attal/eshop/pkg/errors"
	"github.com/govinda-attal/eshop/pkg/errors/codes"
	"github.com/govinda-attal/eshop/pkg/eshop/api"
)

func ErrorHandler(err error, c echo.Context) {
	rq := c.Request()
	log := logctx.From(rq.Context())

	log.WithFields(logrus.Fields{
		"error":  err,
		"method": rq.Method,
		"url":    rq.URL.String(),
	}).Errorf("error message: %s", err.Error())

	if e, ok := err.(*echo.HTTPError); ok {
		err = apperr.New(codes.FromHTTPStatus(e.Code), fmt.Sprintf("%s", e.Message))
	}

	if e, ok := err.(*apperr.Error); ok {
		ae := api.MapAppError(e)
		_ = c.JSON(e.Code.HTTPStatusCode(), ae)
		return
	}
}
