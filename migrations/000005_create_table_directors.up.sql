CREATE TABLE directors (
  id SERIAL PRIMARY KEY,
  director_name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
