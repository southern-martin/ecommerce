import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';

const data = [
  { name: 'Pending', value: 35, color: '#f59e0b' },
  { name: 'Processing', value: 22, color: '#3b82f6' },
  { name: 'Shipped', value: 48, color: '#8b5cf6' },
  { name: 'Delivered', value: 120, color: '#22c55e' },
  { name: 'Cancelled', value: 8, color: '#ef4444' },
  { name: 'Returned', value: 5, color: '#6b7280' },
];

const totalOrders = data.reduce((sum, entry) => sum + entry.value, 0);

export default function OrderStatusChart() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Order Status Breakdown</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={320}>
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="45%"
              innerRadius={70}
              outerRadius={100}
              paddingAngle={3}
              dataKey="value"
              label={({ name, percent }) =>
                `${name} ${(percent * 100).toFixed(0)}%`
              }
            >
              {data.map((entry) => (
                <Cell key={entry.name} fill={entry.color} />
              ))}
            </Pie>
            <Tooltip
              formatter={(value: number) => [`${value} orders`, 'Count']}
            />
            <Legend verticalAlign="bottom" height={36} />
            {/* Center label */}
            <text
              x="50%"
              y="45%"
              textAnchor="middle"
              dominantBaseline="central"
              className="fill-foreground"
            >
              <tspan x="50%" dy="-0.6em" fontSize="24" fontWeight="bold">
                {totalOrders}
              </tspan>
              <tspan x="50%" dy="1.4em" fontSize="12" fill="#6b7280">
                Total Orders
              </tspan>
            </text>
          </PieChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
