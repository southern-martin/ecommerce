import 'package:flutter/material.dart';
import '../theme/app_colors.dart';
import '../theme/app_spacing.dart';

/// The visual variant of an [AppButton].
enum AppButtonVariant { primary, secondary, outline, ghost, destructive }

/// The size of an [AppButton].
enum AppButtonSize { sm, md, lg }

/// A configurable button widget that supports multiple variants, sizes,
/// loading state, and full-width layout.
class AppButton extends StatelessWidget {
  /// Optional child widget. Takes precedence over [text].
  final Widget? child;

  /// Text label for the button. Ignored when [child] is provided.
  final String? text;

  /// Called when the button is pressed. When `null` the button is disabled.
  final VoidCallback? onPressed;

  /// Visual variant of the button.
  final AppButtonVariant variant;

  /// Size of the button.
  final AppButtonSize size;

  /// Whether to show a loading indicator instead of the label.
  final bool isLoading;

  /// Whether the button should expand to fill its parent width.
  final bool isFullWidth;

  const AppButton({
    super.key,
    this.child,
    this.text,
    this.onPressed,
    this.variant = AppButtonVariant.primary,
    this.size = AppButtonSize.md,
    this.isLoading = false,
    this.isFullWidth = false,
  }) : assert(child != null || text != null,
            'Either child or text must be provided');

  EdgeInsetsGeometry get _padding {
    switch (size) {
      case AppButtonSize.sm:
        return const EdgeInsets.symmetric(horizontal: 12, vertical: 6);
      case AppButtonSize.md:
        return const EdgeInsets.symmetric(horizontal: 24, vertical: 12);
      case AppButtonSize.lg:
        return const EdgeInsets.symmetric(horizontal: 32, vertical: 16);
    }
  }

  double get _fontSize {
    switch (size) {
      case AppButtonSize.sm:
        return 12;
      case AppButtonSize.md:
        return 14;
      case AppButtonSize.lg:
        return 16;
    }
  }

  double get _loaderSize {
    switch (size) {
      case AppButtonSize.sm:
        return 14;
      case AppButtonSize.md:
        return 18;
      case AppButtonSize.lg:
        return 22;
    }
  }

  Widget _buildChild(Color foreground) {
    if (isLoading) {
      return SizedBox(
        width: _loaderSize,
        height: _loaderSize,
        child: CircularProgressIndicator(
          strokeWidth: 2,
          valueColor: AlwaysStoppedAnimation<Color>(foreground),
        ),
      );
    }
    return child ??
        Text(
          text!,
          style: TextStyle(
            fontSize: _fontSize,
            fontWeight: FontWeight.w600,
          ),
        );
  }

  @override
  Widget build(BuildContext context) {
    final effectiveOnPressed = isLoading ? null : onPressed;

    Widget button;

    switch (variant) {
      case AppButtonVariant.primary:
        button = ElevatedButton(
          onPressed: effectiveOnPressed,
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.primary,
            foregroundColor: Colors.white,
            padding: _padding,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.md),
            ),
            elevation: 0,
          ),
          child: _buildChild(Colors.white),
        );
        break;
      case AppButtonVariant.secondary:
        button = ElevatedButton(
          onPressed: effectiveOnPressed,
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.secondary,
            foregroundColor: Colors.white,
            padding: _padding,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.md),
            ),
            elevation: 0,
          ),
          child: _buildChild(Colors.white),
        );
        break;
      case AppButtonVariant.outline:
        button = OutlinedButton(
          onPressed: effectiveOnPressed,
          style: OutlinedButton.styleFrom(
            foregroundColor: AppColors.primary,
            side: const BorderSide(color: AppColors.primary),
            padding: _padding,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.md),
            ),
          ),
          child: _buildChild(AppColors.primary),
        );
        break;
      case AppButtonVariant.ghost:
        button = TextButton(
          onPressed: effectiveOnPressed,
          style: TextButton.styleFrom(
            foregroundColor: AppColors.primary,
            padding: _padding,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.md),
            ),
          ),
          child: _buildChild(AppColors.primary),
        );
        break;
      case AppButtonVariant.destructive:
        button = ElevatedButton(
          onPressed: effectiveOnPressed,
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.error,
            foregroundColor: Colors.white,
            padding: _padding,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.md),
            ),
            elevation: 0,
          ),
          child: _buildChild(Colors.white),
        );
        break;
    }

    if (isFullWidth) {
      return SizedBox(width: double.infinity, child: button);
    }

    return button;
  }
}
