import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class LoyaltyTier {
  final String id;
  final String name;
  final int minPoints;
  final double multiplier;
  final List<String> benefits;

  const LoyaltyTier({
    required this.id,
    required this.name,
    required this.minPoints,
    required this.multiplier,
    required this.benefits,
  });

  factory LoyaltyTier.fromJson(Map<String, dynamic> json) {
    return LoyaltyTier(
      id: json['id'] as String,
      name: json['name'] as String,
      minPoints: json['minPoints'] as int,
      multiplier: (json['multiplier'] as num).toDouble(),
      benefits: (json['benefits'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
    );
  }
}

class LoyaltyTransaction {
  final String id;
  final String type;
  final int points;
  final String description;
  final DateTime createdAt;

  const LoyaltyTransaction({
    required this.id,
    required this.type,
    required this.points,
    required this.description,
    required this.createdAt,
  });

  factory LoyaltyTransaction.fromJson(Map<String, dynamic> json) {
    return LoyaltyTransaction(
      id: json['id'] as String,
      type: json['type'] as String,
      points: json['points'] as int,
      description: json['description'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
    );
  }
}

class LoyaltyMembership {
  final String tierId;
  final String tierName;
  final int currentPoints;
  final int lifetimePoints;
  final int pointsToNextTier;
  final String? nextTierName;

  const LoyaltyMembership({
    required this.tierId,
    required this.tierName,
    required this.currentPoints,
    required this.lifetimePoints,
    required this.pointsToNextTier,
    this.nextTierName,
  });

  factory LoyaltyMembership.fromJson(Map<String, dynamic> json) {
    return LoyaltyMembership(
      tierId: json['tierId'] as String,
      tierName: json['tierName'] as String,
      currentPoints: json['currentPoints'] as int,
      lifetimePoints: json['lifetimePoints'] as int,
      pointsToNextTier: json['pointsToNextTier'] as int? ?? 0,
      nextTierName: json['nextTierName'] as String?,
    );
  }
}

class LoyaltyRepository {
  final ApiClient _apiClient;

  LoyaltyRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<LoyaltyMembership> getMembership() async {
    final response = await _apiClient.get('/loyalty/membership');
    return LoyaltyMembership.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<LoyaltyTransaction>> getTransactions({int page = 1}) async {
    final response = await _apiClient.get(
      '/loyalty/transactions',
      queryParameters: {'page': page},
    );
    final List<dynamic> data = response.data is List
        ? response.data as List<dynamic>
        : (response.data as Map<String, dynamic>)['data'] as List<dynamic>;
    return data
        .map((e) => LoyaltyTransaction.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<LoyaltyTier>> getTiers() async {
    final response = await _apiClient.get('/loyalty/tiers');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => LoyaltyTier.fromJson(e as Map<String, dynamic>)).toList();
  }
}
