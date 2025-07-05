CREATE TABLE transaction_details (
  id SERIAL PRIMARY KEY,
  transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
  seat VARCHAR(10) NOT NULL,
  price INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
