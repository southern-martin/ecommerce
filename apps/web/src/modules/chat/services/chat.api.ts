import apiClient from '@/shared/lib/api-client';

export interface Conversation {
  id: string;
  type: string;
  buyer_id: string;
  seller_id: string;
  order_id?: string;
  subject?: string;
  status: string;
  last_message_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: string;
  conversation_id: string;
  sender_id: string;
  sender_role: string;
  content: string;
  message_type: string;
  attachments?: string[];
  is_read: boolean;
  read_at?: string;
  created_at: string;
}

export interface ConversationWithMessages {
  conversation: Conversation;
  messages: Message[];
}

export interface CreateConversationInput {
  buyer_id: string;
  seller_id: string;
  order_id?: string;
  subject?: string;
  type?: string;
}

export const chatApi = {
  getConversations: async (
    params?: { status?: string; page?: number; page_size?: number }
  ): Promise<{ conversations: Conversation[]; total: number; page: number; page_size: number }> => {
    const response = await apiClient.get('/chat/conversations', { params });
    return response.data;
  },

  getConversation: async (conversationId: string): Promise<ConversationWithMessages> => {
    const response = await apiClient.get(`/chat/conversations/${conversationId}`);
    return response.data;
  },

  getMessages: async (
    conversationId: string,
    params: { page: number; page_size: number }
  ): Promise<{ messages: Message[]; total: number; page: number; page_size: number }> => {
    const response = await apiClient.get(
      `/chat/conversations/${conversationId}/messages`,
      { params }
    );
    return response.data;
  },

  sendMessage: async (conversationId: string, content: string): Promise<Message> => {
    const response = await apiClient.post(
      `/chat/conversations/${conversationId}/messages`,
      { content }
    );
    return response.data.message ?? response.data.data ?? response.data;
  },

  createConversation: async (data: CreateConversationInput): Promise<Conversation> => {
    const response = await apiClient.post('/chat/conversations', data);
    return response.data.conversation ?? response.data.data ?? response.data;
  },

  markAsRead: async (conversationId: string): Promise<void> => {
    await apiClient.patch(`/chat/conversations/${conversationId}/messages/read`);
  },

  getUnreadCount: async (conversationId: string): Promise<number> => {
    const response = await apiClient.get(`/chat/conversations/${conversationId}/unread`);
    return response.data.unread_count ?? 0;
  },

  archiveConversation: async (conversationId: string): Promise<void> => {
    await apiClient.patch(`/chat/conversations/${conversationId}/archive`);
  },
};
