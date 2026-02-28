-- Сидеры: тестовые данные для job_stats
-- Idempotent: безопасно запускать повторно — данные очищаются перед вставкой

SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE job_skills;
TRUNCATE TABLE locations;
TRUNCATE TABLE jobs;
TRUNCATE TABLE skills;
TRUNCATE TABLE companies;
SET FOREIGN_KEY_CHECKS = 1;

-- Компании
INSERT INTO companies (name, description) VALUES
    ('VK', 'Социальная сеть и технологическая компания'),
    ('Yandex Infrastructure', 'Подразделение Яндекса, занимающееся разработкой облачной инфраструктуры'),
    ('P2P.ORG', 'Крупнейший институциональный провайдер стейкинга с TVL более $10B'),
    ('Яндекс Foodtech', 'Сервис доставки еды и продуктов Яндекс Еда'),
    ('Ozon Банк', 'Специальное подразделение Ozon в области финансов и IT. Банк, карта, программа лояльности, рассрочка, B2B-продукты');

-- Навыки - Языки программирования
INSERT INTO skills (name, category) VALUES
    ('Golang', 'Язык программирования'),
    ('Go', 'Язык программирования'),
    ('Python', 'Язык программирования'),
    ('Java', 'Язык программирования'),
    ('C++', 'Язык программирования'),
    ('PHP', 'Язык программирования'),
    ('JavaScript', 'Язык программирования'),
    ('TypeScript', 'Язык программирования'),
    ('Rust', 'Язык программирования'),
    ('Solidity', 'Язык программирования');

-- Навыки - Базы данных
INSERT INTO skills (name, category) VALUES
    ('PostgreSQL', 'База данных'),
    ('MySQL', 'База данных'),
    ('MongoDB', 'База данных'),
    ('Redis', 'База данных'),
    ('ClickHouse', 'База данных'),
    ('Cassandra', 'База данных'),
    ('Elasticsearch', 'База данных');

-- Навыки - Инструменты
INSERT INTO skills (name, category) VALUES
    ('Kafka', 'Инструмент'),
    ('Docker', 'Инструмент'),
    ('Kubernetes', 'Инструмент'),
    ('Git', 'Инструмент'),
    ('GitHub', 'Инструмент'),
    ('GitLab', 'Инструмент'),
    ('Jenkins', 'Инструмент'),
    ('CI/CD', 'Инструмент'),
    ('GCP', 'Инструмент'),
    ('AWS', 'Инструмент'),
    ('YTsaurus', 'Инструмент'),
    ('gRPC', 'Инструмент'),
    ('RabbitMQ', 'Инструмент'),
    ('GraphQL', 'Инструмент'),
    ('REST', 'Инструмент'),
    ('Pytest', 'Инструмент'),
    ('Allure', 'Инструмент'),
    ('GitLab CI', 'Инструмент');

-- Навыки - Фреймворки
INSERT INTO skills (name, category) VALUES
    ('React', 'Фреймворк'),
    ('Vue.js', 'Фреймворк'),
    ('Angular', 'Фреймворк'),
    ('Django', 'Фреймворк'),
    ('FastAPI', 'Фреймворк'),
    ('Spring', 'Фреймворк');

-- Навыки - Другое
INSERT INTO skills (name, category) VALUES
    ('Высоконагруженные системы', 'Другое'),
    ('Микросервисы', 'Другое'),
    ('Blockchain', 'Другое'),
    ('DeFi', 'Другое'),
    ('Smart Contracts', 'Другое'),
    ('EVM', 'Другое'),
    ('Отказоустойчивость', 'Другое'),
    ('Алгоритмы и структуры данных', 'Другое'),
    ('Распределенные системы', 'Другое'),
    ('E-commerce', 'Другое'),
    ('Ритейл', 'Другое'),
    ('Логистика', 'Другое'),
    ('Low-latency системы', 'Другое'),
    ('Event-driven архитектура', 'Другое'),
    ('Автоматизация тестирования', 'Другое'),
    ('QA', 'Другое'),
    ('Финтех', 'Другое');

-- Вакансия 1: VK - Ведущий Go-разработчик
INSERT INTO jobs (company_id, title, level, specialization, salary_min, salary_max, salary_currency,
    experience_years, location, remote_available, description, responsibilities, benefits, is_active, posted_date)
VALUES (
    1,
    'Ведущий Go-разработчик',
    'Senior',
    'Golang',
    370000,
    520000,
    'RUB',
    '6+ лет',
    'Москва, Санкт-Петербург',
    TRUE,
    'Работа над самой масштабной контентной платформой в рунете',
    'Проектирование Go-сервисов, оптимизация нагрузки, переход на open-source БД',
    'ДМС, компенсация питания 800₽/день, компенсация спорта 35000₽/год',
    TRUE,
    CURDATE()
);

