import 'package:dio/dio.dart';
import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

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

  factory SellerProduct.fromJson(Map<String, dynamic> json) {
    return SellerProduct(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String,
      price: (json['price'] as num).toDouble(),
      compareAtPrice: json['compare_at_price'] != null
          ? (json['compare_at_price'] as num).toDouble()
          : null,
      categoryId: json['category_id'] as String,
      categoryName: json['category_name'] as String,
      stockQuantity: json['stock_quantity'] as int,
      status: json['status'] as String,
      imageUrls: (json['image_urls'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      variants: (json['variants'] as List<dynamic>)
          .map((e) => ProductVariant.fromJson(e as Map<String, dynamic>))
          .toList(),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }
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

  factory ProductVariant.fromJson(Map<String, dynamic> json) {
    return ProductVariant(
      id: json['id'] as String,
      name: json['name'] as String,
      value: json['value'] as String,
      priceModifier: (json['price_modifier'] as num).toDouble(),
      stock: json['stock'] as int,
    );
  }
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

  factory PaginatedProducts.fromJson(Map<String, dynamic> json) {
    return PaginatedProducts(
      products: (json['products'] as List<dynamic>)
          .map((e) => SellerProduct.fromJson(e as Map<String, dynamic>))
          .toList(),
      totalCount: json['total_count'] as int,
      currentPage: json['current_page'] as int,
      totalPages: json['total_pages'] as int,
    );
  }
}

/// Repository for managing seller products.
class SellerProductRepository {
  final ApiClient _apiClient;

  SellerProductRepository({required ApiClient apiClient})
      : _apiClient = apiClient;

  /// Fetches paginated list of seller products with optional status filter.
  Future<PaginatedProducts> getMyProducts({
    int page = 1,
    String? status,
    String? searchQuery,
  }) async {
    final queryParams = <String, dynamic>{'page': page};
    if (status != null) queryParams['status'] = status;
    if (searchQuery != null && searchQuery.isNotEmpty) {
      queryParams['search'] = searchQuery;
    }

    final response = await _apiClient.get(
      ApiEndpoints.sellerProducts,
      queryParameters: queryParams,
    );
    return PaginatedProducts.fromJson(response.data as Map<String, dynamic>);
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
    final response = await _apiClient.post(
      ApiEndpoints.sellerProducts,
      data: {
        'name': name,
        'description': description,
        'price': price,
        if (compareAtPrice != null) 'compare_at_price': compareAtPrice,
        'category_id': categoryId,
        'stock_quantity': stockQuantity,
        'image_urls': imageUrls,
        'variants': variants
            .map((v) => {
                  'name': v.name,
                  'value': v.value,
                  'price_modifier': v.priceModifier,
                  'stock': v.stock,
                })
            .toList(),
      },
    );
    return SellerProduct.fromJson(response.data as Map<String, dynamic>);
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
    final response = await _apiClient.put(
      '${ApiEndpoints.sellerProducts}/$id',
      data: {
        'name': name,
        'description': description,
        'price': price,
        if (compareAtPrice != null) 'compare_at_price': compareAtPrice,
        'category_id': categoryId,
        'stock_quantity': stockQuantity,
        'image_urls': imageUrls,
        'variants': variants
            .map((v) => {
                  'id': v.id,
                  'name': v.name,
                  'value': v.value,
                  'price_modifier': v.priceModifier,
                  'stock': v.stock,
                })
            .toList(),
      },
    );
    return SellerProduct.fromJson(response.data as Map<String, dynamic>);
  }

  /// Deletes a product by ID.
  Future<void> deleteProduct(String id) async {
    await _apiClient.delete('${ApiEndpoints.sellerProducts}/$id');
  }

  /// Uploads product images and returns their URLs.
  Future<List<String>> uploadImages(List<String> localPaths) async {
    final formData = FormData.fromMap({
      'images': await Future.wait(
        localPaths.map(
          (path) => MultipartFile.fromFile(path),
        ),
      ),
    });

    final response = await _apiClient.upload(
      '${ApiEndpoints.sellerProducts}/images',
      formData: formData,
    );

    final data = response.data as Map<String, dynamic>;
    return (data['urls'] as List<dynamic>).map((e) => e as String).toList();
  }
}
