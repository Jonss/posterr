CREATE TABLE posts(
    id BIGSERIAL PRIMARY KEY,
    content TEXT NULL,
    user_id BIGSERIAL NOT NULL,
    original_post_id BIGINT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    constraint fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    constraint fk_posts_original_post_id FOREIGN KEY (original_post_id) REFERENCES posts (id)
);