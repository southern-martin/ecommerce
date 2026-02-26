import 'package:json_annotation/json_annotation.dart';

part 'product.g.dart';

/// Represents a product listing in the marketplace.
@JsonSerializable()
class Product {
  final String id;
  final String name;
  final String slug;
  final String description;
  final int priceCents;
  final int? compareAtPriceCents;
  final List<ProductImage> images;
  final String? category;
  final double rating;
  final int reviewCount;
  final bool inStock;
  final int stockQuantity;
  final String? seller;
  final List<ProductVariant> variants;

  const Product({
    required this.id,
    required this.name,
    required this.slug,
    required this.description,
    required this.priceCents,
    this.compareAtPriceCents,
    this.images = const [],
    this.category,
    this.rating = 0.0,
    this.reviewCount = 0,
    this.inStock = true,
    this.stockQuantity = 0,
    this.seller,
    this.variants = const [],
  });

  /// Returns the primary image URL, or `null` if there are no images.
  String? get primaryImageUrl {
    final primary = images.where((img) => img.isPrimary).toList();
    if (primary.isNotEmpty) return primary.first.url;
    return images.isNotEmpty ? images.first.url : null;
  }

  /// Whether the product is currently on sale.
  bool get isOnSale =>
      compareAtPriceCents != null && compareAtPriceCents! > priceCents;

  /// Discount percentage when the product is on sale, or 0.
  int get discountPercentage {
    if (!isOnSale) return 0;
    return (((compareAtPriceCents! - priceCents) / compareAtPriceCents!) * 100)
        .round();
  }

  /// Whether the product has variants.
  bool get hasVariants => variants.isNotEmpty;

  factory Product.fromJson(Map<String, dynamic> json) =>
      _$ProductFromJson(json);

  Map<String, dynamic> toJson() => _$ProductToJson(this);

  Product copyWith({
    String? id,
    String? name,
    String? slug,
    String? description,
    int? priceCents,
    int? compareAtPriceCents,
    List<ProductImage>? images,
    String? category,
    double? rating,
    int? reviewCount,
    bool? inStock,
    int? stockQuantity,
    String? seller,
    List<ProductVariant>? variants,
  }) {
    return Product(
      id: id ?? this.id,
      name: name ?? this.name,
      slug: slug ?? this.slug,
      description: description ?? this.description,
      priceCents: priceCents ?? this.priceCents,
      compareAtPriceCents: compareAtPriceCents ?? this.compareAtPriceCents,
      images: images ?? this.images,
      category: category ?? this.category,
      rating: rating ?? this.rating,
      reviewCount: reviewCount ?? this.reviewCount,
      inStock: inStock ?? this.inStock,
      stockQuantity: stockQuantity ?? this.stockQuantity,
      seller: seller ?? this.seller,
      variants: variants ?? this.variants,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Product && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() => 'Product(id: $id, name: $name, price: $priceCents)';
}

/// An image associated with a product.
@JsonSerializable()
class ProductImage {
  final String id;
  final String url;
  final String? alt;
  final bool isPrimary;

  const ProductImage({
    required this.id,
    required this.url,
    this.alt,
    this.isPrimary = false,
  });

  factory ProductImage.fromJson(Map<String, dynamic> json) =>
      _$ProductImageFromJson(json);

  Map<String, dynamic> toJson() => _$ProductImageToJson(this);

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ProductImage &&
          runtimeType == other.runtimeType &&
          id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// A variant of a product (e.g., size, color).
@JsonSerializable()
class ProductVariant {
  final String id;
  final String name;
  final String value;
  final int priceModifier;
  final int stockQuantity;

  const ProductVariant({
    required this.id,
    required this.name,
    required this.value,
    this.priceModifier = 0,
    this.stockQuantity = 0,
  });

  factory ProductVariant.fromJson(Map<String, dynamic> json) =>
      _$ProductVariantFromJson(json);

  Map<String, dynamic> toJson() => _$ProductVariantToJson(this);

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ProductVariant &&
          runtimeType == other.runtimeType &&
          id == other.id;

  @override
  int get hashCode => id.hashCode;
}
