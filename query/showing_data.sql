SELECT 
m.id, m.title as movie, 
d.director_name as director, 
a.actor_name as main_cast, 
mc.role_name as role, 
STRING_AGG(g.genre_name, ', ') as genre 
FROM movies m
JOIN movie_casts mc on mc.id_movie = m.id
JOIN actors a on a.id = mc.id_actor
JOIN movie_directors md on md.id_movie = m.id
JOIN directors d on d.id = md.id_director
JOIN movie_genres mg on mg.id_movie = m.id
JOIN genres g on g.id = mg.id_genre
GROUP BY m.id, m.title, d.director_name, a.actor_name, mc.role_name;

-- Total ticket per user
SELECT 
u.id as id_user,
u.full_name as username,
COUNT(*) as total_ticket,
STRING_AGG(DISTINCT m.title, ', ') as watched_movies
FROM tickets t
JOIN users u on u.id = t.id_user
JOIN movies m on m.id = t.id_movie
GROUP BY u.id, u.full_name
ORDER BY total_ticket DESC;

-- Total ticket sold per movie
SELECT 
m.title, 
COUNT(*) AS total_watched 
FROM tickets t
JOIN movies m ON m.id = t.id_movie
GROUP BY m.title
ORDER BY total_watched DESC;

-- Income per movie
SELECT 
m.title, 
SUM(t.price_per_ticket) as total_income 
FROM tickets t
JOIN movies m ON m.id = t.id_movie
GROUP BY m.title
ORDER BY total_income DESC;

-- Total transaction per payment_method
SELECT 
p.payment_name, 
COUNT(*) as total_transactions 
FROM tickets t
JOIN payment_method p on p.id = t.payment_method
GROUP BY p.payment_name
ORDER BY total_transactions DESC;

-- Total viewer per location
SELECT 
t.location,
COUNT(*) AS total_viewers
FROM tickets t
GROUP BY t.location
ORDER BY total_viewers DESC;

