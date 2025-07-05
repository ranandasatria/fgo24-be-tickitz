DROP TABLE IF EXISTS tickets;

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  id_user INT REFERENCES users(id),
  id_movie INT REFERENCES movies(id),
  show_date DATE NOT NULL,
  show_time TIME NOT NULL,
  location VARCHAR(255) NOT NULL,
  cinema VARCHAR(255) NOT NULL,
  total_price INT NOT NULL,
  payment_method INT REFERENCES payment_method(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);