-- Навыки для вакансии 1
INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 1, id, TRUE, FALSE FROM skills WHERE name IN ('Golang', 'PostgreSQL', 'MongoDB', 'Redis', 'ClickHouse', 'Kafka', 'Высоконагруженные системы');

INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 1, id, FALSE, TRUE FROM skills WHERE name IN ('PHP', 'YTsaurus');

-- Локации для вакансии 1
INSERT INTO locations (job_id, city, metro_station, is_primary) VALUES
    (1, 'Москва', 'Аэропорт', TRUE),
    (1, 'Санкт-Петербург', 'Невский проспект', FALSE);

-- Вакансия 2: Yandex Infrastructure - Go разработчик
INSERT INTO jobs (company_id, title, level, specialization, salary_min, salary_max, salary_currency,
    experience_years, location, remote_available, description, responsibilities, benefits, is_active, posted_date)
VALUES (
    2,
    'Разработчик на Go (Yandex BareMetal)',
    'Senior',
    'Golang',
    410000,
    600000,
    'RUB',
    '5+ лет',
    'Удаленно из РФ',
    TRUE,
    'Разработка сервиса аренды выделенных физических серверов в облаке',
    'Разработка core-компонентов, проектирование архитектуры, оптимизация производительности',
    'ДМС с первого месяца, психотерапия, корпоративный спорт, гибкий график',
    TRUE,
    CURDATE()
);

-- Навыки для вакансии 2
INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 2, id, TRUE, FALSE FROM skills WHERE name IN ('Go', 'Распределенные системы');

INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 2, id, FALSE, TRUE FROM skills WHERE name IN ('Kubernetes', 'GCP');

-- Локации для вакансии 2
INSERT INTO locations (job_id, city, is_primary) VALUES
    (2, 'Удаленно', TRUE);

-- Вакансия 3: P2P.ORG - Backend Engineer (Blockchain)
INSERT INTO jobs (company_id, title, level, specialization, salary_min, salary_max, salary_currency,
    experience_years, location, remote_available, description, responsibilities, benefits, is_active, posted_date)
VALUES (
    3,
    'Backend Engineer (Blockchain)',
    'Senior',
    'Blockchain & Crypto',
    6000,
    9000,
    'USD',
    '5+ лет',
    'Полная удалёнка',
    TRUE,
    'Разработка backend-сервисов для DeFi/on-chain продуктов',
    'Создание протокольных адаптеров, оптимизация RPC-запросов, работа с on-chain данными',
    'Полностью удаленная работа, зарплата в $ или крипте, компенсация образования',
    TRUE,
    CURDATE()
);

-- Навыки для вакансии 3
INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 3, id, TRUE, FALSE FROM skills
WHERE name IN ('Python', 'PostgreSQL', 'ClickHouse', 'Redis', 'Kubernetes', 'Docker', 'Blockchain', 'DeFi', 'Smart Contracts');

INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 3, id, FALSE, TRUE FROM skills WHERE name IN ('Solidity');

-- Локации для вакансии 3
INSERT INTO locations (job_id, city, is_primary) VALUES
    (3, 'Удаленно (European time zone)', TRUE);

-- Вакансия 4: Яндекс Foodtech - C++/Golang разработчик
INSERT INTO jobs (company_id, title, level, specialization, salary_min, salary_max, salary_currency,
    experience_years, location, remote_available, description, responsibilities, benefits, is_active, posted_date)
VALUES (
    4,
    'С++/Golang-разработчик (Еда)',
    'Middle',
    'С++ / Golang',
    300000,
    450000,
    'RUB',
    '4+ лет',
    'Москва, Санкт-Петербург',
    TRUE,
    'Разработка новых сценариев в направлении ритейла Яндекс Еды',
    'Разработка высоконагруженных микросервисов, улучшение качества кода',
    'ДМС, психотерапия, корпоративный спорт, гибкий график, жилищные займы',
    TRUE,
    CURDATE()
);

-- Навыки для вакансии 4
INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 4, id, TRUE, FALSE FROM skills
WHERE name IN ('C++', 'Golang', 'Микросервисы', 'Высоконагруженные системы', 'Алгоритмы и структуры данных');

INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have)
SELECT 4, id, FALSE, TRUE FROM skills WHERE name IN ('E-commerce', 'Ритейл', 'Логистика');

-- Локации для вакансии 4
INSERT INTO locations (job_id, city, metro_station, is_primary) VALUES
    (4, 'Москва', 'Деловой центр', TRUE),
    (4, 'Санкт-Петербург', 'Площадь Ленина', FALSE);
