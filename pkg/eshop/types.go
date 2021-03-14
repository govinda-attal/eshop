package eshop

import "sort"

type CartItem struct {
	Quantity int    `json:"quantity"`
	Sku      string `json:"sku"`
}

type EvaluatedItem struct {
	Item
	Promotions []Promotion `json:"promotions,omitempty"`
	Discount   float32     `json:"discount"`
	SalePrice  float32     `json:"salePrice"`
}

type Item struct {
	Sku       string  `json:"sku" db:"sku"`
	Name      string  `json:"name" db:"price"`
	Quantity  int     `json:"quantity" db:"quantity"`
	ListPrice float32 `json:"listPrice" db:"listPrice"`
}

type InventoryItem struct {
	Sku      string  `json:"sku" db:"sku"`
	Name     string  `json:"name" db:"name"`
	Price    float32 `json:"price" db:"price"`
	Quantity int     `json:"quantity" db:"quantity"`
}

type Promotion struct {
	Buy   int           `json:"buy"`
	Type  PromotionType `json:"type"`
	Info  string        `json:"info"`
	Item  string        `json:"item,omitempty"`
	Rate  float32       `json:"rate,omitempty"`
	Units int           `json:"Units,omitempty"`
}

type ItemPromotions struct {
	Sku        string      `json:"sku"`
	Promotions []Promotion `json:"promotions"`
}

type Cart struct {
	Id    string     `json:"id"`
	State *CartState `json:"state"`
}

type CartState struct {
	BaseAmount float32          `json:"baseAmount"`
	CartAmount float32          `json:"cartAmount"`
	LineItems  []*EvaluatedItem `json:"lineItems"`
}

type (
	PromotionType string
)

const (
	PromotionFree          PromotionType = "FREE"
	PromotionDiscount      PromotionType = "DISCOUNT"
	PromotionDiscountEvery PromotionType = "DISCOUNT-EVERY"
	PromotionPrice         PromotionType = "PRICE"
)

func EvalItemsPriceReverse(ii []*EvaluatedItem) []*EvaluatedItem {
	sort.Slice(ii, func(i, j int) bool {
		return ii[i].ListPrice > ii[j].ListPrice
	})
	return ii
}
