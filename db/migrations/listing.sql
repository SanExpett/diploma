CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email TEXT UNIQUE CHECK (LENGTH(email) <= 30) NOT NULL,
    avatar TEXT DEFAULT '' NOT NULL,
    name TEXT DEFAULT 'user' NOT NULL,
    password TEXT CHECK (LENGTH(password) <= 64) NOT NULL,
    registered_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    birthday TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE NOT NULL,
    subscription_end_date TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    amount DECIMAL NOT NULL,
    description TEXT NOT NULL,
    duration INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS actor
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT UNIQUE NOT NULL,
    avatar TEXT DEFAULT 'https://shorturl.at/ewzP8' NOT NULL,
    birthday TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    career TEXT DEFAULT '' NOT NULL,
    height INTEGER CHECK (height < 300) DEFAULT 192 NOT NULL,
    birth_place TEXT DEFAULT 'Russia, Angarsk' NOT NULL,
    spouse TEXT DEFAULT 'Светлана Ходченкова' NOT NULL
);

CREATE TABLE IF NOT EXISTS director
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT UNIQUE NOT NULL,
    avatar TEXT DEFAULT 'https://shorturl.at/ewzP8' NOT NULL,
    birthday TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS film
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    is_serial BOOLEAN DEFAULT FALSE NOT NULL,
    title TEXT NOT NULL,
    data TEXT DEFAULT '' NOT NULL,
    banner TEXT DEFAULT 'https://shorturl.at/akMR2' NOT NULL,
    s3_link TEXT DEFAULT 'https://shorturl.at/jHIMO' NOT NULL,
    director INTEGER,
    age_limit SMALLINT CHECK (age_limit >= 0 AND age_limit <= 18) DEFAULT 18 NOT NULL,
    duration SMALLINT CHECK (duration > 0) DEFAULT 143 NOT NULL,
    published_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    with_subscription BOOLEAN DEFAULT FALSE NOT NULL,
    FOREIGN KEY (director) REFERENCES director (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS season
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    film_id INTEGER NOT NULL,
    number INTEGER NOT NULL,
    FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS episode
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number INTEGER NOT NULL,
    title TEXT NOT NULL,
    s3_link TEXT DEFAULT 'https://shorturl.at/jHIMO' NOT NULL,
    season_id INTEGER NOT NULL,
    FOREIGN KEY (season_id) REFERENCES season (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    text TEXT NOT NULL,
    score SMALLINT CHECK (score >= 0 AND score <= 10) NOT NULL,
    author_id INTEGER NOT NULL,
    film_id INTEGER NOT NULL,
    added_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS film_actor
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    film_id INTEGER NOT NULL,
    actor_id INTEGER NOT NULL,
    FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actor (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS favorite_film
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    film_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS genre
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS film_genre
(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    film_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES genre (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE
);

-- Индекс на email пользователей
-- Ускоряет авторизацию и проверку уникальности email при регистрации
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Индексы для фильмов
-- Оптимизирует поиск фильмов по названию в поисковой строке
CREATE INDEX IF NOT EXISTS idx_film_title ON film(title);
-- Ускоряет фильтрацию контента по наличию подписки и типу (сериал/фильм)
CREATE INDEX IF NOT EXISTS idx_film_subscription_serial ON film(with_subscription, is_serial);

-- Индекс для комментариев
-- Оптимизирует получение отсортированных по дате комментариев к конкретному фильму
CREATE INDEX IF NOT EXISTS idx_comment_film_date ON comment(film_id, added_at DESC);

-- Индекс для избранного
-- Ускоряет получение списка избранных фильмов пользователя и предотвращает дублирование
CREATE UNIQUE INDEX IF NOT EXISTS idx_favorite_film_user_film ON favorite_film(user_id, film_id);

-- Индекс для актеров
-- Оптимизирует поиск актеров по имени в поисковой строке
CREATE INDEX IF NOT EXISTS idx_actor_name ON actor(name);

-- Индексы для сериалов
-- Ускоряет получение списка сезонов для конкретного сериала
CREATE INDEX IF NOT EXISTS idx_season_film_number ON season(film_id, number);
-- Оптимизирует получение списка серий для конкретного сезона
CREATE INDEX IF NOT EXISTS idx_episode_season_number ON episode(season_id, number);