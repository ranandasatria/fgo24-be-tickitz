CREATE TABLE payment_method (
  id SERIAL PRIMARY KEY,
  payment_name VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
