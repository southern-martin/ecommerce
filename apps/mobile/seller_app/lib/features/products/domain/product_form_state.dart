import '../data/product_repository.dart';

/// Represents the state of a product form used for both creating and editing.
class ProductFormState {
  final String name;
  final String description;
  final double price;
  final double? compareAtPrice;
  final String categoryId;
  final int stockQuantity;
  final List<String> images;
  final List<ProductVariantFormEntry> variants;

  const ProductFormState({
    this.name = '',
    this.description = '',
    this.price = 0.0,
    this.compareAtPrice,
    this.categoryId = '',
    this.stockQuantity = 0,
    this.images = const [],
    this.variants = const [],
  });

  /// Creates a ProductFormState from an existing product for editing.
  factory ProductFormState.fromProduct(SellerProduct product) {
    return ProductFormState(
      name: product.name,
      description: product.description,
      price: product.price,
      compareAtPrice: product.compareAtPrice,
      categoryId: product.categoryId,
      stockQuantity: product.stockQuantity,
      images: product.imageUrls,
      variants: product.variants
          .map((v) => ProductVariantFormEntry(
                name: v.name,
                value: v.value,
                priceModifier: v.priceModifier,
                stock: v.stock,
              ))
          .toList(),
    );
  }

  ProductFormState copyWith({
    String? name,
    String? description,
    double? price,
    double? Function()? compareAtPrice,
    String? categoryId,
    int? stockQuantity,
    List<String>? images,
    List<ProductVariantFormEntry>? variants,
  }) {
    return ProductFormState(
      name: name ?? this.name,
      description: description ?? this.description,
      price: price ?? this.price,
      compareAtPrice:
          compareAtPrice != null ? compareAtPrice() : this.compareAtPrice,
      categoryId: categoryId ?? this.categoryId,
      stockQuantity: stockQuantity ?? this.stockQuantity,
      images: images ?? this.images,
      variants: variants ?? this.variants,
    );
  }

  /// Validates the form and returns a list of error messages.
  /// Returns an empty list if the form is valid.
  List<String> validate() {
    final errors = <String>[];

    if (name.trim().isEmpty) {
      errors.add('Product name is required');
    }

    if (description.trim().isEmpty) {
      errors.add('Product description is required');
    }

    if (price <= 0) {
      errors.add('Price must be greater than zero');
    }

    if (compareAtPrice != null && compareAtPrice! <= price) {
      errors.add('Compare at price must be greater than the selling price');
    }

    if (categoryId.isEmpty) {
      errors.add('Please select a category');
    }

    if (stockQuantity < 0) {
      errors.add('Stock quantity cannot be negative');
    }

    return errors;
  }

  /// Returns `true` if the form state is valid.
  bool get isValid => validate().isEmpty;
}

/// A single variant entry in the product form.
class ProductVariantFormEntry {
  final String name;
  final String value;
  final double priceModifier;
  final int stock;

  const ProductVariantFormEntry({
    this.name = '',
    this.value = '',
    this.priceModifier = 0.0,
    this.stock = 0,
  });

  ProductVariantFormEntry copyWith({
    String? name,
    String? value,
    double? priceModifier,
    int? stock,
  }) {
    return ProductVariantFormEntry(
      name: name ?? this.name,
      value: value ?? this.value,
      priceModifier: priceModifier ?? this.priceModifier,
      stock: stock ?? this.stock,
    );
  }
}
