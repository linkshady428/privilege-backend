-- Enable PostGIS for geohash / radius queries
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum types
CREATE TYPE user_tier AS ENUM ('free', 'privilege');
CREATE TYPE sex_type AS ENUM ('male', 'female', 'non_binary', 'genderqueer', 'prefer_not_to_say', 'other');
CREATE TYPE relationship_status AS ENUM ('single', 'in_relationship', 'married', 'divorced', 'open_relationship', 'complicated');
CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'rejected', 'void');
CREATE TYPE report_reason AS ENUM ('spam', 'inappropriate_content', 'feels_unsafe');

-- Users
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT,                        -- NULL when using OAuth-only login
    name          TEXT NOT NULL,
    birthdate     DATE NOT NULL,               -- immutable after signup
    sex           sex_type NOT NULL,
    bio           TEXT NOT NULL DEFAULT '',
    job           TEXT,
    height_cm     SMALLINT,
    weight_kg     SMALLINT,
    relationship_status relationship_status,
    lifestyle_tags TEXT[] NOT NULL DEFAULT '{}',
    tier          user_tier NOT NULL DEFAULT 'free',
    location      GEOGRAPHY(Point, 4326),      -- lat/lng stored as PostGIS point
    geohash       TEXT,                        -- for fast radius queries
    city          TEXT,
    agree_tos     BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at    TIMESTAMPTZ,                 -- soft-delete; hard-delete after 30 days
    last_active   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_tier ON users(tier) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_geohash ON users(geohash) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_last_active ON users(last_active DESC) WHERE deleted_at IS NULL;

-- Photos (up to 3 for free, 6 for privilege — enforced in application layer)
CREATE TABLE photos (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    position   SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_photos_user_id ON photos(user_id);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- Pass log (Privilege user swiped left on Free user — permanent, append-only)
CREATE TABLE passes (
    privilege_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    free_user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (privilege_user_id, free_user_id)
);

-- Invitations (Privilege swipes right → Free user receives)
CREATE TABLE invitations (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,   -- Privilege user
    recipient_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,   -- Free user
    status        invitation_status NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invitations_recipient ON invitations(recipient_id, status);
CREATE INDEX idx_invitations_sender    ON invitations(sender_id);

-- Matches (created when Free user accepts an invitation)
CREATE TABLE matches (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invitation_id UUID NOT NULL UNIQUE REFERENCES invitations(id),
    user_a_id     UUID NOT NULL REFERENCES users(id),  -- Privilege user
    user_b_id     UUID NOT NULL REFERENCES users(id),  -- Free user
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_matches_user_a ON matches(user_a_id);
CREATE INDEX idx_matches_user_b ON matches(user_b_id);

-- Messages
CREATE TABLE messages (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id   UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    sender_id  UUID NOT NULL REFERENCES users(id),
    body       TEXT NOT NULL,
    sent_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_match_id_sent_at ON messages(match_id, sent_at DESC);

-- Blocks
CREATE TABLE blocks (
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (blocker_id, blocked_id)
);

-- Reports
CREATE TABLE reports (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reported_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason      report_reason NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
