CREATE TABLE movie_casts (
  id SERIAL PRIMARY KEY,
  id_movie INT REFERENCES movies(id) ON DELETE CASCADE,
  id_actor INT REFERENCES actors(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
