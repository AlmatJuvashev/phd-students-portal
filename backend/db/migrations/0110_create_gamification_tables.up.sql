-- Gamification System Tables

-- 1. XP Levels Definition
CREATE TABLE xp_levels (
    level INT PRIMARY KEY,
    xp_required INT NOT NULL,
    title VARCHAR(50),
    perks JSONB
);

-- Seed Initial Levels
INSERT INTO xp_levels (level, xp_required, title, perks) VALUES
(1, 0, 'Newcomer', '{}'),
(2, 100, 'Learner', '{}'),
(3, 300, 'Student', '{}'),
(4, 600, 'Scholar', '{}'),
(5, 1000, 'Advanced', '{}'),
(6, 1500, 'Expert', '{"custom_avatar": true}'),
(7, 2100, 'Master', '{"custom_theme": true}'),
(8, 2800, 'Grandmaster', '{"early_access": true}'),
(9, 3600, 'Legend', '{"mentor_badge": true}'),
(10, 4500, 'Champion', '{"hall_of_fame": true}');

-- 2. User XP and Stats
CREATE TABLE user_xp (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    total_xp INT DEFAULT 0,
    level INT DEFAULT 1 REFERENCES xp_levels(level),
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_activity_date DATE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3. Badges Definitions
CREATE TABLE badges (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url TEXT,
    category VARCHAR(30), -- 'academic', 'engagement', 'milestone', 'special'
    criteria JSONB NOT NULL,
    xp_reward INT DEFAULT 0,
    rarity VARCHAR(20), -- 'common', 'uncommon', 'rare', 'epic', 'legendary'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- 4. User Earned Badges
CREATE TABLE user_badges (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMP DEFAULT NOW(),
    progress INT DEFAULT 100, -- For partial progress badges
    notified BOOLEAN DEFAULT false,
    UNIQUE(user_id, badge_id)
);

-- 5. XP Transaction Log
CREATE TABLE xp_events (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    xp_amount INT NOT NULL,
    source_type VARCHAR(50),
    source_id UUID,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for Leaderboards (filter by tenant + order by xp)
CREATE INDEX idx_user_xp_tenant_total ON user_xp(tenant_id, total_xp DESC);
