import { RotateCcw } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import apiClient from '@/shared/lib/api-client';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { Button } from '@/shared/components/ui/button';
import { Card } from '@/shared/components/ui/card';
import { Badge } from '@/shared/components/ui/badge';

export default function ReturnListPage() {
  const { data: returns = [], isLoading } = useQuery({
    queryKey: ['returns'],
    queryFn: async () => {
      const res = await apiClient.get('/returns');
      const d = res.data;
      return Array.isArray(d) ? d : (d as any).data || [];
    },
  });

  if (isLoading) return <div className="p-6">Loading returns...</div>;

  const requestReturnButton = (
    <Button asChild>
      <Link to="/account/returns/new">Request Return</Link>
    </Button>
  );

  return (
    <PageLayout
      title="My Returns"
      icon={RotateCcw}
      breadcrumbs={[
        { label: 'Account', href: '/account/profile' },
        { label: 'Returns' },
      ]}
      actions={requestReturnButton}
    >
      {returns.length === 0 ? (
        <p className="text-muted-foreground text-center py-8">No returns yet.</p>
      ) : (
        <div className="space-y-4">
          {returns.map((ret: any) => (
            <Card key={ret.id} className="p-4">
              <div className="flex justify-between items-center">
                <div>
                  <p className="font-medium">Return #{ret.id?.slice(0, 8)}</p>
                  <p className="text-sm text-muted-foreground">{ret.reason}</p>
                </div>
                <Badge variant="outline">{ret.status}</Badge>
              </div>
            </Card>
          ))}
        </div>
      )}
    </PageLayout>
  );
}
