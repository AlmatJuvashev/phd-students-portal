DROP INDEX IF EXISTS idx_posts_parent;
DROP INDEX IF EXISTS idx_posts_topic;
DROP INDEX IF EXISTS idx_topics_forum;
DROP INDEX IF EXISTS idx_forums_course;

DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS forums;
DROP TYPE IF EXISTS forum_type;
