package eshop

import "context"

type CartService interface {
	New(ctx context.Context, items []CartItem) (*Cart, error)
	Cart(ctx context.Context, id string) (*Cart, error)
}
