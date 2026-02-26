import { UserPlus, ShoppingCart, ShieldCheck, MessageSquare } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';

interface ActivityItem {
  id: string;
  icon: 'user' | 'order' | 'seller' | 'review';
  title: string;
  description: string;
  timestamp: string;
}

const iconMap = {
  user: UserPlus,
  order: ShoppingCart,
  seller: ShieldCheck,
  review: MessageSquare,
};

const iconColorMap = {
  user: 'text-blue-500',
  order: 'text-green-500',
  seller: 'text-purple-500',
  review: 'text-orange-500',
};

const mockActivities: ActivityItem[] = [
  {
    id: '1',
    icon: 'user',
    title: 'New user registered',
    description: 'john.doe@example.com created an account',
    timestamp: '2 minutes ago',
  },
  {
    id: '2',
    icon: 'order',
    title: 'Order #1234 placed',
    description: 'Order total: $149.99 - 3 items',
    timestamp: '15 minutes ago',
  },
  {
    id: '3',
    icon: 'seller',
    title: 'Seller approved',
    description: 'TechStore has been approved as a seller',
    timestamp: '1 hour ago',
  },
  {
    id: '4',
    icon: 'review',
    title: 'Review submitted',
    description: '5-star review on "Wireless Headphones"',
    timestamp: '2 hours ago',
  },
  {
    id: '5',
    icon: 'order',
    title: 'Order #1230 shipped',
    description: 'Tracking number: 1Z999AA10123456784',
    timestamp: '3 hours ago',
  },
  {
    id: '6',
    icon: 'user',
    title: 'New user registered',
    description: 'jane.smith@example.com created an account',
    timestamp: '4 hours ago',
  },
];

export default function RecentActivityFeed() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="max-h-[400px] space-y-4 overflow-y-auto pr-2">
          {mockActivities.map((activity) => {
            const Icon = iconMap[activity.icon];
            const colorClass = iconColorMap[activity.icon];
            return (
              <div key={activity.id} className="flex items-start gap-3">
                <div className={`mt-0.5 rounded-full bg-muted p-2 ${colorClass}`}>
                  <Icon className="h-4 w-4" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium">{activity.title}</p>
                  <p className="text-sm text-muted-foreground truncate">
                    {activity.description}
                  </p>
                </div>
                <span className="text-xs text-muted-foreground whitespace-nowrap">
                  {activity.timestamp}
                </span>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
