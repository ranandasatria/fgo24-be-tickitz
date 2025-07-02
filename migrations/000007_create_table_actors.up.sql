CREATE TABLE actors (
  id SERIAL PRIMARY KEY,
  actor_name VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
