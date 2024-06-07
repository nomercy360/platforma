package db

import "time"

type LineItem struct {
	ID        int64      `db:"id" json:"id"`
	CartID    *int64     `db:"cart_id" json:"cart_id"`
	OrderID   *int64     `db:"order_id" json:"order_id"`
	VariantID int64      `db:"variant_id" json:"variant_id"`
	Quantity  int        `db:"quantity" json:"quantity"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`

	VariantName string `db:"variant_name" json:"variant_name"`
	Price       int    `db:"price" json:"price"`
	ProductName string `db:"product_name" json:"product_name"`
	ImageURL    string `db:"image_url" json:"image_url"`
}

func (s Storage) SaveLineItem(li LineItem) error {
	query := `
		INSERT INTO line_items (cart_id, order_id, variant_id, quantity)
		VALUES (?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, li.CartID, li.OrderID, li.VariantID, li.Quantity)

	if err != nil && IsDuplicateError(err) {
		return ErrAlreadyExists
	} else if err != nil {
		return err
	}

	return nil
}

func (s Storage) UpdateLineItemsOrderID(cartID, orderID int64) error {
	query := `
		UPDATE line_items
		SET order_id = ?
		WHERE cart_id = ?
	`

	_, err := s.db.Exec(query, orderID, cartID)
	return err
}

func (s Storage) UpdateLineItemQuantity(li int64, quantity int) error {
	query := `
		UPDATE line_items
		SET quantity = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(query, quantity, li)
	return err
}
