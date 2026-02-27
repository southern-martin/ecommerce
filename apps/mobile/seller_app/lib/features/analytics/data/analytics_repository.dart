import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Revenue summary for a given time period.
class RevenueSummary {
  final double totalRevenue;
  final double averageOrderValue;
  final int totalOrders;
  final double revenueTrend;

  const RevenueSummary({
    required this.totalRevenue,
    required this.averageOrderValue,
    required this.totalOrders,
    required this.revenueTrend,
  });

  factory RevenueSummary.fromJson(Map<String, dynamic> json) {
    return RevenueSummary(
      totalRevenue: (json['totalRevenue'] as num).toDouble(),
      averageOrderValue: (json['averageOrderValue'] as num).toDouble(),
      totalOrders: json['totalOrders'] as int,
      revenueTrend: (json['revenueTrend'] as num).toDouble(),
    );
  }
}

/// A single data point for the sales chart.
class SalesDataPoint {
  final String label;
  final double revenue;
  final int orders;

  const SalesDataPoint({
    required this.label,
    required this.revenue,
    required this.orders,
  });

  factory SalesDataPoint.fromJson(Map<String, dynamic> json) {
    return SalesDataPoint(
      label: json['label'] as String,
      revenue: (json['revenue'] as num).toDouble(),
      orders: json['orders'] as int,
    );
  }
}

/// Represents a top-selling product.
class TopProduct {
  final String id;
  final String name;
  final int unitsSold;
  final double revenue;

  const TopProduct({
    required this.id,
    required this.name,
    required this.unitsSold,
    required this.revenue,
  });

  factory TopProduct.fromJson(Map<String, dynamic> json) {
    return TopProduct(
      id: json['id'] as String,
      name: json['name'] as String,
      unitsSold: json['unitsSold'] as int,
      revenue: (json['revenue'] as num).toDouble(),
    );
  }
}

/// Order statistics summary.
class OrderStats {
  final int totalOrders;
  final int pendingOrders;
  final int completedOrders;
  final int cancelledOrders;
  final double fulfillmentRate;

  const OrderStats({
    required this.totalOrders,
    required this.pendingOrders,
    required this.completedOrders,
    required this.cancelledOrders,
    required this.fulfillmentRate,
  });

  factory OrderStats.fromJson(Map<String, dynamic> json) {
    return OrderStats(
      totalOrders: json['totalOrders'] as int,
      pendingOrders: json['pendingOrders'] as int,
      completedOrders: json['completedOrders'] as int,
      cancelledOrders: json['cancelledOrders'] as int,
      fulfillmentRate: (json['fulfillmentRate'] as num).toDouble(),
    );
  }
}

/// Repository for fetching seller analytics data.
class AnalyticsRepository {
  final ApiClient _apiClient;

  AnalyticsRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Fetches the revenue summary for the given [period] (e.g. '7d', '30d', '90d', '1y').
  Future<RevenueSummary> getRevenueSummary(String period) async {
    final response = await _apiClient
        .get('${ApiEndpoints.sellerAnalytics}/revenue?period=$period');
    return RevenueSummary.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches the top-selling products, limited to [limit] results.
  Future<List<TopProduct>> getTopProducts(int limit) async {
    final response = await _apiClient
        .get('${ApiEndpoints.sellerAnalytics}/top-products?limit=$limit');
    final list = response.data as List<dynamic>;
    return list
        .map((e) => TopProduct.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// Fetches sales chart data points for the given [period].
  Future<List<SalesDataPoint>> getSalesChart(String period) async {
    final response = await _apiClient
        .get('${ApiEndpoints.sellerAnalytics}/sales-chart?period=$period');
    final list = response.data as List<dynamic>;
    return list
        .map((e) => SalesDataPoint.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// Fetches overall order statistics.
  Future<OrderStats> getOrderStats() async {
    final response =
        await _apiClient.get('${ApiEndpoints.sellerAnalytics}/orders');
    return OrderStats.fromJson(response.data as Map<String, dynamic>);
  }
}
