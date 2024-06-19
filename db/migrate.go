package db

func (s Storage) Migrate() error {
	createTableQuery := `
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
	        is_published BOOLEAN DEFAULT FALSE 
		);

		CREATE TABLE IF NOT EXISTS product_variants (
		    product_id INTEGER,
		    id INTEGER PRIMARY KEY,
		    name TEXT,
		    available INTEGER DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS product_prices (
		    product_id INTEGER,
		    id INTEGER PRIMARY KEY,
		    price INTEGER,
		    currency TEXT
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
		    email TEXT,
		    phone TEXT,
		    country TEXT,
		    address TEXT,
		    zip TEXT,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS cart (
		    id INTEGER PRIMARY KEY,
		    customer_id INTEGER,
		    discount_id INTEGER,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
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

		
		CREATE TABLE IF NOT EXISTS orders (
		    id INTEGER PRIMARY KEY,
		    customer_id INTEGER,
		    cart_id INTEGER,
		    discount_id INTEGER,
		    status TEXT,
		    payment_status TEXT,
		    currency TEXT,
		    shipping_status TEXT,
		    total INTEGER,
		    subtotal INTEGER,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    metadata TEXT,
		    payment_id TEXT,
		    FOREIGN KEY (customer_id) REFERENCES customers (id),
		    FOREIGN KEY (cart_id) REFERENCES cart (id),
		    FOREIGN KEY (discount_id) REFERENCES discounts (id)
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
		    deleted_at TIMESTAMP
		);
 	`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
