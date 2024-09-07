package db

func (s Storage) Migrate() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS currencies (
		    code TEXT PRIMARY KEY,
		    name TEXT,
		    symbol TEXT,
		    UNIQUE(code)
		);

		INSERT INTO currencies (code, name, symbol) VALUES ('USD', 'US Dollar', '$') ON CONFLICT DO NOTHING;
		INSERT INTO currencies (code, name, symbol) VALUES ('BYN', 'Belarusian Ruble', 'BYN') ON CONFLICT DO NOTHING;
		
		CREATE TABLE IF NOT EXISTS products (
		    id INTEGER PRIMARY KEY,
		    handle TEXT,
		    cover_image_url TEXT,
		    image_urls TEXT,
		    name TEXT,
		    description TEXT,
		    materials TEXT,
	        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	        deleted_at TIMESTAMP,
	        is_published BOOLEAN DEFAULT FALSE,
	        UNIQUE(handle)
		);

		CREATE TABLE IF NOT EXISTS product_variants (
		    product_id INTEGER,
		    id INTEGER PRIMARY KEY,
		    name TEXT,
		    available INTEGER DEFAULT 0,
		    FOREIGN KEY (product_id) REFERENCES products (id)
		);

		CREATE TABLE IF NOT EXISTS variant_prices (
		    variant_id INTEGER,
		    id INTEGER PRIMARY KEY,
		    price INTEGER,
		    currency_code TEXT,
		    FOREIGN KEY (currency_code) REFERENCES currencies (code),
		    FOREIGN KEY (variant_id) REFERENCES product_variants (id),
		    UNIQUE(variant_id, currency_code)
		);
		
		CREATE TABLE IF NOT EXISTS product_translations (
		    product_id INTEGER,
		    id INTEGER PRIMARY KEY,
		    name TEXT,
		    description TEXT,
		    materials TEXT,
		    language TEXT
		);
		
		CREATE TABLE IF NOT EXISTS customers (
		    id INTEGER PRIMARY KEY,
		    name TEXT,
		    email TEXT NOT NULL,
		    phone TEXT,
		    country TEXT,
		    address TEXT,
		    zip TEXT,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    UNIQUE(email)
		);

		CREATE TABLE IF NOT EXISTS cart (
		    id INTEGER PRIMARY KEY,
		    customer_id INTEGER,
		    discount_id INTEGER,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    currency_code TEXT,
		    context TEXT,
		    FOREIGN KEY (customer_id) REFERENCES customers (id),
		    FOREIGN KEY (discount_id) REFERENCES discounts (id)
		);

		CREATE TABLE IF NOT EXISTS line_items (
		    id INTEGER PRIMARY KEY,
		    cart_id INTEGER,
		    order_id INTEGER,
		    variant_id INTEGER,
		    quantity INTEGER,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    FOREIGN KEY (cart_id) REFERENCES cart (id),
		    FOREIGN KEY (variant_id) REFERENCES product_variants (id),
		    FOREIGN KEY (order_id) REFERENCES orders (id),
		    UNIQUE(cart_id, variant_id)
		);

		CREATE TABLE IF NOT EXISTS regions (
			id INTEGER PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			currency_code TEXT,
			deleted_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS countries (
			id INTEGER PRIMARY KEY,
			display_name TEXT,
			iso_code TEXT,
			region_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			FOREIGN KEY (region_id) REFERENCES regions (id),
			UNIQUE(iso_code)
		);

		CREATE TABLE IF NOT EXISTS shipping_methods (
			id INTEGER PRIMARY KEY,
			name TEXT,
			price INTEGER,
			region_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			FOREIGN KEY (region_id) REFERENCES regions (id),
			UNIQUE(name)
		);
		
		CREATE TABLE IF NOT EXISTS orders (
		    id INTEGER PRIMARY KEY,
		    customer_id INTEGER,
		    cart_id INTEGER,
		    discount_id INTEGER,
		    status TEXT,
		    payment_status TEXT,
		    payment_provider TEXT,
			payment_id TEXT,
		    currency_code TEXT,
		    shipping_status TEXT,
		    shipping_method_id INTEGER,
		    total INTEGER,
		    subtotal INTEGER,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    metadata TEXT,
		    FOREIGN KEY (customer_id) REFERENCES customers (id),
		    FOREIGN KEY (cart_id) REFERENCES cart (id),
		    FOREIGN KEY (discount_id) REFERENCES discounts (id),
		    FOREIGN KEY (currency_code) REFERENCES currencies (code),
		    FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods (id)
		);

		CREATE TABLE IF NOT EXISTS discounts (
		    id INTEGER PRIMARY KEY,
		    value INTEGER,
		    code TEXT,
		    type TEXT DEFAULT 'percentage',
		    is_active BOOLEAN DEFAULT TRUE,
		    starts_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    ends_at TIMESTAMP,
		    usage_limit INTEGER DEFAULT 0,
		    usage_count INTEGER DEFAULT 0,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    UNIQUE(code)
		);

		CREATE TABLE IF NOT EXISTS sale_prices (
			variant_id INTEGER,
			sale_price INTEGER,
			starts_at TIMESTAMP,
			ends_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			currency_code TEXT,
			FOREIGN KEY (variant_id) REFERENCES product_variants(id),
			PRIMARY KEY (variant_id, currency_code, starts_at)
		);

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name TEXT,
			avatar_url TEXT,
			deleted_at TIMESTAMP,
			role TEXT NOT NULL DEFAULT 'user',
			UNIQUE(email)
		);

		CREATE TABLE IF NOT EXISTS product_categories (
			id INTEGER PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			UNIQUE(name)
		);

		CREATE TABLE IF NOT EXISTS product_category_products (
			product_id INTEGER,
			category_id INTEGER,
			UNIQUE(product_id, category_id),
			FOREIGN KEY (product_id) REFERENCES products (id),
			FOREIGN KEY (category_id) REFERENCES product_categories (id)
		);
 	`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
