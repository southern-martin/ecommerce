import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Badge } from '@/shared/components/ui/badge';
import { formatDate } from '@/shared/lib/utils';
import { PointsBalance } from '../components/PointsBalance';
import { TierProgressBar } from '../components/TierProgressBar';
import { useMembership } from '../hooks/useMembership';
import { usePointsHistory } from '../hooks/usePointsBalance';

export default function LoyaltyDashboardPage() {
  const { data: membership, isLoading: membershipLoading } = useMembership();
  const { data: history, isLoading: historyLoading } = usePointsHistory();

  if (membershipLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-28 w-full" />
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!membership) return null;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Loyalty Program</h1>

      <PointsBalance balance={membership.points_balance} tier={membership.tier} />

      <Card>
        <CardContent className="p-6">
          <TierProgressBar
            currentTier={membership.tier}
            nextTier={membership.next_tier}
            progress={membership.tier_progress_percentage}
            pointsToNext={membership.points_to_next_tier}
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Points History</CardTitle>
        </CardHeader>
        <CardContent>
          {historyLoading ? (
            <Skeleton className="h-48 w-full" />
          ) : history && history.length > 0 ? (
            <div className="space-y-3">
              {history.map((txn) => (
                <div key={txn.id} className="flex items-center justify-between border-b pb-3 last:border-0">
                  <div>
                    <p className="text-sm font-medium">{txn.description}</p>
                    <p className="text-xs text-muted-foreground">{formatDate(txn.created_at)}</p>
                  </div>
                  <Badge variant={txn.type === 'earned' ? 'default' : txn.type === 'redeemed' ? 'secondary' : 'destructive'}>
                    {txn.type === 'earned' ? '+' : '-'}{txn.points} pts
                  </Badge>
                </div>
              ))}
            </div>
          ) : (
            <p className="py-4 text-center text-sm text-muted-foreground">No points history yet.</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
