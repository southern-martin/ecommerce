import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../data/return_repository.dart';
import '../../../orders/presentation/widgets/order_status_badge.dart';

/// Displays a list of return requests with a status filter dropdown.
class SellerReturnListPage extends StatefulWidget {
  const SellerReturnListPage({super.key});

  @override
  State<SellerReturnListPage> createState() => _SellerReturnListPageState();
}

class _SellerReturnListPageState extends State<SellerReturnListPage> {
  late Future<PaginatedReturns> _returnsFuture;
  String _selectedStatus = 'All';

  static const _statusOptions = ['All', 'Pending', 'Approved', 'Rejected'];

  @override
  void initState() {
    super.initState();
    _loadReturns();
  }

  void _loadReturns() {
    final status = _selectedStatus == 'All' ? null : _selectedStatus.toLowerCase();
    setState(() {
      _returnsFuture = getIt<SellerReturnRepository>().getReturns(status: status);
    });
  }

  String _formatDate(DateTime date) {
    final months = [
      'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
    ];
    return '${months[date.month - 1]} ${date.day}, ${date.year}';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Returns'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Column(
        children: [
          // Status filter
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            child: Row(
              children: [
                Text('Filter by: ', style: theme.textTheme.bodyMedium),
                const SizedBox(width: 8),
                DropdownButton<String>(
                  value: _selectedStatus,
                  onChanged: (value) {
                    if (value != null) {
                      setState(() => _selectedStatus = value);
                      _loadReturns();
                    }
                  },
                  items: _statusOptions.map((s) {
                    return DropdownMenuItem(value: s, child: Text(s));
                  }).toList(),
                ),
              ],
            ),
          ),

          // Returns list
          Expanded(
            child: FutureBuilder<PaginatedReturns>(
              future: _returnsFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                }

                if (snapshot.hasError) {
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.error_outline, size: 48, color: theme.colorScheme.error),
                        const SizedBox(height: 16),
                        Text('Failed to load returns', style: theme.textTheme.titleMedium),
                        const SizedBox(height: 8),
                        FilledButton(onPressed: _loadReturns, child: const Text('Retry')),
                      ],
                    ),
                  );
                }

                final returns = snapshot.data!.returns;

                if (returns.isEmpty) {
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.assignment_return_outlined,
                            size: 64, color: theme.colorScheme.onSurfaceVariant),
                        const SizedBox(height: 16),
                        Text('No returns found', style: theme.textTheme.titleMedium),
                      ],
                    ),
                  );
                }

                return RefreshIndicator(
                  onRefresh: () async => _loadReturns(),
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    itemCount: returns.length,
                    itemBuilder: (context, index) {
                      final ret = returns[index];
                      return Card(
                        margin: const EdgeInsets.only(bottom: 10),
                        child: InkWell(
                          onTap: () => context.push('/returns/${ret.id}'),
                          borderRadius: BorderRadius.circular(12),
                          child: Padding(
                            padding: const EdgeInsets.all(14),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                  children: [
                                    Text(
                                      ret.returnNumber,
                                      style: theme.textTheme.titleSmall?.copyWith(
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                    OrderStatusBadge(status: ret.status),
                                  ],
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  'Order: ${ret.orderNumber}',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: theme.colorScheme.onSurfaceVariant,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Row(
                                  children: [
                                    Icon(Icons.person_outline,
                                        size: 16, color: theme.colorScheme.onSurfaceVariant),
                                    const SizedBox(width: 4),
                                    Text(ret.customerName, style: theme.textTheme.bodyMedium),
                                  ],
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  ret.reason,
                                  style: theme.textTheme.bodySmall,
                                  maxLines: 2,
                                  overflow: TextOverflow.ellipsis,
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  _formatDate(ret.createdAt),
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: theme.colorScheme.onSurfaceVariant,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}
