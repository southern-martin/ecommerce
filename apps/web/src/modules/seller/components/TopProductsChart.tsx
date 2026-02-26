import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from 'recharts';

const topProductsData = [
  { name: 'Wireless Headphones', revenue: 4500 },
  { name: 'Phone Case', revenue: 3200 },
  { name: 'USB-C Cable', revenue: 2800 },
  { name: 'Screen Protector', revenue: 2100 },
  { name: 'Bluetooth Speaker', revenue: 1800 },
];

export function TopProductsChart() {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={topProductsData} layout="vertical">
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis type="number" />
        <YAxis dataKey="name" type="category" width={150} />
        <Tooltip
          formatter={(value: number) => [`$${value.toFixed(2)}`, 'Revenue']}
        />
        <Bar dataKey="revenue" fill="#3b82f6" radius={[0, 4, 4, 0]} />
      </BarChart>
    </ResponsiveContainer>
  );
}
