ALTER TABLE genres
  DROP CONSTRAINT IF EXISTS genres_name_unique;

ALTER TABLE directors
  DROP CONSTRAINT IF EXISTS directors_name_unique;

ALTER TABLE actors
  DROP CONSTRAINT IF EXISTS actors_name_unique;
