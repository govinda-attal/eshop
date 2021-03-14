package eshop_test

import (
	"context"
	"testing"

	. "github.com/govinda-attal/eshop/internal/eshop"
	"github.com/govinda-attal/eshop/pkg/eshop"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateCart(t *testing.T) {

	ctx := context.TODO()

	const (
		macBookPro   = "43N23P"
		googleHome   = "120P90"
		alexaSpeaker = "A304SD"
		raspberryPiB = "234234"
	)

	var (
		inventory = []eshop.InventoryItem{
			{Sku: macBookPro, Name: "MacBook Pro", Price: 5399.99, Quantity: 5},
			{Sku: googleHome, Name: "Google Home", Price: 49.99, Quantity: 10},
			{Sku: alexaSpeaker, Name: "Alexa Speaker", Price: 109.50, Quantity: 10},
			{Sku: raspberryPiB, Name: "Raspberry Pi B", Price: 30.00, Quantity: 2},
		}

		promotions = []eshop.ItemPromotions{
			{
				Sku: macBookPro,
				Promotions: []eshop.Promotion{
					{
						Info: "free raspberry pi with macbook pro",
						Buy:  1, Type: eshop.PromotionFree, Item: raspberryPiB, Units: 1,
					},
				},
			},
			{
				Sku: googleHome,
				Promotions: []eshop.Promotion{
					{
						Info: "buy three google homes at price of two",
						Buy:  3, Type: eshop.PromotionPrice, Units: 2,
					},
				},
			},
			{
				Sku: alexaSpeaker,
				Promotions: []eshop.Promotion{
					{
						Info: "buy three or more alexa speakers and have 10 percent discount on all",
						Buy:  3, Type: eshop.PromotionDiscount, Rate: 10.00,
					},
				},
			},
		}
	)

	scenarios := []struct {
		description string
		cartItems   []eshop.CartItem
		cartAmount  float32
		evalErr     bool
	}{
		{
			description: "Scanned Items: MacBook Pro, Raspberry Pi B; Total: $5,399.99",
			cartItems: []eshop.CartItem{
				{Sku: macBookPro, Quantity: 1},
				{Sku: raspberryPiB, Quantity: 1},
			},
			cartAmount: 5399.99,
			evalErr:    false,
		},
		{
			description: "Scanned Items: Three Google Homes; Total: $99.98",
			cartItems: []eshop.CartItem{
				{Sku: googleHome, Quantity: 3},
			},
			cartAmount: 99.98,
			evalErr:    false,
		},
		{
			description: "Scanned Items: Three Alexa Speakers; Total: $295.65",
			cartItems: []eshop.CartItem{
				{Sku: alexaSpeaker, Quantity: 3},
			},
			cartAmount: 295.65,
			evalErr:    false,
		},
	}

	for i, scn := range scenarios {
		t.Log("-----------------------------\n", "scenario: ", scn.description)
		state, err := EvaluateCart(ctx, scn.cartItems, inventory, promotions)
		isErr := err != nil
		assert.Equal(t, scn.evalErr, isErr, "evaluation error mismatch for scenario #%d", i)
		var amt float32
		if state != nil {
			amt = state.CartAmount
		}
		assert.Equal(t, scn.cartAmount, amt, "cart checkout amount mismatch for scenario #%d", i)
	}
}
