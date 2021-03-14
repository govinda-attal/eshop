package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/govinda-attal/eshop/pkg/eshop"
)

type (
	CartItems      []eshop.CartItem
	EvaluatedItems []eshop.EvaluatedItem
	Promotions     []eshop.Promotion
)

func (ci CartItems) Value() (driver.Value, error) {
	return json.Marshal(ci)
}

func (ci *CartItems) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ci)
}

func (ei EvaluatedItems) Value() (driver.Value, error) {
	return json.Marshal(ei)
}

func (ei *EvaluatedItems) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ei)
}

func (pp Promotions) Value() (driver.Value, error) {
	return json.Marshal(pp)
}

func (pp *Promotions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &pp)
}
