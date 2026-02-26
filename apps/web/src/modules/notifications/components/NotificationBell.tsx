import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Bell } from 'lucide-react';
import { useUnreadCount } from '../hooks/useNotifications';

export function NotificationBell() {
  const { data: unreadCount } = useUnreadCount();

  return (
    <Button asChild variant="ghost" size="icon" className="relative">
      <Link to="/notifications">
        <Bell className="h-5 w-5" />
        {unreadCount && unreadCount > 0 && (
          <span className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-destructive text-xs text-destructive-foreground">
            {unreadCount > 99 ? '99+' : unreadCount}
          </span>
        )}
      </Link>
    </Button>
  );
}
