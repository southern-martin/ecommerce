import { useState } from 'react';
import { MessageCircle, MessageSquare } from 'lucide-react';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ConversationList } from '../components/ConversationList';
import { ChatWindow } from '../components/ChatWindow';
import { useConversations } from '../hooks/useConversations';
import type { Conversation } from '../services/chat.api';

function getConversationTitle(conv: Conversation): string {
  return conv.subject || `Conversation #${conv.id.slice(0, 8)}`;
}

export default function ChatPage() {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const { data: conversations, isLoading } = useConversations();

  const selectedConversation = conversations?.find((c) => c.id === selectedId);

  return (
    <PageLayout
      title="Messages"
      icon={MessageCircle}
      breadcrumbs={[
        { label: 'Account', href: '/account/profile' },
        { label: 'Messages' },
      ]}
    >
      {isLoading ? (
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
      ) : (
        <div className="flex h-[calc(100vh-16rem)] rounded-lg border">
          <div className="w-80 overflow-y-auto border-r p-2">
            {conversations && conversations.length > 0 ? (
              <ConversationList
                conversations={conversations}
                selectedId={selectedId ?? undefined}
                onSelect={setSelectedId}
              />
            ) : (
              <div className="flex flex-col items-center justify-center h-full py-8">
                <MessageSquare className="h-8 w-8 text-muted-foreground/50" />
                <p className="mt-2 text-sm text-muted-foreground">No conversations yet</p>
              </div>
            )}
          </div>

          <div className="flex-1">
            {selectedConversation ? (
              <ChatWindow
                conversationId={selectedConversation.id}
                title={getConversationTitle(selectedConversation)}
              />
            ) : (
              <div className="flex h-full flex-col items-center justify-center text-muted-foreground">
                <MessageCircle className="h-10 w-10 mb-3 opacity-50" />
                <p>Select a conversation to start messaging</p>
              </div>
            )}
          </div>
        </div>
      )}
    </PageLayout>
  );
}
