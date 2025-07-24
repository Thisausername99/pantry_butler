CREATE TABLE pantry_entries (
  id SERIAL PRIMARY KEY,
  name VARCHAR UNIQUE NOT NULL,
  quantity INTEGER,
  quantity_type VARCHAR,
  expiration TIMESTAMP
);

INSERT INTO pantry_entries (name) values ('onion'); 
INSERT INTO pantry_entries (name) values ('rice');
INSERT INTO pantry_entries (name) values ('butter'); 
INSERT INTO pantry_entries (name) values ('garlic'); 