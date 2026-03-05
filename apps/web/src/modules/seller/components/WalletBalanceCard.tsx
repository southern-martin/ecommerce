import { Link } from 'react-router-dom';
import { Card, CardContent } from '@/shared/components/ui/card';
import { Button } from '@/shared/components/ui/button';
import { Wallet, ArrowRight } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';

interface WalletBalanceCardProps {
  availableBalance: number;
  pendingBalance: number;
  currency?: string;
  showLink?: boolean;
}

export function WalletBalanceCard({
  availableBalance,
  pendingBalance,
  currency = 'USD',
  showLink,
}: WalletBalanceCardProps) {
  return (
    <Card>
      <CardContent className="flex items-center gap-6 pt-6">
        <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-primary/10">
          <Wallet className="h-6 w-6 text-primary" />
        </div>
        <div className="flex flex-1 items-center gap-8">
          <div>
            <p className="text-sm text-muted-foreground">Available Balance</p>
            <p className="text-2xl font-bold">{formatPrice(availableBalance, currency)}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Pending Balance</p>
            <p className="text-2xl font-bold text-muted-foreground">
              {formatPrice(pendingBalance, currency)}
            </p>
          </div>
        </div>
        {showLink && (
          <Button variant="outline" size="sm" asChild>
            <Link to="/seller/wallet">
              View Wallet
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
