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
}

/// Repository for fetching seller analytics data.
class AnalyticsRepository {
  /// Fetches the revenue summary for the given [period] (e.g. '7d', '30d', '90d', '1y').
  Future<RevenueSummary> getRevenueSummary(String period) async {
    // TODO: Replace with actual API call to /seller/analytics/revenue?period=
    await Future.delayed(const Duration(seconds: 1));

    final multiplier = _periodMultiplier(period);
    return RevenueSummary(
      totalRevenue: 12450.00 * multiplier,
      averageOrderValue: 67.50 * (1 + multiplier * 0.1),
      totalOrders: (185 * multiplier).round(),
      revenueTrend: 12.5,
    );
  }

  /// Fetches the top-selling products, limited to [limit] results.
  Future<List<TopProduct>> getTopProducts(int limit) async {
    // TODO: Replace with actual API call to /seller/analytics/top-products?limit=
    await Future.delayed(const Duration(milliseconds: 800));

    final products = [
      const TopProduct(
          id: 'prod_1',
          name: 'Wireless Bluetooth Headphones',
          unitsSold: 142,
          revenue: 11218.00),
      const TopProduct(
          id: 'prod_2',
          name: 'USB-C Hub Adapter',
          unitsSold: 98,
          revenue: 4802.00),
      const TopProduct(
          id: 'prod_3',
          name: 'Laptop Stand - Aluminum',
          unitsSold: 76,
          revenue: 3724.00),
      const TopProduct(
          id: 'prod_4',
          name: 'Mechanical Keyboard RGB',
          unitsSold: 64,
          revenue: 5120.00),
      const TopProduct(
          id: 'prod_5',
          name: 'Webcam 1080p HD',
          unitsSold: 53,
          revenue: 2597.00),
      const TopProduct(
          id: 'prod_6',
          name: 'Monitor Light Bar',
          unitsSold: 47,
          revenue: 2209.00),
      const TopProduct(
          id: 'prod_7',
          name: 'Desk Mat XL',
          unitsSold: 41,
          revenue: 1189.00),
      const TopProduct(
          id: 'prod_8',
          name: 'Cable Management Kit',
          unitsSold: 38,
          revenue: 722.00),
      const TopProduct(
          id: 'prod_9',
          name: 'Wireless Mouse Ergonomic',
          unitsSold: 35,
          revenue: 1715.00),
      const TopProduct(
          id: 'prod_10',
          name: 'Screen Protector Pack',
          unitsSold: 29,
          revenue: 435.00),
    ];

    return products.take(limit).toList();
  }

  /// Fetches sales chart data points for the given [period].
  Future<List<SalesDataPoint>> getSalesChart(String period) async {
    // TODO: Replace with actual API call to /seller/analytics/sales-chart?period=
    await Future.delayed(const Duration(milliseconds: 800));

    final count = _periodDataPoints(period);
    final labels = _generateLabels(period, count);

    return List.generate(count, (i) {
      final base = 300.0 + (i * 50.0) + (i % 3 == 0 ? 150 : 0);
      return SalesDataPoint(
        label: labels[i],
        revenue: base + (i * 17.0 % 200),
        orders: (base / 20).round(),
      );
    });
  }

  /// Fetches overall order statistics.
  Future<OrderStats> getOrderStats() async {
    // TODO: Replace with actual API call to /seller/analytics/orders
    await Future.delayed(const Duration(milliseconds: 600));

    return const OrderStats(
      totalOrders: 1247,
      pendingOrders: 23,
      completedOrders: 1189,
      cancelledOrders: 35,
      fulfillmentRate: 95.3,
    );
  }

  double _periodMultiplier(String period) {
    switch (period) {
      case '7d':
        return 0.25;
      case '30d':
        return 1.0;
      case '90d':
        return 2.8;
      case '1y':
        return 10.0;
      default:
        return 1.0;
    }
  }

  int _periodDataPoints(String period) {
    switch (period) {
      case '7d':
        return 7;
      case '30d':
        return 15;
      case '90d':
        return 12;
      case '1y':
        return 12;
      default:
        return 15;
    }
  }

  List<String> _generateLabels(String period, int count) {
    if (period == '7d') {
      return ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
    }
    if (period == '1y') {
      return [
        'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
        'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
      ];
    }
    return List.generate(count, (i) => 'Day ${i + 1}');
  }
}
