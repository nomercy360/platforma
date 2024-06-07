package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type ArrayString []string

func (as *ArrayString) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		*as = make([]string, 0)
	case []byte:
		*as = strings.Split(string(src), ";")
	case string:
		*as = strings.Split(src, ";")
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return nil
}

func (as ArrayString) Value() (driver.Value, error) {
	if len(as) == 0 {
		return nil, nil
	}

	return strings.Join(as, ";"), nil
}

type Product struct {
	ID          int64     `json:"id"`
	Handle      string    `json:"handle"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants"`
	Image       string    `json:"image"`
	Images      []string  `json:"images"`
	Currency    string    `json:"currency"`
	Price       int       `json:"price"`
}

type Variant struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Available int    `json:"available"`
}

func listProductQuery(locale string) string {
	var query string

	switch locale {
	case "ru", "by":
		query = `
		SELECT p.id, p.handle, pt_ru.name, pt_ru.description, pp.currency, pp.price,
		       p.cover_image_url, p.image_urls,
		       json_group_array(distinct json_object('id', pv.id, 'name', pv.name, 'quantity', pv.quantity)) AS variants
		FROM products p
		LEFT JOIN product_translations pt_ru ON p.id = pt_ru.product_id AND pt_ru.language = 'ru'
		LEFT JOIN product_prices pp ON p.id = pp.product_id AND pp.currency = 'BYN'
		LEFT JOIN product_variants pv ON p.id = pv.product_id
	`
	default:
		query = `
		SELECT p.id, p.handle, p.name, p.description, pp.currency, pp.price,
		       p.cover_image_url, p.image_urls,
		       json_group_array(distinct json_object('id', pv.id, 'name', pv.name, 'quantity', pv.quantity)) AS variants
		FROM products p
		LEFT JOIN product_prices pp ON p.id = pp.product_id AND pp.currency = 'USD'
		LEFT JOIN product_variants pv ON p.id = pv.product_id
	`
	}

	return query
}

func (s Storage) ListProducts(locale string) ([]Product, error) {
	query := listProductQuery(locale)

	query += fmt.Sprintf(" WHERE p.is_published = TRUE GROUP BY p.id")

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product
	for rows.Next() {
		var id int64
		var price int
		var handle, name, description, imageUrl, currency, variantsJSON string
		var imageUrls ArrayString

		if err := rows.Scan(
			&id,
			&handle,
			&name,
			&description,
			&currency,
			&price,
			&imageUrl,
			&imageUrls,
			&variantsJSON,
		); err != nil {
			return nil, err
		}

		var variants []Variant
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			return nil, err
		}

		product := Product{
			ID:          id,
			Handle:      handle,
			Name:        name,
			Currency:    currency,
			Price:       price,
			Description: description,
			Image:       imageUrl,
			Variants:    variants,
			Images:      imageUrls,
		}

		products = append(products, product)
	}

	return products, nil
}

type GetProductQuery struct {
	Handle string
	ID     int64
	Locale string
}

func (s Storage) GetProduct(q GetProductQuery) (*Product, error) {
	query := listProductQuery(q.Locale)

	var args []interface{}

	if q.Handle != "" {
		query += " WHERE p.handle = ? AND p.is_published = TRUE"
		args = append(args, q.Handle)
	} else {
		query += " WHERE p.id = ? AND p.is_published = TRUE"
		args = append(args, q.ID)
	}

	query = fmt.Sprintf("%s GROUP BY p.id", query)

	var product Product
	var variantsJSON string
	var imageUrls ArrayString

	if err := s.db.QueryRow(query, args...).Scan(
		&product.ID,
		&product.Handle,
		&product.Name,
		&product.Description,
		&product.Currency,
		&product.Price,
		&product.Image,
		&imageUrls,
		&variantsJSON,
	); err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var variants []Variant
	if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
		return nil, err
	}

	product.Variants = variants
	product.Images = imageUrls

	return &product, nil
}
