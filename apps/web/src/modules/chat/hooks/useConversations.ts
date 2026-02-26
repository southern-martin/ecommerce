import { useQuery } from '@tanstack/react-query';
import { chatApi } from '../services/chat.api';

export function useConversations() {
  return useQuery({
    queryKey: ['conversations'],
    queryFn: () => chatApi.getConversations(),
    refetchInterval: 15000,
  });
}
