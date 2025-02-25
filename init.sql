
CREATE TABLE user_decisions (
        id SERIAL PRIMARY KEY,
        actor_id TEXT NOT NULL,
        recipient_id TEXT NOT NULL,
        liked BOOLEAN NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(actor_id, recipient_id)
);