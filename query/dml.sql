INSERT INTO users (email, password, full_name, phone_number, profile_picture) VALUES
('ali@example.com', 'pass1', 'Ali Rahman', 628111111111, 'ali.jpg'),
('budi@example.com', 'pass2', 'Budi Santoso', 628222222222, 'budi.jpg'),
('citra@example.com', 'pass3', 'Citra Dewi', 628333333333, 'citra.jpg'),
('dina@example.com', 'pass4', 'Dina Kurnia', 628444444444, 'dina.jpg'),
('eko@example.com', 'pass5', 'Eko Prasetyo', 628555555555, 'eko.jpg');

INSERT INTO sessions (id_user, login_time) VALUES
(1, CURRENT_TIMESTAMP),
(2, CURRENT_TIMESTAMP),
(3, CURRENT_TIMESTAMP),
(4, CURRENT_TIMESTAMP),
(5, CURRENT_TIMESTAMP);

INSERT INTO movies (title, description, release_date, duration_minutes, image, horizontal_image) VALUES
('Where is the Friend''s House?', 'A boy searches for his friend''s home to return a notebook.', '1987-03-20', 83, 'friend.jpg', 'friend-h.jpg'),
('Bicycle Thieves', 'A man and his son search for a stolen bicycle.', '1948-11-24', 89, 'bicycle.jpg', 'bicycle-h.jpg'),
('High and Low', 'A businessman is involved in a kidnapping case.', '1963-06-23', 143, 'highlow.jpg', 'highlow-h.jpg'),
('Paris, Texas', 'A man reconnects with his son after years of disappearance.', '1984-05-19', 147, 'paris.jpg', 'paris-h.jpg'),
('A Woman Under the Influence', 'A woman struggles with mental health and family.', '1974-11-18', 155, 'woman.jpg', 'woman-h.jpg');

INSERT INTO genres (genre_name) VALUES
('Drama'),
('Family'),
('Crime'),
('Psychological'),
('Road Movie');

INSERT INTO movie_genres (id_movie, id_genre) VALUES
(1, 1),
(1, 2), 
(2, 1),
(2, 2),
(3, 1),
(3, 3), 
(4, 1),
(4, 5), 
(5, 1),
(5, 4);

INSERT INTO directors (director_name) VALUES
('Abbas Kiarostami'),
('Vittorio De Sica'),
('Akira Kurosawa'),
('Wim Wenders'),
('John Cassavetes');

INSERT INTO movie_directors (id_movie, id_director) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5);

INSERT INTO actors (actor_name) VALUES
('Babek Ahmed Poor'),
('Lamberto Maggiorani'),
('Toshiro Mifune'),
('Harry Dean Stanton'),
('Gena Rowlands');

INSERT INTO movie_casts (id_movie, id_actor, role_name) VALUES
(1, 1, 'Ahmad'),
(2, 2, 'Antonio Ricci'),
(3, 3, 'King Gondo'),
(4, 4, 'Travis'),
(5, 5, 'Mabel Longhetti');

INSERT INTO payment_method (payment_name) VALUES
('BCA'),
('Mandiri'),
('OVO'),
('Gopay'),
('ShopeePay');

INSERT INTO tickets (id_user, id_movie, show_date, show_time, cinema, location, seat, price_per_ticket, payment_method) VALUES
(1, 1, '2025-06-25', '17:00', 'CGV', 'Jakarta', 'A1', 50000, 1),
(2, 2, '2025-06-26', '19:00', 'XXI', 'Bandung', 'B2', 55000, 2),
(3, 3, '2025-06-27', '18:30', 'CGV', 'Surabaya', 'C3', 60000, 3),
(4, 4, '2025-06-28', '20:00', 'XXI', 'Yogyakarta', 'D4', 65000, 4),
(5, 5, '2025-06-29', '21:00', 'Cinepolis', 'Semarang', 'E5', 70000, 5);

INSERT INTO tickets (id_user, id_movie, show_date, show_time, cinema, location, seat, price_per_ticket, payment_method) VALUES
(1, 2, '2025-07-01', '17:30', 'CGV', 'Jakarta', 'A2', 55000, 2),
(1, 3, '2025-07-02', '18:00', 'XXI', 'Jakarta', 'A3', 60000, 3),
(2, 3, '2025-07-03', '20:00', 'Cinepolis', 'Bandung', 'B3', 60000, 1),
(3, 4, '2025-07-04', '19:00', 'XXI', 'Surabaya', 'C4', 65000, 2),
(4, 1, '2025-07-05', '16:00', 'CGV', 'Yogyakarta', 'D1', 50000, 4),
(5, 2, '2025-07-06', '20:30', 'CGV', 'Semarang', 'E2', 55000, 5),
(1, 5, '2025-07-07', '21:00', 'XXI', 'Jakarta', 'F5', 70000, 1),
(2, 4, '2025-07-08', '18:30', 'Cinepolis', 'Bandung', 'B4', 65000, 3),
(3, 5, '2025-07-09', '17:15', 'CGV', 'Surabaya', 'C5', 70000, 2),
(4, 3, '2025-07-10', '19:45', 'XXI', 'Yogyakarta', 'D6', 60000, 4);