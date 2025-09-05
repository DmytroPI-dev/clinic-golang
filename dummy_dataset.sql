-- dummy_data.sql

INSERT INTO `programs` 
(
    `created_at`, `updated_at`, `deleted_at`, 
    `title`, `description`, `results`, 
    `title_pl`, `description_pl`, `results_pl`, 
    `title_en`, `description_en`, `results_en`, 
    `title_uk`, `description_uk`, `results_uk`, 
    `category`
) 
VALUES 
(
    NOW(), NOW(), NULL,
    'Facial Cleansing', 'A deep cleansing facial treatment.', 'Clean and refreshed skin.',
    'Oczyszczanie Twarzy', 'Głęboko oczyszczający zabieg na twarz.', 'Czysta i odświeżona skóra.',
    'Facial Cleansing', 'A deep cleansing facial treatment.', 'Clean and refreshed skin.',
    'Очищення обличчя', 'Глибоко очищувальна процедура для обличчя.', 'Чиста та оновлена шкіра.',
    'KS'
);