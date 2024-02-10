CREATE TABLE recipe (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    description VARCHAR,
    ingredients JSONB,
    rating INTEGER,
    cuisine VARCHAR,
    created_at TIMESTAMP
);

INSERT INTO recipe (name, description, ingredients) VALUES ('burger', 'description', '{"onion": "0", "ground beef": "1"}');

