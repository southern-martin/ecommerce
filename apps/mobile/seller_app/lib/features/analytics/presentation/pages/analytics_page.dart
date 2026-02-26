import 'package:flutter/material.dart';

import '../../../../core/di/injection.dart';
import '../../data/analytics_repository.dart';

/// Analytics dashboard showing revenue summary, sales chart placeholder,
/// and top-selling products for the seller.
class AnalyticsPage extends StatefulWidget {
  const AnalyticsPage({super.key});

  @override
  State<AnalyticsPage> createState() => _AnalyticsPageState();
}

class _AnalyticsPageState extends State<AnalyticsPage> {
  String _selectedPeriod = '30d';

  late Future<RevenueSummary> _revenueFuture;
  late Future<List<SalesDataPoint>> _chartFuture;
  late Future<List<TopProduct>> _topProductsFuture;

  static const _periods = ['7d', '30d', '90d', '1y'];
  static const _periodLabels = {
    '7d': '7D',
    '30d': '30D',
    '90d': '90D',
    '1y': '1Y',
  };

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  void _loadData() {
    final repo = getIt<AnalyticsRepository>();
    _revenueFuture = repo.getRevenueSummary(_selectedPeriod);
    _chartFuture = repo.getSalesChart(_selectedPeriod);
    _topProductsFuture = repo.getTopProducts(10);
  }

  void _onPeriodChanged(Set<String> selected) {
    setState(() {
      _selectedPeriod = selected.first;
      _loadData();
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Analytics'),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          setState(() => _loadData());
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Period selector
              Center(
                child: SegmentedButton<String>(
                  segments: _periods
                      .map((p) => ButtonSegment(
                            value: p,
                            label: Text(_periodLabels[p]!),
                          ))
                      .toList(),
                  selected: {_selectedPeriod},
                  onSelectionChanged: _onPeriodChanged,
                ),
              ),
              const SizedBox(height: 20),

              // Revenue summary cards
              FutureBuilder<RevenueSummary>(
                future: _revenueFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const SizedBox(
                      height: 100,
                      child: Center(child: CircularProgressIndicator()),
                    );
                  }

                  if (snapshot.hasError) {
                    return Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Text(
                          'Failed to load revenue data',
                          style:
                              TextStyle(color: theme.colorScheme.error),
                        ),
                      ),
                    );
                  }

                  final summary = snapshot.data!;
                  return Row(
                    children: [
                      Expanded(
                        child: _SummaryCard(
                          title: 'Total Revenue',
                          value:
                              '\$${summary.totalRevenue.toStringAsFixed(2)}',
                          icon: Icons.attach_money,
                          color: Colors.green,
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: _SummaryCard(
                          title: 'Avg Order',
                          value:
                              '\$${summary.averageOrderValue.toStringAsFixed(2)}',
                          icon: Icons.shopping_cart_outlined,
                          color: Colors.blue,
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: _SummaryCard(
                          title: 'Total Orders',
                          value: summary.totalOrders.toString(),
                          icon: Icons.receipt_long,
                          color: Colors.orange,
                        ),
                      ),
                    ],
                  );
                },
              ),
              const SizedBox(height: 24),

              // Sales chart placeholder
              Text(
                'Sales Overview',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 12),
              FutureBuilder<List<SalesDataPoint>>(
                future: _chartFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return Container(
                      height: 200,
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surfaceContainerLow,
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child:
                          const Center(child: CircularProgressIndicator()),
                    );
                  }

                  if (snapshot.hasError) {
                    return Container(
                      height: 200,
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surfaceContainerLow,
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Center(
                        child: Text(
                          'Failed to load chart data',
                          style:
                              TextStyle(color: theme.colorScheme.error),
                        ),
                      ),
                    );
                  }

                  final data = snapshot.data!;
                  return Container(
                    height: 200,
                    width: double.infinity,
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surfaceContainerLow,
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(
                        color: theme.colorScheme.outlineVariant,
                      ),
                    ),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.bar_chart_rounded,
                            size: 48,
                            color: theme.colorScheme.primary
                                .withOpacity(0.5)),
                        const SizedBox(height: 8),
                        Text(
                          'Line Chart Placeholder',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${data.length} data points loaded',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
              const SizedBox(height: 24),

              // Top products
              Text(
                'Top Products',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 12),
              FutureBuilder<List<TopProduct>>(
                future: _topProductsFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const SizedBox(
                      height: 200,
                      child: Center(child: CircularProgressIndicator()),
                    );
                  }

                  if (snapshot.hasError) {
                    return Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Text(
                          'Failed to load top products',
                          style:
                              TextStyle(color: theme.colorScheme.error),
                        ),
                      ),
                    );
                  }

                  final products = snapshot.data!;

                  if (products.isEmpty) {
                    return const Card(
                      child: Padding(
                        padding: EdgeInsets.all(24),
                        child: Center(child: Text('No product data yet')),
                      ),
                    );
                  }

                  return ListView.separated(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    itemCount: products.length,
                    separatorBuilder: (_, __) => const SizedBox(height: 6),
                    itemBuilder: (context, index) {
                      final product = products[index];
                      final rank = index + 1;
                      return Card(
                        margin: EdgeInsets.zero,
                        child: Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 14, vertical: 12),
                          child: Row(
                            children: [
                              // Rank number
                              SizedBox(
                                width: 32,
                                child: Text(
                                  '#$rank',
                                  style:
                                      theme.textTheme.titleMedium?.copyWith(
                                    fontWeight: FontWeight.bold,
                                    color: rank <= 3
                                        ? theme.colorScheme.primary
                                        : theme
                                            .colorScheme.onSurfaceVariant,
                                  ),
                                ),
                              ),
                              const SizedBox(width: 8),
                              // Product name
                              Expanded(
                                child: Column(
                                  crossAxisAlignment:
                                      CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      product.name,
                                      style: theme.textTheme.bodyMedium
                                          ?.copyWith(
                                        fontWeight: FontWeight.w600,
                                      ),
                                      maxLines: 1,
                                      overflow: TextOverflow.ellipsis,
                                    ),
                                    const SizedBox(height: 2),
                                    Text(
                                      '${product.unitsSold} units sold',
                                      style:
                                          theme.textTheme.bodySmall?.copyWith(
                                        color: theme
                                            .colorScheme.onSurfaceVariant,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              // Revenue
                              Text(
                                '\$${product.revenue.toStringAsFixed(0)}',
                                style:
                                    theme.textTheme.titleSmall?.copyWith(
                                  fontWeight: FontWeight.bold,
                                  color: theme.colorScheme.primary,
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  );
                },
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}

class _SummaryCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  const _SummaryCard({
    required this.title,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(icon, size: 20, color: color),
            const SizedBox(height: 8),
            Text(
              value,
              style: theme.textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 2),
            Text(
              title,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }
}
