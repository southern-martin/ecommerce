import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

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

  factory DashboardStats.fromJson(Map<String, dynamic> json) {
    return DashboardStats(
      totalRevenue: (json['total_revenue'] as num).toDouble(),
      totalOrders: json['total_orders'] as int,
      totalProducts: json['total_products'] as int,
      pendingOrders: json['pending_orders'] as int,
      revenueTrend: (json['revenue_trend'] as num).toDouble(),
      ordersTrend: (json['orders_trend'] as num).toDouble(),
      revenueChart: (json['revenue_chart'] as List<dynamic>)
          .map((e) => RevenueDataPoint.fromJson(e as Map<String, dynamic>))
          .toList(),
      recentOrders: (json['recent_orders'] as List<dynamic>)
          .map((e) => RecentOrder.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}

/// A single data point in the revenue chart.
class RevenueDataPoint {
  final DateTime date;
  final double amount;

  const RevenueDataPoint({required this.date, required this.amount});

  factory RevenueDataPoint.fromJson(Map<String, dynamic> json) {
    return RevenueDataPoint(
      date: DateTime.parse(json['date'] as String),
      amount: (json['amount'] as num).toDouble(),
    );
  }
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

  factory RecentOrder.fromJson(Map<String, dynamic> json) {
    return RecentOrder(
      id: json['id'] as String,
      customerName: json['customer_name'] as String,
      total: (json['total'] as num).toDouble(),
      status: json['status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      itemCount: json['item_count'] as int,
    );
  }
}

/// Repository for fetching seller dashboard data.
class DashboardRepository {
  final ApiClient _apiClient;

  DashboardRepository({required ApiClient apiClient})
      : _apiClient = apiClient;

  /// Fetches all dashboard statistics including revenue, orders, products,
  /// revenue chart data, and recent orders.
  Future<DashboardStats> getDashboardStats() async {
    final response = await _apiClient.get(ApiEndpoints.sellerDashboard);
    return DashboardStats.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches revenue data points for the given [period] (e.g. '7d', '30d', '90d').
  Future<List<RevenueDataPoint>> getRevenueData(String period) async {
    final response = await _apiClient.get(
      '${ApiEndpoints.sellerDashboard}/revenue',
      queryParameters: {'period': period},
    );
    final data = response.data as List<dynamic>;
    return data
        .map((e) => RevenueDataPoint.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// Fetches the most recent orders for the dashboard summary.
  Future<List<RecentOrder>> getRecentOrders({int limit = 10}) async {
    final response = await _apiClient.get(
      '${ApiEndpoints.sellerDashboard}/recent-orders',
      queryParameters: {'limit': limit},
    );
    final data = response.data as List<dynamic>;
    return data
        .map((e) => RecentOrder.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
