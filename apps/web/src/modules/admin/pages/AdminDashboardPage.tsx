import { Skeleton } from '@/shared/components/ui/skeleton';
import { DashboardStats } from '../components/DashboardStats';
import { useAdminDashboard } from '../hooks/useAdminDashboard';

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

      <div className="rounded-lg border p-6">
        <h2 className="mb-4 text-lg font-semibold">Activity Overview</h2>
        <div className="flex h-64 items-center justify-center text-muted-foreground">
          Charts and analytics will be rendered here.
        </div>
      </div>
    </div>
  );
}
