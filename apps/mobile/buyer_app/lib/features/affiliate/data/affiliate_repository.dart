import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class AffiliateStats {
  final int totalClicks;
  final int totalConversions;
  final double totalEarnings;
  final double pendingPayout;
  final String referralCode;
  final String referralLink;

  const AffiliateStats({
    required this.totalClicks,
    required this.totalConversions,
    required this.totalEarnings,
    required this.pendingPayout,
    required this.referralCode,
    required this.referralLink,
  });

  factory AffiliateStats.fromJson(Map<String, dynamic> json) {
    return AffiliateStats(
      totalClicks: json['totalClicks'] as int? ?? 0,
      totalConversions: json['totalConversions'] as int? ?? 0,
      totalEarnings: (json['totalEarnings'] as num? ?? 0).toDouble(),
      pendingPayout: (json['pendingPayout'] as num? ?? 0).toDouble(),
      referralCode: json['referralCode'] as String,
      referralLink: json['referralLink'] as String,
    );
  }
}

class AffiliateClick {
  final String id;
  final DateTime timestamp;
  final String? referrer;
  final bool converted;

  const AffiliateClick({
    required this.id,
    required this.timestamp,
    this.referrer,
    required this.converted,
  });

  factory AffiliateClick.fromJson(Map<String, dynamic> json) {
    return AffiliateClick(
      id: json['id'] as String,
      timestamp: DateTime.parse(json['timestamp'] as String),
      referrer: json['referrer'] as String?,
      converted: json['converted'] as bool? ?? false,
    );
  }
}

class AffiliateRepository {
  final ApiClient _apiClient;

  AffiliateRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<AffiliateStats> getStats() async {
    final response = await _apiClient.get('/affiliate/stats');
    return AffiliateStats.fromJson(response.data as Map<String, dynamic>);
  }

  Future<String> generateLink({String? campaignName}) async {
    final response = await _apiClient.post('/affiliate/links', data: {
      if (campaignName != null) 'campaignName': campaignName,
    });
    return (response.data as Map<String, dynamic>)['link'] as String;
  }

  Future<List<AffiliateClick>> getClicks({int page = 1}) async {
    final response = await _apiClient.get(
      '/affiliate/clicks',
      queryParameters: {'page': page},
    );
    final List<dynamic> data = response.data is List
        ? response.data as List<dynamic>
        : (response.data as Map<String, dynamic>)['data'] as List<dynamic>;
    return data
        .map((e) => AffiliateClick.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> requestPayout() async {
    await _apiClient.post('/affiliate/payout');
  }
}
