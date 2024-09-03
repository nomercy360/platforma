

-- Insert the test product into the products table
INSERT INTO products (id, handle, cover_image_url, image_urls, name, description, materials, is_published)
VALUES (28, 'test-product', '/images/products/test/1.svg',
        '/images/products/test/1.svg;/images/products/test/2.svg',
        'Test Product',
        'This is a test product to verify pricing in the store.',
        '100% cotton', false);

-- Insert a variant for the test product
INSERT INTO product_variants (product_id, id, name, available)
VALUES (28, 55, 'One Size', 10);

-- Insert the pricing for the test variant in USD and BYN
INSERT INTO variant_prices (variant_id, price, currency_code)
VALUES (55, 5, 'USD'),
       (55, 10, 'BYN');

-- Optionally, insert a discount or sale price for the test product
INSERT INTO sale_prices (variant_id, sale_price, currency_code, starts_at, ends_at)
VALUES (55, 1, 'USD', DATETIME('now'), DATETIME('now', '+60 days')),
       (55, 2, 'BYN', DATETIME('now'), DATETIME('now', '+60 days'));

-- Insert the translations for the test product
INSERT INTO product_translations (product_id, name, description, materials, language)
VALUES (28, 'Тестовый продукт',
        'Это тестовый продукт для проверки ценообразования в магазине.',
        '100% хлопок', 'ru');

