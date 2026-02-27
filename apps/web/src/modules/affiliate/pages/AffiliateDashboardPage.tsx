import { Link2, DollarSign } from 'lucide-react';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ReferralStats } from '../components/ReferralStats';
import { ReferralLinkGenerator } from '../components/ReferralLinkGenerator';
import { useAffiliateStats, useRequestPayout } from '../hooks/useAffiliateStats';
import { formatPrice } from '@/shared/lib/utils';

export default function AffiliateDashboardPage() {
  const { data: stats, isLoading } = useAffiliateStats();
  const requestPayout = useRequestPayout();

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-4 md:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-24" />
          ))}
        </div>
        <Skeleton className="h-48" />
      </div>
    );
  }

  if (!stats) return null;

  const payoutButton = stats.pending_payout > 0 ? (
    <Button onClick={() => requestPayout.mutate()} disabled={requestPayout.isPending}>
      <DollarSign className="mr-2 h-4 w-4" />
      Request Payout ({formatPrice(stats.pending_payout)})
    </Button>
  ) : undefined;

  return (
    <PageLayout
      title="Affiliate Dashboard"
      icon={Link2}
      breadcrumbs={[
        { label: 'Account', href: '/account/profile' },
        { label: 'Affiliate' },
      ]}
      actions={payoutButton}
    >
      <div className="space-y-8">
        <ReferralStats stats={stats} />
        <ReferralLinkGenerator />
      </div>
    </PageLayout>
  );
}
