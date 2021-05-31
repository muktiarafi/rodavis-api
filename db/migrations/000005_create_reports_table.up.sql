CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    status status NOT NULL DEFAULT 'Reported',
    image_url VARCHAR(255) NOT NULL,
    classes class[] NOT NULL,
    note TEXT NOT NULL,
    address VARCHAR(255) NOT NULL,
    lat NUMERIC NOT NULL,
    lng NUMERIC NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    date_reported TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);