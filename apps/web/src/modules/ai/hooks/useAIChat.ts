import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { aiApi } from '../services/ai.api';
import type { AIChatMessage, AIChatResponse } from '../services/ai.api';

export function useAIChat() {
  const [messages, setMessages] = useState<AIChatMessage[]>([]);
  const [lastResponse, setLastResponse] = useState<AIChatResponse | null>(null);

  const mutation = useMutation({
    mutationFn: (newMessages: AIChatMessage[]) => aiApi.chat(newMessages),
    onSuccess: (response) => {
      setMessages((prev) => [...prev, { role: 'assistant', content: response.message }]);
      setLastResponse(response);
    },
  });

  const sendMessage = (content: string) => {
    const userMessage: AIChatMessage = { role: 'user', content };
    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);
    mutation.mutate(updatedMessages);
  };

  const clearChat = () => {
    setMessages([]);
    setLastResponse(null);
  };

  return {
    messages,
    lastResponse,
    sendMessage,
    clearChat,
    isLoading: mutation.isPending,
  };
}
