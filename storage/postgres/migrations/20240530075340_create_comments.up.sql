CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,

    post_id BIGINT REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    parent_id BIGINT REFERENCES comments (id),
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE NOT NULL
);