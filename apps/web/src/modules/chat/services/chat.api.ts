import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';

export interface Conversation {
  id: string;
  participant: {
    id: string;
    name: string;
    avatar_url?: string;
  };
  last_message?: Message;
  unread_count: number;
  updated_at: string;
}

export interface Message {
  id: string;
  conversation_id: string;
  sender_id: string;
  content: string;
  is_read: boolean;
  created_at: string;
}

export const chatApi = {
  getConversations: async (): Promise<Conversation[]> => {
    const response = await apiClient.get<ApiResponse<Conversation[]>>('/chat/conversations');
    return response.data.data;
  },

  getMessages: async (
    conversationId: string,
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<Message>> => {
    const response = await apiClient.get<PaginatedResponse<Message>>(
      `/chat/conversations/${conversationId}/messages`,
      { params }
    );
    return response.data;
  },

  sendMessage: async (conversationId: string, content: string): Promise<Message> => {
    const response = await apiClient.post<ApiResponse<Message>>(
      `/chat/conversations/${conversationId}/messages`,
      { content }
    );
    return response.data.data;
  },

  createConversation: async (participantId: string): Promise<Conversation> => {
    const response = await apiClient.post<ApiResponse<Conversation>>('/chat/conversations', {
      participant_id: participantId,
    });
    return response.data.data;
  },
};
