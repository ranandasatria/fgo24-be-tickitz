ALTER TABLE genres ADD CONSTRAINT genres_name_unique UNIQUE (genre_name);
ALTER TABLE directors ADD CONSTRAINT directors_name_unique UNIQUE (director_name);
ALTER TABLE actors ADD CONSTRAINT actors_name_unique UNIQUE (actor_name);
