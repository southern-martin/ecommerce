import { Avatar, AvatarFallback } from '@/shared/components/ui/avatar';
import { cn, formatDateTime } from '@/shared/lib/utils';
import type { Conversation } from '../services/chat.api';

interface ConversationListProps {
  conversations: Conversation[];
  selectedId?: string;
  onSelect: (id: string) => void;
}

function getDisplayName(conv: Conversation): string {
  return conv.subject || `Conversation #${conv.id.slice(0, 8)}`;
}

function getInitial(conv: Conversation): string {
  if (conv.subject) return conv.subject[0].toUpperCase();
  return conv.type === 'support' ? 'S' : 'C';
}

export function ConversationList({
  conversations,
  selectedId,
  onSelect,
}: ConversationListProps) {
  return (
    <div className="space-y-1">
      {conversations.map((conv) => (
        <button
          key={conv.id}
          onClick={() => onSelect(conv.id)}
          className={cn(
            'flex w-full items-center gap-3 rounded-md px-3 py-3 text-left transition-colors',
            selectedId === conv.id ? 'bg-muted' : 'hover:bg-muted/50'
          )}
        >
          <Avatar className="h-10 w-10">
            <AvatarFallback>{getInitial(conv)}</AvatarFallback>
          </Avatar>
          <div className="flex-1 overflow-hidden">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium truncate">
                {getDisplayName(conv)}
              </span>
              <span className="text-xs text-muted-foreground shrink-0 ml-2">
                {conv.last_message_at
                  ? formatDateTime(conv.last_message_at)
                  : formatDateTime(conv.updated_at)}
              </span>
            </div>
            <p className="truncate text-xs text-muted-foreground">
              {conv.status === 'active' ? 'Active' : conv.status}
            </p>
          </div>
        </button>
      ))}
    </div>
  );
}
