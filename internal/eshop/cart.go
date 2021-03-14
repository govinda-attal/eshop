package eshop

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"

	"github.com/govinda-attal/eshop/internal/eshop/db"
	apperr "github.com/govinda-attal/eshop/pkg/errors"
	"github.com/govinda-attal/eshop/pkg/eshop"
)

type CartService struct {
	db *sqlx.DB
}

var _ eshop.CartService = &CartService{}

func NewCartService(ctx context.Context, db *sqlx.DB) *CartService {
	return &CartService{
		db: db,
	}
}

func (cs *CartService) New(ctx context.Context, items []eshop.CartItem) (*eshop.Cart, error) {
	state, err := cs.evaluateCart(ctx, items)
	if err != nil {
		return nil, err
	}
	id, err := db.NewCart(ctx, cs.db, items)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return &eshop.Cart{
		Id:    id,
		State: state,
	}, nil
}

func (cs *CartService) Cart(ctx context.Context, id string) (*eshop.Cart, error) {
	items, err := db.GetCartItems(ctx, cs.db, id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	state, err := cs.evaluateCart(ctx, items)
	if err != nil {
		return nil, err
	}
	return &eshop.Cart{
		Id:    id,
		State: state,
	}, nil
}

func (cs *CartService) evaluateCart(ctx context.Context, items []eshop.CartItem) (*eshop.CartState, error) {
	skus := make([]string, len(items))
	for i, ci := range items {
		skus[i] = ci.Sku
	}
	invItems, err := db.InventoryBySkus(ctx, cs.db, skus)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	itemProms, err := db.ItemPromotionsBySkus(ctx, cs.db, skus)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return EvaluateCart(ctx, items, invItems, itemProms)
}

func EvaluateCart(ctx context.Context, items []eshop.CartItem, invItems []eshop.InventoryItem, itemProms []eshop.ItemPromotions) (_ *eshop.CartState, err error) {

	skuInvItems := make(map[string]eshop.InventoryItem, len(invItems))
	for _, stock := range invItems {
		skuInvItems[stock.Sku] = stock
	}

	itemQuantities := make(map[string]int, len(items))
	for _, ci := range items {
		stock := skuInvItems[ci.Sku]
		if stock.Quantity < ci.Quantity {
			if err == nil {
				err = apperr.BadRequest().WithMessage("cart item(s) stock is less than requested quantity")
			}
			var e *apperr.Error
			if errors.As(err, &e) {
				err = e.WithAdditionalMessages(fmt.Sprintf("item (%s) current availability %d", stock.Name, stock.Quantity))
			}
		}
		itemQuantities[ci.Sku] = ci.Quantity
	}
	if err != nil {
		return nil, err
	}
	skuProms := make(map[string][]eshop.Promotion, len(itemProms))
	for _, item := range itemProms {
		skuProms[item.Sku] = item.Promotions
	}

	return evaluateCart(ctx, itemQuantities, skuInvItems, skuProms)
}

func evaluateCart(ctx context.Context, itemQuantities map[string]int, invItems map[string]eshop.InventoryItem, skuProms map[string][]eshop.Promotion) (*eshop.CartState, error) {

	var eii []*eshop.EvaluatedItem
	eiiBySku := make(map[string]*eshop.EvaluatedItem, len(itemQuantities))
	state := new(eshop.CartState)
	applyPromotions := func(ei *eshop.EvaluatedItem, proms []eshop.Promotion) *eshop.EvaluatedItem {
		// promotions are in order of priority
		// only one will be applied
		for _, prom := range proms {
			switch prom.Type {
			case eshop.PromotionPrice:
				quantity := itemQuantities[ei.Sku]
				if quantity < prom.Buy {
					continue
				}
				prom.Info = fmt.Sprintf("buy %d of (%s) for a price of %d", prom.Buy, ei.Name, prom.Units)
				q, r := divmod(quantity, prom.Buy)
				ei.Promotions = append(ei.Promotions, prom)
				ei.SalePrice += ei.ListPrice * float32(prom.Units) * float32(q)
				itemQuantities[ei.Sku] = r
				return ei
			case eshop.PromotionDiscount:
				quantity := itemQuantities[ei.Sku]
				if quantity < prom.Buy {
					continue
				}
				prom.Info = fmt.Sprintf("buy %d or more (%s) and have %.2f percent discount on all", prom.Buy, ei.Name, prom.Rate)
				atRate := (100.00 - prom.Rate) / 100.00
				ei.Promotions = append(ei.Promotions, prom)
				ei.SalePrice += (ei.ListPrice * float32(quantity)) * atRate
				itemQuantities[ei.Sku] = 0
				return ei
			case eshop.PromotionDiscountEvery:
				quantity := itemQuantities[ei.Sku]
				if quantity < prom.Buy {
					continue
				}
				prom.Info = fmt.Sprintf("buy %d or more (%s) and have %.2f percent discount on every %d", prom.Buy, ei.Name, prom.Rate, prom.Buy)
				q, r := divmod(quantity, prom.Buy)
				atRate := (100.00 - prom.Rate) / 100.00
				ei.Promotions = append(ei.Promotions, prom)
				ei.SalePrice += (ei.ListPrice * float32(prom.Buy) * float32(q)) * atRate
				itemQuantities[ei.Sku] = r
				return ei
			case eshop.PromotionFree:
				quantity := itemQuantities[ei.Sku]
				if quantity < prom.Buy {
					continue
				}
				related, ok := eiiBySku[prom.Item]
				if !ok {
					continue
				}
				relQuantity := itemQuantities[related.Sku]
				if relQuantity < prom.Units {
					continue
				}
				prom.Info = fmt.Sprintf("free %d (%s) with %d of (%s)", prom.Units, related.Name, prom.Buy, ei.Name)
				q, r := divmod(quantity, prom.Buy)
				relR := relQuantity - q
				ei.Promotions = append(ei.Promotions, prom)
				ei.SalePrice += ei.ListPrice * float32(q*prom.Buy)
				itemQuantities[ei.Sku] = r
				itemQuantities[related.Sku] = relR
				return ei
			}
		}
		return ei
	}

	for sku, quantity := range itemQuantities {
		invItem, ok := invItems[sku]
		if !ok {
			return nil, apperr.BadRequest().WithMessage("sku (%s) is not found", sku)
		}
		ei := &eshop.EvaluatedItem{
			Item: eshop.Item{
				Sku:       sku,
				Name:      invItem.Name,
				Quantity:  quantity,
				ListPrice: invItem.Price,
			},
		}
		state.BaseAmount += ei.ListPrice * float32(ei.Quantity)
		eii = append(eii, ei)
		eiiBySku[ei.Sku] = ei
	}
	eii = eshop.EvalItemsPriceReverse(eii)
	for i, ei := range eii {
		eii[i] = applyPromotions(ei, skuProms[ei.Sku])
	}

	for sku, quantity := range itemQuantities {
		ei := eiiBySku[sku]
		if quantity > 0 {
			ei.SalePrice += (ei.ListPrice * float32(quantity))
		}
		state.CartAmount += ei.SalePrice
		ei.Discount = float32(math.Round(float64(ei.ListPrice*float32(ei.Quantity)-ei.SalePrice)*100) / 100)

	}
	state.LineItems = eii
	return state, nil
}

func divmod(num, d int) (q, r int) {
	q = num / d
	r = num % d
	return q, r
}
