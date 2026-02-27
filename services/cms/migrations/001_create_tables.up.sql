CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Banners
CREATE TABLE banners (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image_url TEXT NOT NULL,
    link_url TEXT,
    position VARCHAR(50),
    sort_order INTEGER DEFAULT 0,
    target_audience VARCHAR(100),
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_banners_position ON banners(position);
CREATE INDEX idx_banners_ends_at ON banners(ends_at);
CREATE INDEX idx_banners_is_active ON banners(is_active);

-- Pages
CREATE TABLE pages (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content_html TEXT,
    meta_title VARCHAR(255),
    meta_description TEXT,
    status VARCHAR(20) DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pages_status ON pages(status);

-- Content schedules
CREATE TABLE content_schedules (
    id UUID PRIMARY KEY,
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    executed BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_content_schedules_scheduled_at ON content_schedules(scheduled_at);
CREATE INDEX idx_content_schedules_executed ON content_schedules(executed);
