import 'package:flutter/material.dart';

/// Displays a row of star icons representing a rating value.
class RatingBar extends StatelessWidget {
  /// The rating value (0.0 to 5.0).
  final double rating;

  /// Size of each star icon.
  final double size;

  /// Color of filled and half-filled stars.
  final Color color;

  /// Whether to show the review count label.
  final bool showCount;

  /// Number of reviews (displayed when [showCount] is true).
  final int? count;

  const RatingBar({
    super.key,
    required this.rating,
    this.size = 16,
    this.color = Colors.amber,
    this.showCount = false,
    this.count,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        ...List.generate(5, (index) {
          final starValue = index + 1;
          IconData icon;
          if (rating >= starValue) {
            icon = Icons.star;
          } else if (rating >= starValue - 0.5) {
            icon = Icons.star_half;
          } else {
            icon = Icons.star_border;
          }
          return Icon(icon, size: size, color: color);
        }),
        if (showCount && count != null) ...[
          const SizedBox(width: 4),
          Text(
            '($count)',
            style: TextStyle(
              fontSize: size * 0.75,
              color: Colors.grey[600],
            ),
          ),
        ],
      ],
    );
  }
}
