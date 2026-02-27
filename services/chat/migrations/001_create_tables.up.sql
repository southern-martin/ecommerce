CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Conversations
CREATE TABLE conversations (
    id UUID PRIMARY KEY,
    type VARCHAR(20) NOT NULL DEFAULT 'buyer_seller',
    participant_ids TEXT[],
    buyer_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    order_id UUID,
    subject VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    last_message_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_conversations_buyer_id ON conversations(buyer_id);
CREATE INDEX idx_conversations_seller_id ON conversations(seller_id);
CREATE INDEX idx_conversations_order_id ON conversations(order_id);
CREATE INDEX idx_conversations_status ON conversations(status);

-- Messages
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    sender_role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    attachments TEXT[],
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_messages_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);

-- Conversation participants
CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL,
    joined_at TIMESTAMPTZ,
    last_read_at TIMESTAMPTZ,
    CONSTRAINT fk_conversation_participants_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT uq_conversation_participants UNIQUE (conversation_id, user_id)
);

CREATE INDEX idx_conversation_participants_user_id ON conversation_participants(user_id);
