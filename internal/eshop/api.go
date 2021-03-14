package eshop

import (
	"context"
	"net/http"

	"github.com/go-logrusutil/logrusutil/logctx"

	"github.com/labstack/echo/v4"

	apperr "github.com/govinda-attal/eshop/pkg/errors"
	"github.com/govinda-attal/eshop/pkg/eshop"
	"github.com/govinda-attal/eshop/pkg/eshop/api"
)

type Api struct {
	cartSrv eshop.CartService
}

func NewApi(ctx context.Context, cartSrv eshop.CartService) *Api {
	return &Api{
		cartSrv: cartSrv,
	}
}

// NewCart receives http.POST on /cart
func (a *Api) NewCart(ec echo.Context) (err error) {
	var (
		ctx = ec.Request().Context()
		log = logctx.From(ctx).WithField("operation", "NewCart")
	)
	defer func() {
		if err != nil {
			log.WithField("error", err).Error("error encountered when creating a new cart")
		}
	}()

	input := new(api.NewCartRq)
	if err := ec.Bind(input); err != nil {
		return apperr.BadRequest(err)
	}
	if len(input.LineItems) == 0 {
		return apperr.BadRequest().WithMessage("atleast one line item must be specified")
	}
	cart, err := a.cartSrv.New(ctx, api.MapToCartItems(input.LineItems))
	if err != nil {
		return err
	}
	return ec.JSON(http.StatusOK, api.MapFromCart(cart))
}

// GetCart receives http.GET on /cart/{id}
func (a *Api) GetCart(ec echo.Context, id string) (err error) {
	var (
		ctx = ec.Request().Context()
		log = logctx.From(ctx).WithField("operation", "GetCart")
	)
	defer func() {
		if err != nil {
			log.WithField("error", err).Error("error encountered when fetching an existing cart")
		}
	}()
	cart, err := a.cartSrv.Cart(ctx, id)
	if err != nil {
		return err
	}
	return ec.JSON(http.StatusOK, api.MapFromCart(cart))
}
