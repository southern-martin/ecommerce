import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from 'recharts';

const revenueData = [
  { date: 'Mon', revenue: 1200 },
  { date: 'Tue', revenue: 1800 },
  { date: 'Wed', revenue: 1400 },
  { date: 'Thu', revenue: 2200 },
  { date: 'Fri', revenue: 2800 },
  { date: 'Sat', revenue: 3200 },
  { date: 'Sun', revenue: 2600 },
];

export function RevenueChart() {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={revenueData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis />
        <Tooltip
          formatter={(value: number) => [`$${value.toFixed(2)}`, 'Revenue']}
        />
        <Line
          type="monotone"
          dataKey="revenue"
          stroke="#3b82f6"
          strokeWidth={2}
          dot={{ fill: '#3b82f6' }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
