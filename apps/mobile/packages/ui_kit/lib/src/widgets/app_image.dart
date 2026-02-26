import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';
import '../theme/app_colors.dart';

/// A network image widget with shimmer loading placeholder and error fallback.
class AppImage extends StatelessWidget {
  /// The URL of the image to display.
  final String url;

  /// Optional width constraint.
  final double? width;

  /// Optional height constraint.
  final double? height;

  /// How to fit the image within its bounds.
  final BoxFit fit;

  const AppImage({
    super.key,
    required this.url,
    this.width,
    this.height,
    this.fit = BoxFit.cover,
  });

  @override
  Widget build(BuildContext context) {
    return CachedNetworkImage(
      imageUrl: url,
      width: width,
      height: height,
      fit: fit,
      placeholder: (context, url) => Shimmer.fromColors(
        baseColor: Colors.grey[300]!,
        highlightColor: Colors.grey[100]!,
        child: Container(
          width: width,
          height: height,
          color: Colors.white,
        ),
      ),
      errorWidget: (context, url, error) => Container(
        width: width,
        height: height,
        color: AppColors.divider,
        child: const Icon(
          Icons.broken_image_outlined,
          color: AppColors.textSecondary,
          size: 32,
        ),
      ),
    );
  }
}
