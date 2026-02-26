import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../../orders/data/order_repository.dart';
import '../../data/return_repository.dart';

class ReturnRequestPage extends StatefulWidget {
  const ReturnRequestPage({super.key});

  @override
  State<ReturnRequestPage> createState() => _ReturnRequestPageState();
}

class _ReturnRequestPageState extends State<ReturnRequestPage> {
  final OrderRepository _orderRepo = getIt<OrderRepository>();
  final ReturnRepository _returnRepo = getIt<ReturnRepository>();
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  final TextEditingController _descriptionController = TextEditingController();

  List<Order> _orders = [];
  Order? _selectedOrder;
  final Map<String, bool> _selectedItems = {};
  String? _selectedReason;
  bool _isLoadingOrders = true;
  bool _isSubmitting = false;

  static const List<String> _reasons = [
    'Defective product',
    'Wrong item received',
    'Item not as described',
    'Changed my mind',
    'Better price found',
    'Arrived too late',
    'Other',
  ];

  @override
  void initState() {
    super.initState();
    _loadOrders();
  }

  @override
  void dispose() {
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _loadOrders() async {
    try {
      final result = await _orderRepo.getOrders(status: 'delivered');
      if (mounted) {
        setState(() {
          _orders = result.orders;
          _isLoadingOrders = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoadingOrders = false);
      }
    }
  }

  void _onOrderSelected(Order? order) {
    setState(() {
      _selectedOrder = order;
      _selectedItems.clear();
      if (order != null) {
        for (final item in order.items) {
          _selectedItems[item.id] = false;
        }
      }
    });
  }

  Future<void> _submitReturn() async {
    if (!_formKey.currentState!.validate()) return;

    final selectedItemIds = _selectedItems.entries
        .where((e) => e.value)
        .map((e) => e.key)
        .toList();

    if (selectedItemIds.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select at least one item to return')),
      );
      return;
    }

    setState(() => _isSubmitting = true);

    try {
      final items = selectedItemIds.map((id) {
        final orderItem = _selectedOrder!.items.firstWhere((i) => i.id == id);
        return {
          'orderItemId': id,
          'quantity': orderItem.quantity,
        };
      }).toList();

      await _returnRepo.createReturn(
        orderId: _selectedOrder!.id,
        items: items,
        reason: _selectedReason!,
        description: _descriptionController.text.trim().isNotEmpty
            ? _descriptionController.text.trim()
            : null,
      );

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Return request submitted successfully')),
        );
        context.pop(true);
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isSubmitting = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to submit return: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Request Return'),
      ),
      body: _isLoadingOrders
          ? const Center(child: CircularProgressIndicator())
          : Form(
              key: _formKey,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // Order dropdown
                  DropdownButtonFormField<Order>(
                    value: _selectedOrder,
                    decoration: const InputDecoration(
                      labelText: 'Select Order',
                      prefixIcon: Icon(Icons.receipt_long),
                      border: OutlineInputBorder(),
                    ),
                    isExpanded: true,
                    items: _orders.map((order) {
                      return DropdownMenuItem<Order>(
                        value: order,
                        child: Text(
                          'Order #${order.orderNumber} - \$${order.total.toStringAsFixed(2)}',
                          overflow: TextOverflow.ellipsis,
                        ),
                      );
                    }).toList(),
                    onChanged: _onOrderSelected,
                    validator: (value) {
                      if (value == null) return 'Please select an order';
                      return null;
                    },
                  ),
                  const SizedBox(height: 20),

                  // Items selection
                  if (_selectedOrder != null) ...[
                    Text(
                      'Select Items to Return',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    ..._selectedOrder!.items.map((item) {
                      return CheckboxListTile(
                        value: _selectedItems[item.id] ?? false,
                        onChanged: (value) {
                          setState(() {
                            _selectedItems[item.id] = value ?? false;
                          });
                        },
                        title: Text(item.name),
                        subtitle: Text(
                          'Qty: ${item.quantity} - \$${item.price.toStringAsFixed(2)}',
                        ),
                        secondary: ClipRRect(
                          borderRadius: BorderRadius.circular(8),
                          child: Image.network(
                            item.imageUrl,
                            width: 48,
                            height: 48,
                            fit: BoxFit.cover,
                            errorBuilder: (_, __, ___) => Container(
                              width: 48,
                              height: 48,
                              color: Colors.grey.shade200,
                              child: const Icon(Icons.image, size: 24),
                            ),
                          ),
                        ),
                        controlAffinity: ListTileControlAffinity.leading,
                        contentPadding: EdgeInsets.zero,
                      );
                    }),
                    const SizedBox(height: 20),
                  ],

                  // Reason dropdown
                  DropdownButtonFormField<String>(
                    value: _selectedReason,
                    decoration: const InputDecoration(
                      labelText: 'Reason for Return',
                      prefixIcon: Icon(Icons.help_outline),
                      border: OutlineInputBorder(),
                    ),
                    items: _reasons.map((reason) {
                      return DropdownMenuItem<String>(
                        value: reason,
                        child: Text(reason),
                      );
                    }).toList(),
                    onChanged: (value) {
                      setState(() => _selectedReason = value);
                    },
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please select a reason';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 20),

                  // Description
                  TextField(
                    controller: _descriptionController,
                    decoration: const InputDecoration(
                      labelText: 'Additional Details (optional)',
                      hintText: 'Describe the issue in more detail...',
                      prefixIcon: Icon(Icons.description),
                      border: OutlineInputBorder(),
                      alignLabelWithHint: true,
                    ),
                    maxLines: 4,
                    textInputAction: TextInputAction.newline,
                  ),
                  const SizedBox(height: 32),

                  // Submit button
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: _isSubmitting ? null : _submitReturn,
                      style: ElevatedButton.styleFrom(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: _isSubmitting
                          ? const SizedBox(
                              width: 24,
                              height: 24,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : const Text('Submit Return Request'),
                    ),
                  ),
                ],
              ),
            ),
    );
  }
}
