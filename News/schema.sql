
CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pubtime BIGINT DEFAULT 0,
    link TEXT NOT NULL
);