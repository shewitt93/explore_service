-- Clear existing data (if any)
DELETE FROM user_decisions;
DELETE FROM user;

-- Insert sample users
INSERT INTO user (id, email, name) VALUES
(1, 'john@example.com', 'John Smith'),
(2, 'sarah@example.com', 'Sarah Johnson'),
(3, 'mike@example.com', 'Mike Williams'),
(4, 'emily@example.com', 'Emily Brown'),
(5, 'david@example.com', 'David Lee'),
(6, 'lisa@example.com', 'Lisa Garcia'),
(7, 'james@example.com', 'James Wilson'),
(8, 'jessica@example.com', 'Jessica Martinez'),
(9, 'robert@example.com', 'Robert Taylor'),
(10, 'jennifer@example.com', 'Jennifer Anderson');

-- Insert sample decisions
-- User 1 likes several users
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('1', '2', TRUE, '2025-01-15 10:30:00', '2025-01-15 10:30:00'),
('1', '3', TRUE, '2025-01-16 11:45:00', '2025-01-16 11:45:00'),
('1', '4', FALSE, '2025-01-17 09:15:00', '2025-01-17 09:15:00'),
('1', '5', TRUE, '2025-01-18 14:20:00', '2025-01-18 14:20:00');

-- User 2 likes several users including User 1 (mutual like)
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('2', '1', TRUE, '2025-01-20 16:30:00', '2025-01-20 16:30:00'),
('2', '3', FALSE, '2025-01-21 09:45:00', '2025-01-21 09:45:00'),
('2', '4', TRUE, '2025-01-22 11:10:00', '2025-01-22 11:10:00');

-- User 3 likes User 1 (mutual like)
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('3', '1', TRUE, '2025-01-23 08:30:00', '2025-01-23 08:30:00'),
('3', '2', TRUE, '2025-01-24 10:15:00', '2025-01-24 10:15:00');

-- User 4 doesn't like User 1
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('4', '1', FALSE, '2025-01-25 14:40:00', '2025-01-25 14:40:00'),
('4', '2', TRUE, '2025-01-26 16:20:00', '2025-01-26 16:20:00'),
('4', '3', TRUE, '2025-01-27 11:30:00', '2025-01-27 11:30:00');

-- User 5 likes User 1 but User 1 already liked User 5 (mutual like)
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('5', '1', TRUE, '2025-01-28 09:50:00', '2025-01-28 09:50:00'),
('5', '2', FALSE, '2025-01-29 13:25:00', '2025-01-29 13:25:00');

-- More users like User 1
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('6', '1', TRUE, '2025-02-01 10:10:00', '2025-02-01 10:10:00'),
('7', '1', TRUE, '2025-02-02 15:45:00', '2025-02-02 15:45:00'),
('8', '1', FALSE, '2025-02-03 12:30:00', '2025-02-03 12:30:00'),
('9', '1', TRUE, '2025-02-04 16:20:00', '2025-02-04 16:20:00'),
('10', '1', TRUE, '2025-02-05 11:15:00', '2025-02-05 11:15:00');

-- Add a few additional connections
INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES
('6', '7', TRUE, '2025-02-06 14:30:00', '2025-02-06 14:30:00'),
('7', '6', TRUE, '2025-02-07 09:20:00', '2025-02-07 09:20:00'),
('8', '9', TRUE, '2025-02-08 10:45:00', '2025-02-08 10:45:00'),
('9', '8', FALSE, '2025-02-09 13:10:00', '2025-02-09 13:10:00'),
('10', '5', TRUE, '2025-02-10 15:30:00', '2025-02-10 15:30:00');

-- Select statements to verify the data
SELECT 'Users:' as '';
SELECT * FROM user;

SELECT 'Decisions:' as '';
SELECT * FROM user_decisions;

SELECT 'Users who liked User 1:' as '';
SELECT u.name, d.liked, d.created_at 
FROM user_decisions d 
JOIN user u ON d.actor_id = CAST(u.id AS CHAR) 
WHERE d.recipient_id = '1' AND d.liked = TRUE;

SELECT 'Mutual likes for User 1:' as '';
SELECT u.name 
FROM user_decisions d1
JOIN user_decisions d2 ON d1.actor_id = d2.recipient_id AND d1.recipient_id = d2.actor_id
JOIN user u ON CAST(u.id AS CHAR) = d1.recipient_id
WHERE d1.actor_id = '1' AND d2.actor_id = d1.recipient_id 
AND d1.liked = TRUE AND d2.liked = TRUE;
