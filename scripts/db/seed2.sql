INSERT INTO products (id, handle, cover_image_url, image_urls, name, description, materials, is_published)
VALUES (15, 'silk-drawstring-pouch', '/images/products/15/1.png',
        '/images/products/15/1.png;/images/products/15/2.png;/images/products/15/3.png;/images/products/15/4.png',
        'Silk drawstring pouch',
        'A petite silk pouch accented with a delicate rose, perfect for adding a touch of romance to your summer outfits.',
        '100% silk', true),
       (16, 'loose-shirt-with-ribbons', '/images/products/16/1.png',
        '/images/products/16/1.png;/images/products/16/2.png;/images/products/16/3.png;/images/products/16/4.png',
        'Loose shirt with ribbons',
        'A whimsical shirt made of rustling silk in a delightful pistachio ice cream hue, fastened with thin ribbons instead of traditional buttons.',
        '100% cotton', true),
       (17, 'blush-pink-dress', '/images/products/17/1.png',
        '/images/products/17/1.png;/images/products/17/2.png;/images/products/17/3.png;/images/products/17/4.png;/images/products/17/5.png',
        'Blush pink dress',
        'A dress made of dense blush pink summer cotton, beautifully complemented by a cherry red ribbon trim.',
        '100% silk', true),
       (18, 'top-with-diagonal-cutout', '/images/products/18/1.png',
        '/images/products/18/1.png',
        'Top with diagonal cutout',
        'A light and airy summer cotton top inspired by our dress design, featuring diagonal lines and cutouts perfect for hot weather.',
        '100% cotton', true),
       (19, 'pouch-1', '/images/products/19/1.png',
        '/images/products/19/1.png;/images/products/19/2.png',
        'pouch',
        'pouch',
        '100% cotton', true),
       (20, 'embroidered-silk-pouch', '/images/products/20/1.png',
        '/images/products/20/1.png;/images/products/20/2.png',
        'Embroidered silk pouch',
        'An elegant accessory for your neck or as a slender belt, this embroidered creamy pouch is made from sheer silk organza with intricate detailing.',
        '100% silk', true),
       (21, 'pouch', '/images/products/21/1.png',
        '/images/products/21/1.png;/images/products/21/2.png',
        'pouch',
        'pouch',
        '100% cotton', true),
       (22, 'textured-floral-top', '/images/products/22/1.png',
        '/images/products/22/1.png;/images/products/22/2.png;/images/products/22/3.png;/images/products/22/4.png',
        'Textured floral top',
        'A breezy cotton top adorned with a textured floral pattern and charming cream-colored mother-of-pearl buttons.',
        '100% cotton', true),
       (23, 'sheer-dress-with-twisted-straps', '/images/products/23/1.png',
        '/images/products/23/1.png;/images/products/23/2.png;/images/products/23/3.png;/images/products/23/4.png',
        'Sheer dress with twisted straps',
        'A dress of semi-sheer vanilla viscose with twisted straps and an additional mini-length slip dress on top.',
        '100% silk', true),
       (24, 'ribbon-dress-with-pearls', '/images/products/24/1.png',
        '/images/products/24/1.png;/images/products/24/2.png;/images/products/24/3.png;/images/products/24/4.png',
        'Ribbon dress with pearls',
        'A loose-fitting dress with deep front and back necklines in milky cotton, featuring a wide ruffle and large baroque pearls hanging from delicate ribbons.',
        '100% silk', true),
       (25, 'weightless-blouse-with-watercolor-print', '/images/products/25/1.png',
        '/images/products/25/1.png;/images/products/25/2.png;/images/products/25/3.png;/images/products/25/4.png',
        'Weightless blouse with watercolor print',
        'A feather-light, floaty blouse crafted from the finest silk, decorated with a dreamy watercolor candy print and delicate ribbons.',
        '100% silk', true),
       (26, 'puff-sleeve-top', '/images/products/26/1.png',
        '/images/products/26/1.png;/images/products/26/2.png;/images/products/26/3.png;/images/products/26/4.png',
        'Puff sleeve top',
        'An ultra-voluminous, cloud-like top with lush gathers around the neckline, available in a sweet biscuit shade.',
        '100% silk', true);

