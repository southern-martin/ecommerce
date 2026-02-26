CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE auth_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'buyer',
    oauth_provider VARCHAR(50),
    oauth_provider_id VARCHAR(255),
    refresh_token TEXT,
    reset_token VARCHAR(255),
    reset_token_exp TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_users_email ON auth_users(email);
CREATE INDEX idx_auth_users_oauth ON auth_users(oauth_provider, oauth_provider_id) WHERE oauth_provider IS NOT NULL;
