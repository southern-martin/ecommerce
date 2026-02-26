import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';

import '../../../../core/di/injection.dart';
import '../../data/order_repository.dart';

class OrderDetailPage extends StatefulWidget {
  final String orderId;

  const OrderDetailPage({super.key, required this.orderId});

  @override
  State<OrderDetailPage> createState() => _OrderDetailPageState();
}

class _OrderDetailPageState extends State<OrderDetailPage> {
  final OrderRepository _orderRepo = getIt<OrderRepository>();

  Order? _order;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadOrder();
  }

  Future<void> _loadOrder() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final order = await _orderRepo.getOrderById(widget.orderId);
      if (mounted) {
        setState(() {
          _order = order;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _error = e.toString();
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _cancelOrder() async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Cancel Order'),
        content: const Text('Are you sure you want to cancel this order?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('No'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Cancel Order'),
          ),
        ],
      ),
    );

    if (confirm == true) {
      try {
        await _orderRepo.cancelOrder(widget.orderId);
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Order cancelled successfully')),
          );
          _loadOrder();
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to cancel order: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(_order != null ? 'Order #${_order!.orderNumber}' : 'Order Details'),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(Icons.error_outline, size: 64, color: Colors.grey.shade400),
                      const SizedBox(height: 16),
                      Text('Failed to load order', style: theme.textTheme.titleMedium),
                      const SizedBox(height: 8),
                      FilledButton(onPressed: _loadOrder, child: const Text('Retry')),
                    ],
                  ),
                )
              : _buildContent(theme),
    );
  }

  Widget _buildContent(ThemeData theme) {
    final order = _order!;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Order items
          Text(
            'Items',
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          ...order.items.map((item) => _buildItemRow(theme, item)),
          const SizedBox(height: 16),

          // Shipping address
          if (order.shippingAddress != null) ...[
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Shipping Address',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(order.shippingAddress!.name),
                    Text(order.shippingAddress!.street),
                    Text(
                      '${order.shippingAddress!.city}, ${order.shippingAddress!.state} ${order.shippingAddress!.zip}',
                    ),
                    if (order.shippingAddress!.phone.isNotEmpty)
                      Text(order.shippingAddress!.phone),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),
          ],

          // Order timeline
          if (order.timeline.isNotEmpty) ...[
            Text(
              'Order Timeline',
              style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 12),
            Stepper(
              physics: const NeverScrollableScrollPhysics(),
              currentStep: order.timeline.lastIndexWhere((s) => s.isCompleted),
              controlsBuilder: (context, details) => const SizedBox.shrink(),
              steps: order.timeline.map((step) {
                return Step(
                  title: Text(step.label),
                  subtitle: step.date != null
                      ? Text(
                          '${step.date!.day}/${step.date!.month}/${step.date!.year} '
                          '${step.date!.hour}:${step.date!.minute.toString().padLeft(2, '0')}',
                        )
                      : null,
                  isActive: step.isCompleted,
                  state: step.isCompleted ? StepState.complete : StepState.indexed,
                  content: const SizedBox.shrink(),
                );
              }).toList(),
            ),
            const SizedBox(height: 16),
          ],

          // Action buttons
          _buildActionButtons(theme, order),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  Widget _buildItemRow(ThemeData theme, OrderItem item) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: CachedNetworkImage(
              imageUrl: item.imageUrl,
              width: 64,
              height: 64,
              fit: BoxFit.cover,
              placeholder: (_, __) => Container(
                width: 64,
                height: 64,
                color: Colors.grey.shade100,
              ),
              errorWidget: (_, __, ___) => Container(
                width: 64,
                height: 64,
                color: Colors.grey.shade200,
                child: const Icon(Icons.image_not_supported, size: 20),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  item.name,
                  style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                if (item.variantLabel != null)
                  Text(
                    item.variantLabel!,
                    style: theme.textTheme.bodySmall?.copyWith(color: Colors.grey),
                  ),
                Text(
                  'Qty: ${item.quantity}',
                  style: theme.textTheme.bodySmall?.copyWith(color: Colors.grey),
                ),
              ],
            ),
          ),
          Text(
            '\$${(item.price * item.quantity).toStringAsFixed(2)}',
            style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }

  Widget _buildActionButtons(ThemeData theme, Order order) {
    final status = order.status.toLowerCase();

    return Wrap(
      spacing: 12,
      runSpacing: 12,
      children: [
        if (status == 'shipped' || status == 'in_transit')
          OutlinedButton.icon(
            onPressed: () => context.push('/tracking/${order.id}'),
            icon: const Icon(Icons.local_shipping_outlined),
            label: const Text('Track Shipment'),
          ),
        if (status == 'delivered')
          OutlinedButton.icon(
            onPressed: () => context.push('/account/returns/new'),
            icon: const Icon(Icons.assignment_return_outlined),
            label: const Text('Request Return'),
          ),
        if (status == 'processing' || status == 'pending')
          OutlinedButton.icon(
            onPressed: _cancelOrder,
            icon: const Icon(Icons.cancel_outlined, color: Colors.red),
            label: const Text('Cancel Order', style: TextStyle(color: Colors.red)),
            style: OutlinedButton.styleFrom(
              side: const BorderSide(color: Colors.red),
            ),
          ),
      ],
    );
  }
}
