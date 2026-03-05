import { useQuery } from '@tanstack/react-query';
import { chatApi } from '../services/chat.api';

export function useConversations() {
  return useQuery({
    queryKey: ['conversations'],
    queryFn: async () => {
      const result = await chatApi.getConversations();
      return result.conversations ?? [];
    },
    refetchInterval: 15000,
  });
}
