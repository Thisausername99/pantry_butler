CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    description VARCHAR,
    ingredients JSONB,
    rating INTEGER,
    cuisine VARCHAR,
    difficulty INTEGER,
    created_at TIMESTAMP
);

INSERT INTO recipes (name, description, cuisine, ingredients) VALUES ('burger', 'description', 'bach_kun', '{"onion": "0", "ground beef": "1"}');

