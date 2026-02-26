import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class BannerItem {
  final String id;
  final String imageUrl;
  final String title;
  final String? subtitle;
  final String? actionUrl;

  const BannerItem({
    required this.id,
    required this.imageUrl,
    required this.title,
    this.subtitle,
    this.actionUrl,
  });

  factory BannerItem.fromJson(Map<String, dynamic> json) {
    return BannerItem(
      id: json['id'] as String,
      imageUrl: json['imageUrl'] as String,
      title: json['title'] as String,
      subtitle: json['subtitle'] as String?,
      actionUrl: json['actionUrl'] as String?,
    );
  }
}

class Category {
  final String id;
  final String name;
  final String? iconUrl;
  final String slug;

  const Category({
    required this.id,
    required this.name,
    this.iconUrl,
    required this.slug,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'] as String,
      name: json['name'] as String,
      iconUrl: json['iconUrl'] as String?,
      slug: json['slug'] as String,
    );
  }
}

class ProductSummary {
  final String id;
  final String name;
  final String slug;
  final double price;
  final double? compareAtPrice;
  final String imageUrl;
  final double rating;
  final int reviewCount;

  const ProductSummary({
    required this.id,
    required this.name,
    required this.slug,
    required this.price,
    this.compareAtPrice,
    required this.imageUrl,
    required this.rating,
    required this.reviewCount,
  });

  bool get hasDiscount => compareAtPrice != null && compareAtPrice! > price;

  double get discountPercentage {
    if (!hasDiscount) return 0;
    return ((compareAtPrice! - price) / compareAtPrice! * 100).roundToDouble();
  }

  factory ProductSummary.fromJson(Map<String, dynamic> json) {
    return ProductSummary(
      id: json['id'] as String,
      name: json['name'] as String,
      slug: json['slug'] as String,
      price: (json['price'] as num).toDouble(),
      compareAtPrice: json['compareAtPrice'] != null
          ? (json['compareAtPrice'] as num).toDouble()
          : null,
      imageUrl: json['imageUrl'] as String,
      rating: (json['rating'] as num? ?? 0).toDouble(),
      reviewCount: json['reviewCount'] as int? ?? 0,
    );
  }
}

class HomeRepository {
  final ApiClient _apiClient;

  HomeRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<List<BannerItem>> getBanners() async {
    final response = await _apiClient.get('/home/banners');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => BannerItem.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<List<ProductSummary>> getFeaturedProducts() async {
    final response = await _apiClient.get('/home/featured-products');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => ProductSummary.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<List<ProductSummary>> getTrendingProducts() async {
    final response = await _apiClient.get('/home/trending-products');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => ProductSummary.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<List<Category>> getCategories() async {
    final response = await _apiClient.get('/categories');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => Category.fromJson(e as Map<String, dynamic>)).toList();
  }
}
