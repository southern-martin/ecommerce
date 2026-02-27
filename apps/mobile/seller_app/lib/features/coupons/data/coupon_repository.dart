import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Represents a coupon created by the seller.
class Coupon {
  final String id;
  final String code;
  final String discountType; // percentage, fixed
  final double discountValue;
  final double? minimumOrderAmount;
  final int? maxUses;
  final int usedCount;
  final DateTime? expiryDate;
  final String? productScope;
  final String? categoryScope;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Coupon({
    required this.id,
    required this.code,
    required this.discountType,
    required this.discountValue,
    this.minimumOrderAmount,
    this.maxUses,
    required this.usedCount,
    this.expiryDate,
    this.productScope,
    this.categoryScope,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Coupon.fromJson(Map<String, dynamic> json) {
    return Coupon(
      id: json['id'] as String,
      code: json['code'] as String,
      discountType: json['discountType'] as String,
      discountValue: (json['discountValue'] as num).toDouble(),
      minimumOrderAmount: json['minimumOrderAmount'] != null
          ? (json['minimumOrderAmount'] as num).toDouble()
          : null,
      maxUses: json['maxUses'] as int?,
      usedCount: json['usedCount'] as int,
      expiryDate: json['expiryDate'] != null
          ? DateTime.parse(json['expiryDate'] as String)
          : null,
      productScope: json['productScope'] as String?,
      categoryScope: json['categoryScope'] as String?,
      isActive: json['isActive'] as bool,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  String get discountDisplay {
    if (discountType == 'percentage') {
      return '${discountValue.toStringAsFixed(0)}%';
    }
    return '\$${discountValue.toStringAsFixed(2)}';
  }

  bool get isExpired {
    if (expiryDate == null) return false;
    return expiryDate!.isBefore(DateTime.now());
  }
}

/// Coupon statistics summary.
class CouponStats {
  final int totalCoupons;
  final int activeCoupons;
  final int expiredCoupons;
  final int totalRedemptions;
  final double totalDiscountGiven;

  const CouponStats({
    required this.totalCoupons,
    required this.activeCoupons,
    required this.expiredCoupons,
    required this.totalRedemptions,
    required this.totalDiscountGiven,
  });

  factory CouponStats.fromJson(Map<String, dynamic> json) {
    return CouponStats(
      totalCoupons: json['totalCoupons'] as int,
      activeCoupons: json['activeCoupons'] as int,
      expiredCoupons: json['expiredCoupons'] as int,
      totalRedemptions: json['totalRedemptions'] as int,
      totalDiscountGiven: (json['totalDiscountGiven'] as num).toDouble(),
    );
  }
}

/// Repository for managing seller coupons.
class CouponRepository {
  final ApiClient _apiClient;

  CouponRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Fetches all coupons for the seller.
  Future<List<Coupon>> getCoupons() async {
    final response = await _apiClient.get(ApiEndpoints.sellerCoupons);
    final list = response.data as List<dynamic>;
    return list
        .map((e) => Coupon.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// Creates a new coupon.
  Future<Coupon> createCoupon(Map<String, dynamic> data) async {
    final response =
        await _apiClient.post(ApiEndpoints.sellerCoupons, data: data);
    return Coupon.fromJson(response.data as Map<String, dynamic>);
  }

  /// Updates an existing coupon.
  Future<Coupon> updateCoupon(String id, Map<String, dynamic> data) async {
    final response =
        await _apiClient.put('${ApiEndpoints.sellerCoupons}/$id', data: data);
    return Coupon.fromJson(response.data as Map<String, dynamic>);
  }

  /// Deletes a coupon by ID.
  Future<void> deleteCoupon(String id) async {
    await _apiClient.delete('${ApiEndpoints.sellerCoupons}/$id');
  }

  /// Fetches coupon statistics.
  Future<CouponStats> getCouponStats() async {
    final response =
        await _apiClient.get('${ApiEndpoints.sellerCoupons}/stats');
    return CouponStats.fromJson(response.data as Map<String, dynamic>);
  }
}
