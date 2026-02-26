import { Avatar, AvatarFallback, AvatarImage } from '@/shared/components/ui/avatar';
import { Badge } from '@/shared/components/ui/badge';
import { cn } from '@/shared/lib/utils';
import { formatDateTime } from '@/shared/lib/utils';
import type { Conversation } from '../services/chat.api';

interface ConversationListProps {
  conversations: Conversation[];
  selectedId?: string;
  onSelect: (id: string) => void;
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
            <AvatarImage src={conv.participant.avatar_url} />
            <AvatarFallback>{conv.participant.name[0]}</AvatarFallback>
          </Avatar>
          <div className="flex-1 overflow-hidden">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">{conv.participant.name}</span>
              {conv.last_message && (
                <span className="text-xs text-muted-foreground">
                  {formatDateTime(conv.last_message.created_at)}
                </span>
              )}
            </div>
            {conv.last_message && (
              <p className="truncate text-xs text-muted-foreground">
                {conv.last_message.content}
              </p>
            )}
          </div>
          {conv.unread_count > 0 && (
            <Badge className="h-5 min-w-[20px] justify-center rounded-full px-1.5 text-xs">
              {conv.unread_count}
            </Badge>
          )}
        </button>
      ))}
    </div>
  );
}
