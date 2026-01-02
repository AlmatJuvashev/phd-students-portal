CREATE TYPE forum_type AS ENUM ('ANNOUNCEMENT', 'QNA', 'DISCUSSION');

CREATE TABLE forums (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    forum_type forum_type NOT NULL DEFAULT 'DISCUSSION',
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    forum_id UUID NOT NULL REFERENCES forums(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    views_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES posts(id) ON DELETE CASCADE, -- For nested replies
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_forums_course ON forums(course_offering_id);
CREATE INDEX idx_topics_forum ON topics(forum_id);
CREATE INDEX idx_posts_topic ON posts(topic_id);
CREATE INDEX idx_posts_parent ON posts(parent_id);
