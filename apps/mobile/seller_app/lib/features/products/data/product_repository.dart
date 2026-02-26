/// Represents a seller's product.
class SellerProduct {
  final String id;
  final String name;
  final String description;
  final double price;
  final double? compareAtPrice;
  final String categoryId;
  final String categoryName;
  final int stockQuantity;
  final String status; // active, draft, out_of_stock
  final List<String> imageUrls;
  final List<ProductVariant> variants;
  final DateTime createdAt;
  final DateTime updatedAt;

  const SellerProduct({
    required this.id,
    required this.name,
    required this.description,
    required this.price,
    this.compareAtPrice,
    required this.categoryId,
    required this.categoryName,
    required this.stockQuantity,
    required this.status,
    required this.imageUrls,
    required this.variants,
    required this.createdAt,
    required this.updatedAt,
  });
}

/// A variant of a product (e.g. size, color).
class ProductVariant {
  final String id;
  final String name;
  final String value;
  final double priceModifier;
  final int stock;

  const ProductVariant({
    required this.id,
    required this.name,
    required this.value,
    required this.priceModifier,
    required this.stock,
  });
}

/// Paginated result wrapper.
class PaginatedProducts {
  final List<SellerProduct> products;
  final int totalCount;
  final int currentPage;
  final int totalPages;

  const PaginatedProducts({
    required this.products,
    required this.totalCount,
    required this.currentPage,
    required this.totalPages,
  });
}

/// Repository for managing seller products.
class SellerProductRepository {
  /// Fetches paginated list of seller products with optional status filter.
  Future<PaginatedProducts> getMyProducts({
    int page = 1,
    String? status,
    String? searchQuery,
  }) async {
    // TODO: Replace with actual API call
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    final allProducts = List.generate(
      15,
      (i) => SellerProduct(
        id: 'prod_${i + 1}',
        name: 'Product ${i + 1}',
        description: 'Description for product ${i + 1}',
        price: 19.99 + (i * 10),
        compareAtPrice: i % 3 == 0 ? 29.99 + (i * 10) : null,
        categoryId: 'cat_${(i % 5) + 1}',
        categoryName: ['Electronics', 'Clothing', 'Home', 'Sports', 'Books'][i % 5],
        stockQuantity: i % 4 == 0 ? 0 : 10 + i * 5,
        status: i % 4 == 0 ? 'out_of_stock' : (i % 3 == 0 ? 'draft' : 'active'),
        imageUrls: ['https://picsum.photos/200?random=$i'],
        variants: i % 2 == 0
            ? [
                ProductVariant(
                  id: 'var_${i}_1',
                  name: 'Size',
                  value: 'Large',
                  priceModifier: 5.0,
                  stock: 10,
                ),
              ]
            : [],
        createdAt: now.subtract(Duration(days: i * 3)),
        updatedAt: now.subtract(Duration(days: i)),
      ),
    );

    final filtered = status != null
        ? allProducts.where((p) => p.status == status).toList()
        : allProducts;

    return PaginatedProducts(
      products: filtered,
      totalCount: filtered.length,
      currentPage: page,
      totalPages: 1,
    );
  }

  /// Creates a new product.
  Future<SellerProduct> createProduct({
    required String name,
    required String description,
    required double price,
    double? compareAtPrice,
    required String categoryId,
    required int stockQuantity,
    required List<String> imageUrls,
    required List<ProductVariant> variants,
  }) async {
    await Future.delayed(const Duration(seconds: 1));
    return SellerProduct(
      id: 'prod_new_${DateTime.now().millisecondsSinceEpoch}',
      name: name,
      description: description,
      price: price,
      compareAtPrice: compareAtPrice,
      categoryId: categoryId,
      categoryName: 'Category',
      stockQuantity: stockQuantity,
      status: 'draft',
      imageUrls: imageUrls,
      variants: variants,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    );
  }

  /// Updates an existing product by ID.
  Future<SellerProduct> updateProduct({
    required String id,
    required String name,
    required String description,
    required double price,
    double? compareAtPrice,
    required String categoryId,
    required int stockQuantity,
    required List<String> imageUrls,
    required List<ProductVariant> variants,
  }) async {
    await Future.delayed(const Duration(seconds: 1));
    return SellerProduct(
      id: id,
      name: name,
      description: description,
      price: price,
      compareAtPrice: compareAtPrice,
      categoryId: categoryId,
      categoryName: 'Category',
      stockQuantity: stockQuantity,
      status: 'active',
      imageUrls: imageUrls,
      variants: variants,
      createdAt: DateTime.now().subtract(const Duration(days: 30)),
      updatedAt: DateTime.now(),
    );
  }

  /// Deletes a product by ID.
  Future<void> deleteProduct(String id) async {
    await Future.delayed(const Duration(milliseconds: 500));
  }

  /// Uploads product images and returns their URLs.
  Future<List<String>> uploadImages(List<String> localPaths) async {
    await Future.delayed(const Duration(seconds: 2));
    return localPaths
        .map((path) =>
            'https://picsum.photos/400?random=${DateTime.now().millisecondsSinceEpoch}')
        .toList();
  }
}
