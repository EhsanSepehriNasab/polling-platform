CREATE TABLE IF NOT EXISTS polls (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    options TEXT[] NOT NULL,
    tags TEXT[] NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);