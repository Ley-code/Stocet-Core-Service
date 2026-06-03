CREATE TABLE article_likes (
    article_id  uuid        NOT NULL,
    user_id     uuid        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (article_id, user_id)
);

CREATE TABLE article_comments (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id  uuid        NOT NULL,
    user_id     uuid        NOT NULL,
    user_name   text        NOT NULL,
    content     text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_article_comments_article ON article_comments (article_id, created_at DESC);
