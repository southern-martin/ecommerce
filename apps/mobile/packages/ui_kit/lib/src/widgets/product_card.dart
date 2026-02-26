import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';
import 'package:flutter/material.dart';
import '../theme/app_colors.dart';
import '../theme/app_spacing.dart';
import 'app_card.dart';
import 'app_image.dart';
import 'price_text.dart';
import 'rating_bar.dart';

/// A card displaying a product's image, name, price, rating, and an
/// add-to-cart action.
class ProductCard extends StatelessWidget {
  /// The product to display.
  final Product product;

  /// Called when the card is tapped.
  final VoidCallback? onTap;

  /// Called when the add-to-cart button is pressed.
  final VoidCallback? onAddToCart;

  const ProductCard({
    super.key,
    required this.product,
    this.onTap,
    this.onAddToCart,
  });

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: EdgeInsets.zero,
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Product image
          AspectRatio(
            aspectRatio: 1,
            child: product.primaryImageUrl != null
                ? AppImage(
                    url: product.primaryImageUrl!,
                    fit: BoxFit.cover,
                  )
                : Container(
                    color: AppColors.divider,
                    child: const Icon(
                      Icons.image_outlined,
                      size: 48,
                      color: AppColors.textSecondary,
                    ),
                  ),
          ),

          Padding(
            padding: const EdgeInsets.all(AppSpacing.sm),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Product name
                Text(
                  product.name,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: AppSpacing.xs),

                // Price
                PriceText(
                  priceCents: product.priceCents,
                  compareAtPriceCents: product.compareAtPriceCents,
                  fontSize: 14,
                ),
                const SizedBox(height: AppSpacing.xs),

                // Rating and cart button row
                Row(
                  children: [
                    Expanded(
                      child: RatingBar(
                        rating: product.rating,
                        size: 14,
                        showCount: true,
                        count: product.reviewCount,
                      ),
                    ),
                    if (onAddToCart != null)
                      SizedBox(
                        width: 32,
                        height: 32,
                        child: IconButton(
                          onPressed: onAddToCart,
                          icon: const Icon(Icons.add_shopping_cart, size: 18),
                          padding: EdgeInsets.zero,
                          color: AppColors.primary,
                          tooltip: 'Add to cart',
                        ),
                      ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
