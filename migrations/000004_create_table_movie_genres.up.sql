CREATE TABLE movie_genres (
  id SERIAL PRIMARY KEY,
  id_movie INT REFERENCES movies(id) ON DELETE CASCADE,
  id_genre INT REFERENCES genres(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
