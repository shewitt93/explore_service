CREATE TABLE `user` (
     id INT PRIMARY KEY,
     email TEXT NOT NULL,
     name TEXT NOT NULL


);


CREATE TABLE user_decisions (
    actor_id VARCHAR(255) NOT NULL,
    recipient_id VARCHAR(255) NOT NULL,
    liked BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (actor_id, recipient_id),
    INDEX idx_recipient_liked (recipient_id, liked),
    INDEX idx_recipient_updated (recipient_id, updated_at, actor_id)

);