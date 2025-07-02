CREATE TABLE movie_directors (
  id SERIAL PRIMARY KEY,
  id_movie INT REFERENCES movies(id) ON DELETE CASCADE,
  id_director INT REFERENCES directors(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
