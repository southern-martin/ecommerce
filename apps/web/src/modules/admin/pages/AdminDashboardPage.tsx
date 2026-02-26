import { lazy, Suspense } from 'react';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { DashboardStats } from '../components/DashboardStats';
import { useAdminDashboard } from '../hooks/useAdminDashboard';

const RevenueOverviewChart = lazy(() => import('../components/RevenueOverviewChart'));
const RecentActivityFeed = lazy(() => import('../components/RecentActivityFeed'));

export default function AdminDashboardPage() {
  const { data: stats, isLoading } = useAdminDashboard();

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-28" />
          ))}
        </div>
        <Skeleton className="h-64" />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Admin Dashboard</h1>
      {stats && <DashboardStats stats={stats} />}

      <div className="grid gap-6 lg:grid-cols-2">
        <Suspense fallback={<Skeleton className="h-[380px]" />}>
          <RevenueOverviewChart />
        </Suspense>
        <Suspense fallback={<Skeleton className="h-[380px]" />}>
          <RecentActivityFeed />
        </Suspense>
      </div>
    </div>
  );
}
