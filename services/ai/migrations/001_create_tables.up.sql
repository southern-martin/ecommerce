CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Embeddings
CREATE TABLE embeddings (
    id UUID PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    embedding_vector TEXT,
    model_version VARCHAR(100),
    dimensions INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_embeddings_entity UNIQUE (entity_type, entity_id)
);

-- Recommendations
CREATE TABLE recommendations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    reason VARCHAR(100),
    is_viewed BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_recommendations_user_id ON recommendations(user_id);

-- AI conversations
CREATE TABLE ai_conversations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(255),
    messages_json TEXT,
    model VARCHAR(50),
    token_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ai_conversations_user_id ON ai_conversations(user_id);

-- Generated contents
CREATE TABLE generated_contents (
    id UUID PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    content TEXT,
    model VARCHAR(50),
    prompt_tokens INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_generated_contents_entity UNIQUE (entity_type, entity_id)
);
