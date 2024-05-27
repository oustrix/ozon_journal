CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at BIGINT NOT NULL,
    author_id INTEGER NOT NULL,
    commentable BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    content VARCHAR(200) NOT NULL,
    author_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    published_at BIGINT NOT NULL,
    parent_comment_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_posts_published_at ON posts(published_at);