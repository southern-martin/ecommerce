import 'package:flutter/material.dart';

/// Overlays a numeric badge on top of a child widget (e.g., a cart icon).
/// The badge is hidden when [count] is zero or negative.
class BadgeWidget extends StatelessWidget {
  /// The widget to display behind the badge.
  final Widget child;

  /// The number to show in the badge.
  final int count;

  /// Background color of the badge circle.
  final Color color;

  const BadgeWidget({
    super.key,
    required this.child,
    required this.count,
    this.color = Colors.red,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        child,
        if (count > 0)
          Positioned(
            right: -6,
            top: -6,
            child: Container(
              padding: const EdgeInsets.all(4),
              constraints: const BoxConstraints(minWidth: 18, minHeight: 18),
              decoration: BoxDecoration(
                color: color,
                shape: BoxShape.circle,
              ),
              child: Center(
                child: Text(
                  count > 99 ? '99+' : count.toString(),
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 10,
                    fontWeight: FontWeight.bold,
                  ),
                  textAlign: TextAlign.center,
                ),
              ),
            ),
          ),
      ],
    );
  }
}
