import 'package:flutter/material.dart';

/// A badge widget that displays the order status with color-coded styling.
///
/// Colors are mapped as follows:
/// - pending: amber
/// - processing: blue
/// - shipped: indigo
/// - delivered: green
/// - cancelled: red
class OrderStatusBadge extends StatelessWidget {
  final String status;

  const OrderStatusBadge({super.key, required this.status});

  Color _backgroundColor() {
    switch (status.toLowerCase()) {
      case 'pending':
        return Colors.amber;
      case 'processing':
        return Colors.blue;
      case 'shipped':
        return Colors.indigo;
      case 'delivered':
        return Colors.green;
      case 'cancelled':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  Color _textColor() {
    switch (status.toLowerCase()) {
      case 'pending':
        return Colors.amber.shade900;
      case 'processing':
        return Colors.blue.shade900;
      case 'shipped':
        return Colors.indigo.shade900;
      case 'delivered':
        return Colors.green.shade900;
      case 'cancelled':
        return Colors.red.shade900;
      default:
        return Colors.grey.shade900;
    }
  }

  String _displayText() {
    if (status.isEmpty) return '';
    return status[0].toUpperCase() + status.substring(1).toLowerCase();
  }

  @override
  Widget build(BuildContext context) {
    final bgColor = _backgroundColor();
    final txtColor = _textColor();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor.withOpacity(0.15),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        _displayText(),
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: txtColor,
        ),
      ),
    );
  }
}
