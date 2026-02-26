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
}

/// Repository for managing seller coupons.
class CouponRepository {
  /// Fetches all coupons for the seller.
  Future<List<Coupon>> getCoupons() async {
    // TODO: Replace with actual API call to /seller/coupons
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return [
      Coupon(
        id: 'coup_1',
        code: 'SUMMER25',
        discountType: 'percentage',
        discountValue: 25,
        minimumOrderAmount: 50.0,
        maxUses: 100,
        usedCount: 43,
        expiryDate: now.add(const Duration(days: 30)),
        isActive: true,
        createdAt: now.subtract(const Duration(days: 15)),
        updatedAt: now.subtract(const Duration(days: 2)),
      ),
      Coupon(
        id: 'coup_2',
        code: 'FLAT10',
        discountType: 'fixed',
        discountValue: 10.0,
        minimumOrderAmount: 30.0,
        maxUses: 200,
        usedCount: 128,
        expiryDate: now.add(const Duration(days: 60)),
        isActive: true,
        createdAt: now.subtract(const Duration(days: 30)),
        updatedAt: now.subtract(const Duration(days: 5)),
      ),
      Coupon(
        id: 'coup_3',
        code: 'WELCOME15',
        discountType: 'percentage',
        discountValue: 15,
        maxUses: 50,
        usedCount: 50,
        expiryDate: now.subtract(const Duration(days: 5)),
        isActive: false,
        createdAt: now.subtract(const Duration(days: 60)),
        updatedAt: now.subtract(const Duration(days: 5)),
      ),
      Coupon(
        id: 'coup_4',
        code: 'FREESHIP',
        discountType: 'fixed',
        discountValue: 5.99,
        maxUses: null,
        usedCount: 312,
        expiryDate: null,
        isActive: true,
        createdAt: now.subtract(const Duration(days: 90)),
        updatedAt: now.subtract(const Duration(days: 1)),
      ),
      Coupon(
        id: 'coup_5',
        code: 'HOLIDAY50',
        discountType: 'percentage',
        discountValue: 50,
        minimumOrderAmount: 100.0,
        maxUses: 20,
        usedCount: 20,
        expiryDate: now.subtract(const Duration(days: 30)),
        isActive: false,
        createdAt: now.subtract(const Duration(days: 120)),
        updatedAt: now.subtract(const Duration(days: 30)),
      ),
    ];
  }

  /// Creates a new coupon.
  Future<Coupon> createCoupon(Map<String, dynamic> data) async {
    // TODO: Replace with actual API call to /seller/coupons
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return Coupon(
      id: 'coup_new_${now.millisecondsSinceEpoch}',
      code: data['code'] as String,
      discountType: data['discountType'] as String,
      discountValue: data['discountValue'] as double,
      minimumOrderAmount: data['minimumOrderAmount'] as double?,
      maxUses: data['maxUses'] as int?,
      usedCount: 0,
      expiryDate: data['expiryDate'] as DateTime?,
      productScope: data['productScope'] as String?,
      categoryScope: data['categoryScope'] as String?,
      isActive: true,
      createdAt: now,
      updatedAt: now,
    );
  }

  /// Updates an existing coupon.
  Future<Coupon> updateCoupon(String id, Map<String, dynamic> data) async {
    // TODO: Replace with actual API call to /seller/coupons/:id
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return Coupon(
      id: id,
      code: data['code'] as String,
      discountType: data['discountType'] as String,
      discountValue: data['discountValue'] as double,
      minimumOrderAmount: data['minimumOrderAmount'] as double?,
      maxUses: data['maxUses'] as int?,
      usedCount: 0,
      expiryDate: data['expiryDate'] as DateTime?,
      productScope: data['productScope'] as String?,
      categoryScope: data['categoryScope'] as String?,
      isActive: true,
      createdAt: now.subtract(const Duration(days: 10)),
      updatedAt: now,
    );
  }

  /// Deletes a coupon by ID.
  Future<void> deleteCoupon(String id) async {
    // TODO: Replace with actual API call to /seller/coupons/:id
    await Future.delayed(const Duration(milliseconds: 500));
  }

  /// Fetches coupon statistics.
  Future<CouponStats> getCouponStats() async {
    // TODO: Replace with actual API call to /seller/coupons/stats
    await Future.delayed(const Duration(milliseconds: 500));

    return const CouponStats(
      totalCoupons: 5,
      activeCoupons: 3,
      expiredCoupons: 2,
      totalRedemptions: 553,
      totalDiscountGiven: 4250.75,
    );
  }
}
