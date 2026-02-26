/// Data model for dashboard statistics.
class DashboardStats {
  final double totalRevenue;
  final int totalOrders;
  final int totalProducts;
  final int pendingOrders;
  final double revenueTrend;
  final double ordersTrend;
  final List<RevenueDataPoint> revenueChart;
  final List<RecentOrder> recentOrders;

  const DashboardStats({
    required this.totalRevenue,
    required this.totalOrders,
    required this.totalProducts,
    required this.pendingOrders,
    required this.revenueTrend,
    required this.ordersTrend,
    required this.revenueChart,
    required this.recentOrders,
  });
}

/// A single data point in the revenue chart.
class RevenueDataPoint {
  final DateTime date;
  final double amount;

  const RevenueDataPoint({required this.date, required this.amount});
}

/// A recent order shown on the dashboard.
class RecentOrder {
  final String id;
  final String customerName;
  final double total;
  final String status;
  final DateTime createdAt;
  final int itemCount;

  const RecentOrder({
    required this.id,
    required this.customerName,
    required this.total,
    required this.status,
    required this.createdAt,
    required this.itemCount,
  });
}

/// Repository for fetching seller dashboard data.
class DashboardRepository {
  /// Fetches all dashboard statistics including revenue, orders, products,
  /// revenue chart data for the last 30 days, and recent orders.
  Future<DashboardStats> getDashboardStats() async {
    // TODO: Replace with actual API call
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    final revenueChart = List.generate(
      30,
      (i) => RevenueDataPoint(
        date: now.subtract(Duration(days: 29 - i)),
        amount: 200.0 + (i * 15.0) + (i % 5 * 30.0),
      ),
    );

    final recentOrders = [
      RecentOrder(
        id: 'ORD-1001',
        customerName: 'Alice Johnson',
        total: 129.99,
        status: 'pending',
        createdAt: now.subtract(const Duration(hours: 2)),
        itemCount: 3,
      ),
      RecentOrder(
        id: 'ORD-1002',
        customerName: 'Bob Smith',
        total: 79.50,
        status: 'processing',
        createdAt: now.subtract(const Duration(hours: 5)),
        itemCount: 1,
      ),
      RecentOrder(
        id: 'ORD-1003',
        customerName: 'Carol Davis',
        total: 249.00,
        status: 'shipped',
        createdAt: now.subtract(const Duration(days: 1)),
        itemCount: 4,
      ),
      RecentOrder(
        id: 'ORD-1004',
        customerName: 'Dan Wilson',
        total: 55.00,
        status: 'delivered',
        createdAt: now.subtract(const Duration(days: 2)),
        itemCount: 2,
      ),
    ];

    return DashboardStats(
      totalRevenue: 15420.50,
      totalOrders: 187,
      totalProducts: 45,
      pendingOrders: 12,
      revenueTrend: 12.5,
      ordersTrend: 8.3,
      revenueChart: revenueChart,
      recentOrders: recentOrders,
    );
  }
}
