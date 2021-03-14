package db

import (
	"context"
	"strings"

	"github.com/govinda-attal/eshop/pkg/eshop"
	"github.com/jmoiron/sqlx"
)

func InventoryBySkus(ctx context.Context, db *sqlx.DB, skus []string) ([]eshop.InventoryItem, error) {
	var (
		query = `SELECT sku, name, price, quantity 
						FROM eshop.inventory WHERE sku in ('` + strings.Join(skus, `','`) + `')`
		items []eshop.InventoryItem
	)
	if err := db.SelectContext(ctx, &items, query); err != nil {
		return nil, err
	}

	return items, nil
}

func NewCart(ctx context.Context, db *sqlx.DB, cii []eshop.CartItem) (id string, err error) {
	var (
		query = `INSERT INTO eshop.carts (cart) VALUES ($1) RETURNING id`
	)
	err = db.QueryRowContext(ctx, query, CartItems(cii)).Scan(&id)
	return id, err
}

func ItemPromotionsBySkus(ctx context.Context, db *sqlx.DB, skus []string) ([]eshop.ItemPromotions, error) {

	type skuPromRow struct {
		Sku        string     `db:"sku"`
		Promotions Promotions `db:"promotions"`
	}
	var (
		query = `SELECT sku, promotions 
						FROM eshop.promotions WHERE sku in ('` + strings.Join(skus, `','`) + `')`
		rows []skuPromRow
	)
	if err := db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}
	proms := make([]eshop.ItemPromotions, len(rows))

	for i, row := range rows {
		proms[i] = eshop.ItemPromotions{
			Sku:        row.Sku,
			Promotions: []eshop.Promotion(row.Promotions),
		}
	}
	return proms, nil
}
