import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

const orderStatusData = [
  { name: 'Pending', value: 12 },
  { name: 'Processing', value: 25 },
  { name: 'Shipped', value: 18 },
  { name: 'Delivered', value: 42 },
  { name: 'Cancelled', value: 3 },
];

const COLORS = ['#f59e0b', '#6366f1', '#8b5cf6', '#22c55e', '#ef4444'];

export function OrderStatusChart() {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={orderStatusData}
          cx="50%"
          cy="50%"
          outerRadius={100}
          dataKey="value"
          label={({ name, percent }) =>
            `${name} ${(percent * 100).toFixed(0)}%`
          }
        >
          {orderStatusData.map((_, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
}
