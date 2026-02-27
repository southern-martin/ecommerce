CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Media files
CREATE TABLE media_files (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    owner_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    original_name VARCHAR(500) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT DEFAULT 0,
    url TEXT,
    thumbnail_url TEXT,
    width INTEGER DEFAULT 0,
    height INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_media_files_owner_id ON media_files(owner_id);
CREATE INDEX idx_media_files_owner_type ON media_files(owner_type);
