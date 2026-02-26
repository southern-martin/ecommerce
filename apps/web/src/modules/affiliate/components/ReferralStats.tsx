import { Card, CardContent } from '@/shared/components/ui/card';
import { MousePointer, Target, DollarSign, Percent } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import type { AffiliateStats } from '../services/affiliate.api';

interface ReferralStatsProps {
  stats: AffiliateStats;
}

export function ReferralStats({ stats }: ReferralStatsProps) {
  const cards = [
    { label: 'Total Clicks', value: stats.total_clicks.toLocaleString(), icon: MousePointer },
    { label: 'Conversions', value: stats.total_conversions.toLocaleString(), icon: Target },
    { label: 'Total Earnings', value: formatPrice(stats.total_earnings), icon: DollarSign },
    { label: 'Conversion Rate', value: `${stats.conversion_rate.toFixed(1)}%`, icon: Percent },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card) => (
        <Card key={card.label}>
          <CardContent className="flex items-center gap-4 p-6">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
              <card.icon className="h-5 w-5 text-primary" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground">{card.label}</p>
              <p className="text-xl font-bold">{card.value}</p>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
