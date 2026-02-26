import { useState } from 'react';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { MessageSquare } from 'lucide-react';
import { ConversationList } from '../components/ConversationList';
import { ChatWindow } from '../components/ChatWindow';
import { useConversations } from '../hooks/useConversations';

export default function MessagesPage() {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const { data: conversations, isLoading } = useConversations();

  const selectedConversation = conversations?.find((c) => c.id === selectedId);

  if (isLoading) {
    return (
      <div className="flex h-[600px] gap-0 rounded-lg border">
        <div className="w-80 border-r p-4 space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={i} className="h-16 w-full" />
          ))}
        </div>
        <div className="flex-1">
          <Skeleton className="h-full w-full" />
        </div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Messages</h1>
      <div className="flex h-[600px] rounded-lg border">
        <div className="w-80 overflow-y-auto border-r p-2">
          {conversations && conversations.length > 0 ? (
            <ConversationList
              conversations={conversations}
              selectedId={selectedId ?? undefined}
              onSelect={setSelectedId}
            />
          ) : (
            <div className="flex flex-col items-center py-8">
              <MessageSquare className="h-8 w-8 text-muted-foreground/50" />
              <p className="mt-2 text-sm text-muted-foreground">No conversations</p>
            </div>
          )}
        </div>

        <div className="flex-1">
          {selectedConversation ? (
            <ChatWindow
              conversationId={selectedConversation.id}
              participantName={selectedConversation.participant.name}
            />
          ) : (
            <div className="flex h-full items-center justify-center">
              <p className="text-muted-foreground">Select a conversation to start messaging</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
