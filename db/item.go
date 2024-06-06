package db

// id INTEGER PRIMARY KEY,
// cart_id INTEGER,
// order_id INTEGER,
// variant_id INTEGER,
// quantity INTEGER,
// price INTEGER,
// currency TEXT,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// deleted_at TIMESTAMP,
// FOREIGN KEY (cart_id) REFERENCES cart (id),
// FOREIGN KEY (variant_id) REFERENCES product_variants (id),
// FOREIGN KEY (order_id) REFERENCES orders (id)

type LineItem struct {
	ID        int64  `db:"id" json:"id"`
	CartID    int64  `db:"cart_id" json:"cart_id"`
	OrderID   int64  `db:"order_id" json:"order_id"`
	VariantID int64  `db:"variant_id" json:"variant_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
	Price     int    `db:"price" json:"price"`
	Currency  string `db:"currency" json:"currency"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	DeletedAt string `db:"deleted_at" json:"deleted_at"`
}

func (s Storage) SaveLineItem(li LineItem) error {
	query := `
		INSERT INTO line_items (cart_id, order_id, variant_id, quantity, price, currency)
		VALUES (?, ?, ?, ?, ?, ?);
	`

	_, err := s.db.Exec(query, li.CartID, li.OrderID, li.VariantID, li.Quantity, li.Price, li.Currency)
	return err
}
