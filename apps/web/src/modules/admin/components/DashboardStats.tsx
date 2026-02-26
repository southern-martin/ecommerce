import { Card, CardContent } from '@/shared/components/ui/card';
import { DollarSign, Users, ShoppingCart, Package } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import type { AdminDashboardStats as Stats } from '../services/admin-user.api';

interface DashboardStatsProps {
  stats: Stats;
}

export function DashboardStats({ stats }: DashboardStatsProps) {
  const cards = [
    { label: 'Total Revenue', value: formatPrice(stats.total_revenue), icon: DollarSign },
    { label: 'Total Users', value: stats.total_users.toLocaleString(), icon: Users },
    { label: 'Total Orders', value: stats.total_orders.toLocaleString(), icon: ShoppingCart },
    { label: 'Total Products', value: stats.total_products.toLocaleString(), icon: Package },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card) => (
        <Card key={card.label}>
          <CardContent className="flex items-center gap-4 p-6">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
              <card.icon className="h-6 w-6 text-primary" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{card.label}</p>
              <p className="text-2xl font-bold">{card.value}</p>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
