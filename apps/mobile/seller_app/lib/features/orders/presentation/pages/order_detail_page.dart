import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../data/order_repository.dart';
import '../widgets/order_status_badge.dart';

/// Displays detailed information about a single order, including items,
/// customer info, status timeline, and action buttons.
class SellerOrderDetailPage extends StatefulWidget {
  final String orderId;

  const SellerOrderDetailPage({super.key, required this.orderId});

  @override
  State<SellerOrderDetailPage> createState() => _SellerOrderDetailPageState();
}

class _SellerOrderDetailPageState extends State<SellerOrderDetailPage> {
  late Future<SellerOrder> _orderFuture;
  final _trackingController = TextEditingController();
  bool _isUpdating = false;

  @override
  void initState() {
    super.initState();
    _loadOrder();
  }

  @override
  void dispose() {
    _trackingController.dispose();
    super.dispose();
  }

  void _loadOrder() {
    setState(() {
      _orderFuture = getIt<SellerOrderRepository>().getOrderById(widget.orderId);
    });
  }

  Future<void> _updateStatus(String newStatus, {String? trackingNumber}) async {
    setState(() => _isUpdating = true);
    try {
      await getIt<SellerOrderRepository>().updateOrderStatus(
        widget.orderId,
        newStatus,
        trackingNumber: trackingNumber,
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Order updated to $newStatus')),
        );
        _loadOrder();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to update order: $e'),
            backgroundColor: Theme.of(context).colorScheme.error,
          ),
        );
      }
    } finally {
      if (mounted) setState(() => _isUpdating = false);
    }
  }

  void _showShipDialog() {
    _trackingController.clear();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Ship Order'),
        content: TextField(
          controller: _trackingController,
          decoration: const InputDecoration(
            labelText: 'Tracking Number',
            hintText: 'Enter tracking number',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              _updateStatus('shipped', trackingNumber: _trackingController.text.trim());
            },
            child: const Text('Ship Order'),
          ),
        ],
      ),
    );
  }

  int _statusStepIndex(String status) {
    switch (status.toLowerCase()) {
      case 'pending':
        return 0;
      case 'processing':
        return 1;
      case 'shipped':
        return 2;
      case 'delivered':
        return 3;
      case 'cancelled':
        return -1;
      default:
        return 0;
    }
  }

  String _formatDate(DateTime date) {
    final months = [
      'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
    ];
    return '${months[date.month - 1]} ${date.day}, ${date.year} at ${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Order Details'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: FutureBuilder<SellerOrder>(
        future: _orderFuture,
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
                  Text('Failed to load order', style: theme.textTheme.titleMedium),
                  const SizedBox(height: 8),
                  FilledButton(onPressed: _loadOrder, child: const Text('Retry')),
                ],
              ),
            );
          }

          final order = snapshot.data!;
          final stepIndex = _statusStepIndex(order.status);

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Order header
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          order.orderNumber,
                          style: theme.textTheme.headlineSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          _formatDate(order.createdAt),
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                    OrderStatusBadge(status: order.status),
                  ],
                ),
                const SizedBox(height: 24),

                // Order items
                Text(
                  'Items',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                ListView.separated(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: order.items.length,
                  separatorBuilder: (_, __) => const Divider(height: 1),
                  itemBuilder: (context, index) {
                    final item = order.items[index];
                    return Padding(
                      padding: const EdgeInsets.symmetric(vertical: 10),
                      child: Row(
                        children: [
                          ClipRRect(
                            borderRadius: BorderRadius.circular(8),
                            child: item.imageUrl != null
                                ? Image.network(
                                    item.imageUrl!,
                                    width: 56,
                                    height: 56,
                                    fit: BoxFit.cover,
                                    errorBuilder: (_, __, ___) => Container(
                                      width: 56,
                                      height: 56,
                                      color: Colors.grey.shade200,
                                      child: const Icon(Icons.image, color: Colors.grey),
                                    ),
                                  )
                                : Container(
                                    width: 56,
                                    height: 56,
                                    color: Colors.grey.shade200,
                                    child: const Icon(Icons.image, color: Colors.grey),
                                  ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  item.productName,
                                  style: theme.textTheme.bodyMedium?.copyWith(
                                    fontWeight: FontWeight.w600,
                                  ),
                                  maxLines: 2,
                                  overflow: TextOverflow.ellipsis,
                                ),
                                const SizedBox(height: 2),
                                Text(
                                  'Qty: ${item.quantity}  x  \$${item.price.toStringAsFixed(2)}',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: theme.colorScheme.onSurfaceVariant,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Text(
                            '\$${item.total.toStringAsFixed(2)}',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
                const Divider(),
                Align(
                  alignment: Alignment.centerRight,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    child: Text(
                      'Total: \$${order.total.toStringAsFixed(2)}',
                      style: theme.textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: theme.colorScheme.primary,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                // Customer info
                if (order.customerInfo != null) ...[
                  Text(
                    'Customer Information',
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              const Icon(Icons.person_outline, size: 18),
                              const SizedBox(width: 8),
                              Text(order.customerInfo!.name,
                                  style: theme.textTheme.bodyMedium?.copyWith(
                                    fontWeight: FontWeight.w600,
                                  )),
                            ],
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              const Icon(Icons.email_outlined, size: 18),
                              const SizedBox(width: 8),
                              Text(order.customerInfo!.email,
                                  style: theme.textTheme.bodyMedium),
                            ],
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              const Icon(Icons.phone_outlined, size: 18),
                              const SizedBox(width: 8),
                              Text(order.customerInfo!.phone,
                                  style: theme.textTheme.bodyMedium),
                            ],
                          ),
                          const SizedBox(height: 8),
                          Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Icon(Icons.location_on_outlined, size: 18),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  order.customerInfo!.fullAddress,
                                  style: theme.textTheme.bodyMedium,
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),
                ],

                // Status timeline
                Text(
                  'Status Timeline',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                if (order.status.toLowerCase() == 'cancelled')
                  Card(
                    color: Colors.red.shade50,
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Row(
                        children: [
                          Icon(Icons.cancel, color: Colors.red.shade700),
                          const SizedBox(width: 12),
                          Text(
                            'This order has been cancelled',
                            style: TextStyle(
                              color: Colors.red.shade700,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                  )
                else
                  Stepper(
                    currentStep: stepIndex,
                    controlsBuilder: (context, details) => const SizedBox.shrink(),
                    physics: const NeverScrollableScrollPhysics(),
                    steps: [
                      Step(
                        title: const Text('Pending'),
                        subtitle: const Text('Order placed by customer'),
                        content: const SizedBox.shrink(),
                        isActive: stepIndex >= 0,
                        state: stepIndex > 0 ? StepState.complete : StepState.indexed,
                      ),
                      Step(
                        title: const Text('Processing'),
                        subtitle: const Text('Order accepted and being prepared'),
                        content: const SizedBox.shrink(),
                        isActive: stepIndex >= 1,
                        state: stepIndex > 1 ? StepState.complete : StepState.indexed,
                      ),
                      Step(
                        title: const Text('Shipped'),
                        subtitle: Text(
                          order.trackingNumber != null
                              ? 'Tracking: ${order.trackingNumber}'
                              : 'Package shipped to customer',
                        ),
                        content: const SizedBox.shrink(),
                        isActive: stepIndex >= 2,
                        state: stepIndex > 2 ? StepState.complete : StepState.indexed,
                      ),
                      Step(
                        title: const Text('Delivered'),
                        subtitle: const Text('Order delivered successfully'),
                        content: const SizedBox.shrink(),
                        isActive: stepIndex >= 3,
                        state: stepIndex >= 3 ? StepState.complete : StepState.indexed,
                      ),
                    ],
                  ),
                const SizedBox(height: 24),

                // Action buttons
                if (!_isUpdating) ...[
                  if (order.status.toLowerCase() == 'pending')
                    SizedBox(
                      width: double.infinity,
                      child: FilledButton.icon(
                        onPressed: () => _updateStatus('processing'),
                        icon: const Icon(Icons.check_circle_outline),
                        label: const Text('Accept Order'),
                      ),
                    ),
                  if (order.status.toLowerCase() == 'processing')
                    SizedBox(
                      width: double.infinity,
                      child: FilledButton.icon(
                        onPressed: _showShipDialog,
                        icon: const Icon(Icons.local_shipping_outlined),
                        label: const Text('Ship Order'),
                      ),
                    ),
                  if (order.status.toLowerCase() == 'shipped')
                    SizedBox(
                      width: double.infinity,
                      child: FilledButton.icon(
                        onPressed: () => _updateStatus('delivered'),
                        icon: const Icon(Icons.done_all),
                        label: const Text('Mark Delivered'),
                      ),
                    ),
                ] else
                  const Center(child: CircularProgressIndicator()),
                const SizedBox(height: 24),
              ],
            ),
          );
        },
      ),
    );
  }
}
