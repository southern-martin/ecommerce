import { useQuery } from '@tanstack/react-query';
import apiClient from '@/shared/lib/api-client';
import { Link } from 'react-router-dom';
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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">My Returns</h1>
        <Button asChild><Link to="/account/returns/new">Request Return</Link></Button>
      </div>
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
    </div>
  );
}
