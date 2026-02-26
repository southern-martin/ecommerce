import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

/// Displays a formatted price with optional compare-at price and discount badge.
class PriceText extends StatelessWidget {
  /// Current price in cents.
  final int priceCents;

  /// Original price in cents (displayed with strikethrough when provided).
  final int? compareAtPriceCents;

  /// Currency symbol.
  final String currency;

  /// Font size for the current price.
  final double fontSize;

  const PriceText({
    super.key,
    required this.priceCents,
    this.compareAtPriceCents,
    this.currency = '\$',
    this.fontSize = 16,
  });

  String _formatPrice(int cents) {
    final dollars = cents / 100;
    return '$currency${dollars.toStringAsFixed(2)}';
  }

  bool get _isOnSale =>
      compareAtPriceCents != null && compareAtPriceCents! > priceCents;

  int get _discountPercentage {
    if (!_isOnSale) return 0;
    return (((compareAtPriceCents! - priceCents) / compareAtPriceCents!) * 100)
        .round();
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Text(
          _formatPrice(priceCents),
          style: TextStyle(
            fontSize: fontSize,
            fontWeight: FontWeight.bold,
            color: _isOnSale ? AppColors.error : AppColors.textPrimary,
          ),
        ),
        if (_isOnSale) ...[
          const SizedBox(width: 6),
          Text(
            _formatPrice(compareAtPriceCents!),
            style: TextStyle(
              fontSize: fontSize * 0.8,
              color: AppColors.textSecondary,
              decoration: TextDecoration.lineThrough,
            ),
          ),
          const SizedBox(width: 6),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(
              color: AppColors.error.withOpacity(0.1),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              '-$_discountPercentage%',
              style: TextStyle(
                fontSize: fontSize * 0.7,
                fontWeight: FontWeight.w600,
                color: AppColors.error,
              ),
            ),
          ),
        ],
      ],
    );
  }
}
