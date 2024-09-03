package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Cart struct {
	ID             int64           `json:"id" db:"id"`
	Items          []LineItem      `json:"items" db:"items"`
	CustomerID     *int64          `json:"customer_id" db:"customer_id"`
	CurrencyCode   string          `json:"currency_code" db:"currency_code"`
	CurrencySymbol string          `json:"currency_symbol" db:"currency_symbol"`
	Total          int             `json:"total" db:"total"`
	Count          int             `json:"count" db:"count"`
	Subtotal       int             `json:"subtotal" db:"subtotal"`
	Discount       *Discount       `json:"discount" db:"discount"`
	DiscountID     *int64          `json:"discount_id" db:"discount_id"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at" db:"deleted_at"`
	Context        CustomerContext `json:"context" db:"context"`
	DiscountAmount int             `json:"discount_amount" db:"-"`
	Customer       *Customer       `json:"customer" db:"-"`
}

type CustomerContext struct {
	UserAgent string  `json:"user_agent" db:"user_agent"`
	IP        string  `json:"ip" db:"ip"`
	Country   *string `json:"country" db:"country"`
}

func currencyFromLocale(locale string) string {
	switch locale {
	case "en":
		return "USD"
	case "ru", "by":
		return "BYN"
	default:
		return "USD"
	}
}

func (cc CustomerContext) Value() (driver.Value, error) {
	if cc == (CustomerContext{}) {
		return nil, nil
	}

	b, err := json.Marshal(cc)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (cc *CustomerContext) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		*cc = CustomerContext{}
	case []byte:
		if err := json.Unmarshal(src, cc); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(src), cc); err != nil {
			return err
		}
	default:
		return errors.New("unsupported type")
	}

	return nil
}

func (s Storage) GetCartByID(id int64, locale string) (*Cart, error) {
	var cart Cart
	q := `
		SELECT 
			c.id,
			c.customer_id,
			c.created_at,
			c.updated_at,
			c.deleted_at,
			c.context,
			c.currency_code,
			COALESCE(cr.symbol, '$') AS currency_symbol,
			c.discount_id
		FROM
			cart c
		LEFT JOIN currencies cr ON c.currency_code = cr.code
		WHERE
			c.id = ?;`

	row := s.db.QueryRow(q, id)

	err := row.Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.DeletedAt,
		&cart.Context,
		&cart.CurrencyCode,
		&cart.CurrencySymbol,
		&cart.DiscountID,
	)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	query := LineItemQuery{
		Locale:   locale,
		CartID:   id,
		Currency: cart.CurrencyCode,
	}

	items, err := s.GetLineItems(query)

	if err != nil {
		return nil, err
	}

	cart.Items = items
	for _, item := range items {
		salePrice := item.Price
		if item.SalePrice != nil {
			salePrice = *item.SalePrice
		}

		cart.Subtotal += item.Price * item.Quantity
		cart.Count += item.Quantity
		cart.Total += salePrice * item.Quantity
	}

	if cart.DiscountID != nil {
		discount, err := s.GetDiscount(DiscountQuery{ID: *cart.DiscountID})
		if err != nil {
			return nil, err
		}

		if discount.Value > 0 {
			var discountAmount int
			switch discount.Type {
			case "percentage":
				discountAmount = cart.Total * discount.Value / 100
			case "fixed":
				discountAmount = discount.Value
			}

			cart.Total -= discountAmount
			cart.DiscountAmount = discountAmount
		}

		cart.Discount = discount
	}

	// delivery
	if cart.CurrencyCode == "USD" {
		cart.Total += 10
	} else {
		cart.Total += 25
	}

	// only for testing purposes
	if len(cart.Items) == 1 && cart.Items[0].ProductName == "Test Product" {
		cart.Total = 1
		cart.Subtotal = 1
	}

	if cart.CustomerID != nil {
		customer, err := s.GetCustomerByID(*cart.CustomerID)
		if err != nil {
			return nil, err
		}

		cart.Customer = customer
	}

	return &cart, nil
}

func lineItemQuery() string {
	return `
			SELECT li.id,
				   li.cart_id,
				   li.order_id,
				   li.variant_id,
				   li.quantity,
				   li.created_at,
				   li.updated_at,
				   li.deleted_at,
				   pv.name AS variant_name,
				   COALESCE(pt.name, p.name) AS product_name,
				   p.cover_image_url AS image_url,
				   vp.price,
				   sp.sale_price
			FROM line_items li
			JOIN product_variants pv on li.variant_id = pv.id
			JOIN products p on pv.product_id = p.id
			LEFT JOIN sale_prices sp on li.variant_id = sp.variant_id AND sp.currency_code = ?
			JOIN variant_prices vp on li.variant_id = vp.variant_id AND vp.currency_code = ?
			LEFT JOIN product_translations pt on p.id = pt.product_id AND pt.language = ?
		`
}

type LineItemQuery struct {
	Locale   string
	CartID   int64
	OrderID  int64
	Currency string
}

func (s Storage) GetLineItems(query LineItemQuery) ([]LineItem, error) {
	q := lineItemQuery()

	var currency string
	if query.Currency != "" {
		currency = query.Currency
	} else {
		currency = currencyFromLocale(query.Locale)
	}

	args := []interface{}{currency, currency, query.Locale}

	if query.CartID > 0 {
		q = fmt.Sprintf("%s WHERE li.cart_id = %d", q, query.CartID)
		args = append(args, query.CartID)
	} else if query.OrderID > 0 {
		q = fmt.Sprintf("%s WHERE li.order_id = %d", q, query.OrderID)
		args = append(args, query.OrderID)
	}

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []LineItem
	for rows.Next() {
		var item LineItem
		if err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.OrderID,
			&item.VariantID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
			&item.VariantName,
			&item.ProductName,
			&item.ImageURL,
			&item.Price,
			&item.SalePrice,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (s Storage) CreateCart(cart Cart, locale string) (*Cart, error) {
	res, err := s.db.Exec("INSERT INTO cart (customer_id, context, currency_code) VALUES (?, ?, ?)",
		cart.CustomerID, cart.Context, cart.CurrencyCode)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	if len(cart.Items) > 0 {
		for _, item := range cart.Items {
			li := LineItem{
				CartID:    &id,
				VariantID: item.VariantID,
				Quantity:  item.Quantity,
			}

			if err := s.SaveLineItem(li); err != nil {
				return nil, err
			}
		}
	}

	return s.GetCartByID(id, locale)
}

func (s Storage) UpdateCartCurrency(cartID int64, currency string) error {
	_, err := s.db.Exec("UPDATE cart SET currency_code = ? WHERE id = ?", currency, cartID)

	return err
}

func (s Storage) DeleteCart(id int64) error {
	_, err := s.db.Exec("UPDATE cart SET deleted_at = ? WHERE id = ?", time.Now(), id)

	return err
}

func (s Storage) UpdateCartDiscount(cartID int64, discountID int64) error {
	_, err := s.db.Exec("UPDATE cart SET discount_id = ? WHERE id = ?", discountID, cartID)

	return err
}

func (s Storage) DropCartDiscount(cartID int64) error {
	_, err := s.db.Exec("UPDATE cart SET discount_id = NULL WHERE id = ?", cartID)

	return err
}

func (s Storage) UpdateCartCustomer(cartID int64, customerID int64) error {
	_, err := s.db.Exec("UPDATE cart SET customer_id = ? WHERE id = ?", customerID, cartID)

	return err
}
