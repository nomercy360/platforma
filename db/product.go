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

type Prices struct {
	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`
	Price          int    `json:"price"`
	SalePrice      *int   `json:"sale_price"`
	IsOnSale       bool   `json:"is_on_sale"`
}

type Product struct {
	ID          int64            `json:"id"`
	Handle      string           `json:"handle"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Variants    []ProductVariant `json:"variants"`
	Image       string           `json:"image"`
	Images      []string         `json:"images"`
	Materials   string           `json:"materials"`
}

type ProductVariant struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Available int      `json:"available"`
	Prices    []Prices `json:"prices"`
}

func listProductQuery() string {

	return `
		SELECT p.id,
			   p.handle,
			   COALESCE(pt.name, p.name)               AS name,
			   COALESCE(pt.description, p.description) AS name,
			   COALESCE(pt.materials, p.materials)     AS name,
			   p.cover_image_url,
			   p.image_urls,
			   json_group_array(
					   json_object(
							   'id', pv.id,
							   'name', pv.name,
							   'available', pv.available,
							   'prices', (SELECT json_group_array(
														 json_object(
																 'currency_code', vp.currency_code,
																 'currency_symbol', c.symbol,
																 'price', vp.price,
																 'sale_price', sp.sale_price
														 )
												 )
										  FROM variant_prices vp
												   JOIN currencies c ON vp.currency_code = c.code
												   LEFT JOIN sale_prices sp
															 ON vp.variant_id = sp.variant_id AND c.code = sp.currency_code
										  WHERE vp.variant_id = pv.id
										  GROUP BY vp.variant_id)
					   )
			   )                                       AS variants
		FROM products p
				 LEFT JOIN product_variants pv ON p.id = pv.product_id
				 LEFT JOIN product_translations pt ON p.id = pt.product_id AND pt.language = ?
`
}

func (s Storage) ListProducts(locale string) ([]Product, error) {
	query := listProductQuery()

	query += fmt.Sprintf(" WHERE p.is_published = TRUE GROUP BY p.id")

	rows, err := s.db.Query(query, locale)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var id int64
		var handle, name, description, materials, imageUrl, variantsJSON string
		var imageUrls ArrayString

		if err := rows.Scan(
			&id,
			&handle,
			&name,
			&description,
			&materials,
			&imageUrl,
			&imageUrls,
			&variantsJSON,
		); err != nil {
			return nil, err
		}

		var variants []ProductVariant
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			return nil, err
		}

		product := Product{
			ID:          id,
			Handle:      handle,
			Name:        name,
			Description: description,
			Materials:   materials,
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
	query := listProductQuery()

	args := []interface{}{q.Locale}

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
		&product.Materials,
		&product.Image,
		&imageUrls,
		&variantsJSON,
	); err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var variants []ProductVariant
	if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
		return nil, err
	}

	product.Variants = variants
	product.Images = imageUrls

	return &product, nil
}
