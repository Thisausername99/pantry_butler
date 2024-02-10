CREATE TABLE pantry_entry (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  quantity INTEGER,
  quantity_type VARCHAR,
  created_at TIMESTAMP
);

INSERT INTO pantry_entry (name) values ('onion'); 
INSERT INTO pantry_entry (name) values ('rice');
INSERT INTO pantry_entry (name) values ('butter'); 
INSERT INTO pantry_entry (name) values ('garlic'); 