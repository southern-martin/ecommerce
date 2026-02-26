import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import '../../home/data/home_repository.dart';

class ProductVariant {
  final String id;
  final String name;
  final String value;
  final double? priceAdjustment;
  final bool inStock;

  const ProductVariant({
    required this.id,
    required this.name,
    required this.value,
    this.priceAdjustment,
    required this.inStock,
  });

  factory ProductVariant.fromJson(Map<String, dynamic> json) {
    return ProductVariant(
      id: json['id'] as String,
      name: json['name'] as String,
      value: json['value'] as String,
      priceAdjustment: json['priceAdjustment'] != null
          ? (json['priceAdjustment'] as num).toDouble()
          : null,
      inStock: json['inStock'] as bool? ?? true,
    );
  }
}

class SellerInfo {
  final String id;
  final String name;
  final double rating;
  final String? avatarUrl;

  const SellerInfo({
    required this.id,
    required this.name,
    required this.rating,
    this.avatarUrl,
  });

  factory SellerInfo.fromJson(Map<String, dynamic> json) {
    return SellerInfo(
      id: json['id'] as String,
      name: json['name'] as String,
      rating: (json['rating'] as num).toDouble(),
      avatarUrl: json['avatarUrl'] as String?,
    );
  }
}

class ProductDetail {
  final String id;
  final String name;
  final String slug;
  final String description;
  final double price;
  final double? compareAtPrice;
  final List<String> images;
  final double rating;
  final int reviewCount;
  final bool inStock;
  final int stockQuantity;
  final List<ProductVariant> variants;
  final SellerInfo seller;
  final String category;
  final List<String> tags;

  const ProductDetail({
    required this.id,
    required this.name,
    required this.slug,
    required this.description,
    required this.price,
    this.compareAtPrice,
    required this.images,
    required this.rating,
    required this.reviewCount,
    required this.inStock,
    required this.stockQuantity,
    required this.variants,
    required this.seller,
    required this.category,
    required this.tags,
  });

  bool get hasDiscount => compareAtPrice != null && compareAtPrice! > price;

  factory ProductDetail.fromJson(Map<String, dynamic> json) {
    return ProductDetail(
      id: json['id'] as String,
      name: json['name'] as String,
      slug: json['slug'] as String,
      description: json['description'] as String,
      price: (json['price'] as num).toDouble(),
      compareAtPrice: json['compareAtPrice'] != null
          ? (json['compareAtPrice'] as num).toDouble()
          : null,
      images: (json['images'] as List<dynamic>).map((e) => e as String).toList(),
      rating: (json['rating'] as num? ?? 0).toDouble(),
      reviewCount: json['reviewCount'] as int? ?? 0,
      inStock: json['inStock'] as bool? ?? true,
      stockQuantity: json['stockQuantity'] as int? ?? 0,
      variants: (json['variants'] as List<dynamic>?)
              ?.map((e) => ProductVariant.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      seller: SellerInfo.fromJson(json['seller'] as Map<String, dynamic>),
      category: json['category'] as String? ?? '',
      tags: (json['tags'] as List<dynamic>?)?.map((e) => e as String).toList() ?? [],
    );
  }
}

class PaginatedProducts {
  final List<ProductSummary> products;
  final int total;
  final int page;
  final int pageSize;

  const PaginatedProducts({
    required this.products,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  bool get hasMore => page * pageSize < total;

  factory PaginatedProducts.fromJson(Map<String, dynamic> json) {
    return PaginatedProducts(
      products: (json['data'] as List<dynamic>)
          .map((e) => ProductSummary.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      page: json['page'] as int,
      pageSize: json['pageSize'] as int,
    );
  }
}

class ProductRepository {
  final ApiClient _apiClient;

  ProductRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<PaginatedProducts> getProducts({
    String? category,
    double? minPrice,
    double? maxPrice,
    double? rating,
    bool? inStock,
    int page = 1,
    int pageSize = 20,
    String? sort,
    String? search,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'pageSize': pageSize,
    };

    if (category != null) queryParams['category'] = category;
    if (minPrice != null) queryParams['minPrice'] = minPrice;
    if (maxPrice != null) queryParams['maxPrice'] = maxPrice;
    if (rating != null) queryParams['rating'] = rating;
    if (inStock != null) queryParams['inStock'] = inStock;
    if (sort != null) queryParams['sort'] = sort;
    if (search != null) queryParams['search'] = search;

    final response = await _apiClient.get('/products', queryParameters: queryParams);
    return PaginatedProducts.fromJson(response.data as Map<String, dynamic>);
  }

  Future<ProductDetail> getProductBySlug(String slug) async {
    final response = await _apiClient.get('/products/$slug');
    return ProductDetail.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<Category>> getCategories() async {
    final response = await _apiClient.get('/categories');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => Category.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<List<ProductSummary>> getFeaturedProducts() async {
    final response = await _apiClient.get('/products/featured');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => ProductSummary.fromJson(e as Map<String, dynamic>)).toList();
  }
}
