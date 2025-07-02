CREATE TABLE tickets (
  id SERIAL PRIMARY KEY,
  id_user INT REFERENCES users(id),
  id_movie INT REFERENCES movies(id),
  show_date DATE,
  show_time TIME,
  cinema VARCHAR(255),
  location VARCHAR(255),
  seat VARCHAR(255),
  price_per_ticket INT,
  payment_method INT REFERENCES payment_method(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
