-- dummy_data.sql
INSERT INTO
    `programs` (
        `created_at`,
        `updated_at`,
        `deleted_at`,
        `title`,
        `description`,
        `results`,
        `title_pl`,
        `description_pl`,
        `results_pl`,
        `title_en`,
        `description_en`,
        `results_en`,
        `title_uk`,
        `description_uk`,
        `results_uk`,
        `category`
    )
VALUES
    (
        NOW(),
        NOW(),
        NULL,
        'Face Cleaning',
        'A deep cleansing facial treatment.',
        'Clean and refreshed skin.',
        'Oczyszczanie Twarzy',
        'Głęboko oczyszczający zabieg na twarz.',
        'Czysta i odświeżona skóra.',
        'Facial Cleansing',
        'A deep cleansing facial treatment.',
        'Clean and refreshed skin.',
        'Очищення обличчя',
        'Глибоко очищувальна процедура для обличчя.',
        'Чиста та оновлена шкіра.',
        'KS'
    );

INSERT INTO
    `prices` (
        `created_at`,
        `updated_at`,
        `item_name`,
        `price`,
        `item_name_pl`,
        `item_name_en`,
        `item_name_uk`,
        `category`
    )
VALUES
    (
        NOW(),
        NOW(),
        'Consultation',
        50.00,
        'Konsultacja',
        'Consultation',
        'Консультація',
        'KS'
    ),
    (
        NOW(),
        NOW(),
        'Laser Facial',
        150.00,
        'Zabieg Laserowy na Twarz',
        'Laser Facial',
        'Лазерна чистка обличчя',
        'LS'
    ),
    (
        NOW(),
        NOW(),
        'Manicure',
        80.00,
        'Manicure',
        'Manicure',
        'Манікюр',
        'KT'
    );

INSERT INTO
    `news` (
        `created_at`,
        `updated_at`,
        `title`,
        `description`,
        `header`,
        `features`,
        `image_left`,
        `image_right`,
        `posted_on`,
        `title_pl`,
        `description_pl`,
        `header_pl`,
        `features_pl`,
        `title_en`,
        `description_en`,
        `header_en`,
        `features_en`,
        `title_uk`,
        `description_uk`,
        `header_uk`,
        `features_uk`
    )
VALUES
    (
        NOW(),
        NOW(),
        'New Spring Promotions',
        'Check out our new promotions for the spring season!',
        'Spring Sale',
        'Discount on all laser treatments.',
        'http://example.com/image1.jpg',
        'http://example.com/image2.jpg',
        '2025-09-01',
        'Nowe Promocje Wiosenne',
        'Sprawdź nasze nowe promocje na sezon wiosenny!',
        'Wiosenna Wyprzedaż',
        'Zniżka na wszystkie zabiegi laserowe.',
        'New Spring Promotions',
        'Check out our new promotions for the spring season!',
        'Spring Sale',
        'Discount on all laser treatments.',
        'Нові весняні акції',
        'Ознайомтеся з нашими новими акціями на весняний сезон!',
        'Весняний розпродаж',
        'Знижка на всі лазерні процедури.'
    ),
    (
        NOW(),
        NOW(),
        'We Are Open on Saturdays',
        'We are happy to announce we are now open on Saturdays.',
        'New Opening Hours',
        'Longer hours for your convenience.',
        'http://example.com/image3.jpg',
        'http://example.com/image4.jpg',
        '2025-08-15',
        'Jesteśmy Otwarci w Soboty',
        'Z przyjemnością informujemy, że jesteśmy teraz otwarci w soboty.',
        'Nowe Godziny Otwarcia',
        'Dłuższe godziny dla Twojej wygody.',
        'We Are Open on Saturdays',
        'We are happy to announce we are now open on Saturdays.',
        'New Opening Hours',
        'Longer hours for your convenience.',
        'Ми працюємо по суботах',
        'Раді повідомити, що тепер ми працюємо по суботах.',
        'Новий графік роботи',
        'Довший робочий день для вашої зручності.'
    );