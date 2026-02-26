import { Card, CardContent } from '@/shared/components/ui/card';
import { Coins } from 'lucide-react';

interface PointsBalanceProps {
  balance: number;
  tier: string;
}

export function PointsBalance({ balance, tier }: PointsBalanceProps) {
  return (
    <Card>
      <CardContent className="flex items-center gap-4 p-6">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-yellow-100">
          <Coins className="h-6 w-6 text-yellow-600" />
        </div>
        <div>
          <p className="text-sm text-muted-foreground">Points Balance</p>
          <p className="text-2xl font-bold">{balance.toLocaleString()}</p>
        </div>
        <div className="ml-auto">
          <span className="rounded-full bg-primary/10 px-3 py-1 text-sm font-medium capitalize text-primary">
            {tier}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}