INSERT INTO product_variants (product_id, id, name, available)
VALUES (15, 31, 'XS-S', 10),
       (15, 32, 'M-L', 10),
       (16, 33, 'XS-S', 10),
       (16, 34, 'M-L', 10),
       (17, 35, 'XS-S', 10),
       (17, 36, 'M-L', 10),
       (18, 37, 'XS-S', 10),
       (18, 38, 'M-L', 10),
       (19, 39, 'XS-S', 10),
       (19, 40, 'M-L', 10),
       (20, 41, 'XS-S', 10),
       (20, 42, 'M-L', 10),
       (21, 43, 'XS-S', 10),
       (21, 44, 'M-L', 10),
       (22, 45, 'XS-S', 10),
       (22, 46, 'M-L', 10),
       (23, 47, 'XS-S', 10),
       (23, 48, 'M-L', 10),
       (24, 49, 'XS-S', 10),
       (24, 50, 'M-L', 10),
       (25, 51, 'XS-S', 10),
       (25, 52, 'M-L', 10),
       (26, 53, 'XS-S', 10),
       (26, 54, 'M-L', 10);


INSERT INTO product_prices (product_id, id, price, currency_code)
VALUES (15, 29, 60, 'USD'),
       (15, 30, 264, 'BYN'),
       (16, 31, 250, 'USD'),
       (16, 32, 820, 'BYN'),
       (17, 33, 220, 'USD'),
       (17, 34, 730, 'BYN'),
       (18, 35, 135, 'USD'),
       (18, 36, 430, 'BYN'),
       (19, 37, 250, 'USD'),
       (19, 38, 820, 'BYN'),
       (20, 39, 60, 'USD'),
       (20, 40, 405, 'BYN'),
       (21, 41, 150, 'USD'),
       (21, 42, 400, 'BYN'),
       (22, 43, 250, 'USD'),
       (22, 44, 880, 'BYN'),
       (23, 45, 280, 'USD'),
       (23, 46, 990, 'BYN'),
       (24, 47, 380, 'USD'),
       (24, 48, 1150, 'BYN'),
       (25, 49, 225, 'USD'),
       (25, 50, 735, 'BYN'),
       (26, 51, 250, 'USD'),
       (26, 52, 840, 'BYN');


INSERT INTO product_translations (product_id, id, name, description, materials, language)
VALUES (15, 29, 'Шёлковый мешочек на тонких ремешках',
        'Небольшой шёлковый мешочек с маленькой розочкой для дополнения летних образов.', '100% шелк', 'ru'),
       (16, 30, 'Свободная рубашка на лентах',
        'Рубашка из шуршащего шелка в оттенке фисташкового мороженого на тонких лентах вместо классических пуговиц.',
        '100% шелк', 'ru'),
       (17, 31, 'Платье в пудрово - розовом оттенке',
        'Платье из плотного пудрово - розового летнего хлопка с дополняющей край вишневой лентой.', '100% шелк', 'ru'),
       (18, 32, 'Топ с диагональным рызрезом',
        'Топ из тонкого летнего хлопка по образцу нашего платья с диагональными линиями и вырезами для жаркого лета.',
        '100% хлопок', 'ru'),
       (19, 33, 'Подсумок', 'Подсумок', '100% хлопок', 'ru'),
       (20, 34, 'Шелковый мешочек с вышивкой',
        'Дополнительный аксессуар на шею или в виде тонкого пояса - Расшитый тонким узором сливочный мешочек из прозрачной шёлковой органзы.',
        '100% шелк', 'ru'),
       (21, 35, 'Подсумок', 'Подсумок', '100% хлопок', 'ru'),
       (22, 36, 'Топ с фактурным цветочным узором',
        'Свободный топ с лифом из хлопка с фактурным цветочным узором с застежкой на пуговицы из кремового перламутра.',
        '100% хлопок', 'ru'),
       (23, 37, 'Полупрозрачное платье с перевернутыми бретелями',
        'Свободное платье с глубокими вырезами спереди и сзади в молочном хлопке с широким воланом и  крупными барочными жемчужинами подвешенными на тонкие ленты. ',
        '100% шелк', 'ru'),
       (24, 38, 'Платье на лентах с жемчужинами',
        'Свободный ультра-пышный топ-облако с густой сборкой вокруг шеи в бисквитном оттенке.', '100% шелк', 'ru'),
       (25, 39, 'Невесомая блуза с акварельным рисунком',
        'Невесомая летящая блуза из тончайшего шёлка с акварельным конфетным рисунком на тонких лентах.', '100% шелк',
        'ru'),
       (26, 40, 'Топ с пышными рукавами',
        'Свободный ультра-пышный топ-облако с густой сборкой вокруг шеи в бисквитном оттенке.', '100% шелк', 'ru');

INSERT INTO discounts (id, value, code, type, is_active)
VALUES (2, 20, 'XH832KAY', 'percentage', true);
