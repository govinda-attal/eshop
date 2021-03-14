package api

import (
	apperr "github.com/govinda-attal/eshop/pkg/errors"
	"github.com/govinda-attal/eshop/pkg/eshop"
)

func MapAppError(err *apperr.Error) *Error {
	e := &Error{
		Code:    int(err.Code),
		Message: err.Message,
	}
	var dtls []ErrorDetail
	for _, d := range err.Details {

		dtls = append(dtls, ErrorDetail{
			Message: d.Message,
			Code:    strPtr(d.Code),
		})
	}
	if len(dtls) > 0 {
		e.Details = &dtls
	}
	return e
}

func MapToCartItems(items []CartItem) []eshop.CartItem {
	out := make([]eshop.CartItem, len(items))
	for i, item := range items {
		out[i] = eshop.CartItem{
			Quantity: item.Quantity,
			Sku:      item.Sku,
		}
	}
	return out
}

func MapFromCart(c *eshop.Cart) *Cart {
	if c == nil {
		return nil
	}

	out := &Cart{
		Id:    c.Id,
		State: MapFromCartState(c.State),
	}
	return out
}

func MapFromCartState(state *eshop.CartState) (out CartState) {
	if state == nil {
		return
	}
	out.BaseAmount = state.BaseAmount
	out.CartAmount = state.CartAmount
	out.LineItems = make([]EvaluatedItem, len(state.LineItems))

	for i, item := range state.LineItems {
		out.LineItems[i] = EvaluatedItem{
			Sku:       item.Sku,
			Quantity:  item.Quantity,
			ListPrice: item.ListPrice,
			Discount:  item.Discount,
			SalePrice: item.SalePrice,
		}
		proms := make([]Promotion, len(item.Promotions))
		for j, prom := range item.Promotions {
			proms[j] = Promotion{
				Buy:   prom.Buy,
				Info:  prom.Info,
				Type:  string(prom.Type),
				Item:  strPtr(prom.Item),
				Rate:  floatPtr(prom.Rate),
				Units: intPtr(prom.Units),
			}
		}
		out.LineItems[i].Promotions = proms
	}
	return out
}

func strPtr(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
func floatPtr(f float32) *float32 {
	if f == 0 {
		return nil
	}
	return &f
}
func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
