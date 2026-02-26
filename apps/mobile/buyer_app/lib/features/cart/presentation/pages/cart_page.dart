import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/di/injection.dart';
import '../../data/cart_repository.dart';
import '../widgets/cart_item_widget.dart';

class CartPage extends StatefulWidget {
  const CartPage({super.key});

  @override
  State<CartPage> createState() => _CartPageState();
}

class _CartPageState extends State<CartPage> {
  final CartRepository _cartRepo = getIt<CartRepository>();
  final TextEditingController _couponController = TextEditingController();

  Cart? _cart;
  bool _isLoading = true;
  bool _useLoyaltyPoints = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadCart();
  }

  @override
  void dispose() {
    _couponController.dispose();
    super.dispose();
  }

  Future<void> _loadCart() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final cart = await _cartRepo.getCart();
      if (mounted) {
        setState(() {
          _cart = cart;
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

  Future<void> _updateQuantity(String itemId, int quantity) async {
    try {
      final cart = await _cartRepo.updateQuantity(itemId: itemId, quantity: quantity);
      if (mounted) setState(() => _cart = cart);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to update quantity: $e')),
        );
      }
    }
  }

  Future<void> _removeItem(String itemId) async {
    try {
      final cart = await _cartRepo.removeFromCart(itemId);
      if (mounted) setState(() => _cart = cart);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to remove item: $e')),
        );
      }
    }
  }

  Future<void> _applyCoupon() async {
    final code = _couponController.text.trim();
    if (code.isEmpty) return;

    try {
      final cart = await _cartRepo.applyCoupon(code);
      if (mounted) {
        setState(() => _cart = cart);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Coupon applied successfully')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Invalid coupon code: $e')),
        );
      }
    }
  }

  Future<void> _toggleLoyaltyPoints(bool value) async {
    setState(() => _useLoyaltyPoints = value);
    if (value) {
      try {
        final cart = await _cartRepo.applyPoints(0); // apply max available
        if (mounted) setState(() => _cart = cart);
      } catch (e) {
        if (mounted) {
          setState(() => _useLoyaltyPoints = false);
        }
      }
    } else {
      try {
        final cart = await _cartRepo.applyPoints(0);
        if (mounted) setState(() => _cart = cart);
      } catch (_) {}
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(_cart != null ? 'Cart (${_cart!.itemCount})' : 'Cart'),
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
                      Text('Failed to load cart', style: theme.textTheme.titleMedium),
                      const SizedBox(height: 8),
                      FilledButton(onPressed: _loadCart, child: const Text('Retry')),
                    ],
                  ),
                )
              : _cart == null || _cart!.items.isEmpty
                  ? _buildEmptyState(theme)
                  : _buildCartContent(theme),
    );
  }

  Widget _buildEmptyState(ThemeData theme) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.shopping_cart_outlined, size: 80, color: Colors.grey.shade300),
          const SizedBox(height: 16),
          Text(
            'Your cart is empty',
            style: theme.textTheme.titleLarge?.copyWith(color: Colors.grey),
          ),
          const SizedBox(height: 8),
          Text(
            'Add items to start shopping',
            style: theme.textTheme.bodyMedium?.copyWith(color: Colors.grey),
          ),
          const SizedBox(height: 24),
          FilledButton(
            onPressed: () => context.go('/'),
            child: const Text('Start Shopping'),
          ),
        ],
      ),
    );
  }

  Widget _buildCartContent(ThemeData theme) {
    final cart = _cart!;

    return Column(
      children: [
        Expanded(
          child: ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Cart items
              ...cart.items.map(
                (item) => CartItemWidget(
                  item: item,
                  onQuantityChanged: (qty) => _updateQuantity(item.id, qty),
                  onRemove: () => _removeItem(item.id),
                ),
              ),
              const SizedBox(height: 16),

              // Coupon code
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _couponController,
                      decoration: const InputDecoration(
                        hintText: 'Enter coupon code',
                        prefixIcon: Icon(Icons.local_offer_outlined),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  FilledButton(
                    onPressed: _applyCoupon,
                    child: const Text('Apply'),
                  ),
                ],
              ),
              if (cart.couponCode != null) ...[
                const SizedBox(height: 8),
                Chip(
                  label: Text('Coupon: ${cart.couponCode}'),
                  onDeleted: () {},
                  deleteIcon: const Icon(Icons.close, size: 16),
                ),
              ],
              const SizedBox(height: 16),

              // Loyalty points toggle
              SwitchListTile(
                title: const Text('Use Loyalty Points'),
                subtitle: Text(
                  cart.loyaltyPointsApplied > 0
                      ? '${cart.loyaltyPointsApplied} points applied'
                      : 'Apply your available loyalty points',
                ),
                value: _useLoyaltyPoints,
                onChanged: _toggleLoyaltyPoints,
                contentPadding: EdgeInsets.zero,
              ),
              const SizedBox(height: 16),

              // Order summary
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Order Summary',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 12),
                      _buildSummaryRow('Subtotal', '\$${cart.subtotal.toStringAsFixed(2)}'),
                      const SizedBox(height: 8),
                      _buildSummaryRow(
                        'Shipping',
                        cart.shippingEstimate > 0
                            ? '\$${cart.shippingEstimate.toStringAsFixed(2)}'
                            : 'Free',
                      ),
                      if (cart.discount > 0) ...[
                        const SizedBox(height: 8),
                        _buildSummaryRow(
                          'Discount',
                          '-\$${cart.discount.toStringAsFixed(2)}',
                          valueColor: Colors.green,
                        ),
                      ],
                      if (cart.tax > 0) ...[
                        const SizedBox(height: 8),
                        _buildSummaryRow('Tax', '\$${cart.tax.toStringAsFixed(2)}'),
                      ],
                      const Divider(height: 24),
                      _buildSummaryRow(
                        'Total',
                        '\$${cart.total.toStringAsFixed(2)}',
                        isBold: true,
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),

        // Checkout button
        SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: SizedBox(
              width: double.infinity,
              height: 52,
              child: ElevatedButton(
                onPressed: () => context.push('/checkout'),
                style: ElevatedButton.styleFrom(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: Text('Checkout (\$${cart.total.toStringAsFixed(2)})'),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSummaryRow(String label, String value,
      {bool isBold = false, Color? valueColor}) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: TextStyle(
            fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
            fontSize: isBold ? 16 : 14,
          ),
        ),
        Text(
          value,
          style: TextStyle(
            fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
            fontSize: isBold ? 16 : 14,
            color: valueColor,
          ),
        ),
      ],
    );
  }
}
