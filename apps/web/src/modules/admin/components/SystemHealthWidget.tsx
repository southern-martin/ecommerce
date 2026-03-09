import { Activity } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';

type Status = 'healthy' | 'degraded' | 'down';

interface Service {
  name: string;
  status: Status;
}

const services: Service[] = [
  { name: 'API Gateway', status: 'healthy' },
  { name: 'Auth', status: 'healthy' },
  { name: 'Products', status: 'healthy' },
  { name: 'Orders', status: 'degraded' },
  { name: 'Payments', status: 'healthy' },
  { name: 'Search', status: 'healthy' },
  { name: 'Database', status: 'healthy' },
  { name: 'Cache', status: 'down' },
];

const statusColors: Record<Status, string> = {
  healthy: 'bg-green-500',
  degraded: 'bg-yellow-500',
  down: 'bg-red-500',
};

const statusLabels: Record<Status, string> = {
  healthy: 'Healthy',
  degraded: 'Degraded',
  down: 'Down',
};

export default function SystemHealthWidget() {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            System Health
          </CardTitle>
          <span className="text-sm text-muted-foreground">
            Uptime: <span className="font-semibold text-foreground">99.9%</span>
          </span>
        </div>
      </CardHeader>
      <CardContent>
        <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-4">
          {services.map((service) => (
            <div
              key={service.name}
              className="flex items-center gap-3 rounded-lg border p-3"
            >
              <span
                className={`h-2.5 w-2.5 shrink-0 rounded-full ${statusColors[service.status]}`}
              />
              <div className="min-w-0">
                <p className="truncate text-sm font-medium">{service.name}</p>
                <p className="text-xs text-muted-foreground">
                  {statusLabels[service.status]}
                </p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
