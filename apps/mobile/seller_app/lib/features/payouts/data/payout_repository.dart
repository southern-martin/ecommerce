import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Represents the seller's current payout balance.
class PayoutBalance {
  final double available;
  final double pending;
  final double totalEarned;

  const PayoutBalance({
    required this.available,
    required this.pending,
    required this.totalEarned,
  });

  factory PayoutBalance.fromJson(Map<String, dynamic> json) {
    return PayoutBalance(
      available: (json['available'] as num).toDouble(),
      pending: (json['pending'] as num).toDouble(),
      totalEarned: (json['totalEarned'] as num).toDouble(),
    );
  }
}

/// Represents a single payout record.
class Payout {
  final String id;
  final double amount;
  final String status; // pending, processing, completed, failed
  final String method; // bank_transfer, paypal, stripe
  final DateTime createdAt;
  final DateTime? completedAt;
  final String? reference;

  const Payout({
    required this.id,
    required this.amount,
    required this.status,
    required this.method,
    required this.createdAt,
    this.completedAt,
    this.reference,
  });

  factory Payout.fromJson(Map<String, dynamic> json) {
    return Payout(
      id: json['id'] as String,
      amount: (json['amount'] as num).toDouble(),
      status: json['status'] as String,
      method: json['method'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
      completedAt: json['completedAt'] != null
          ? DateTime.parse(json['completedAt'] as String)
          : null,
      reference: json['reference'] as String?,
    );
  }

  String get methodDisplay {
    switch (method) {
      case 'bank_transfer':
        return 'Bank Transfer';
      case 'paypal':
        return 'PayPal';
      case 'stripe':
        return 'Stripe';
      default:
        return method;
    }
  }

  String get statusDisplay {
    return status[0].toUpperCase() + status.substring(1);
  }
}

/// Repository for managing seller payouts.
class PayoutRepository {
  final ApiClient _apiClient;

  PayoutRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Fetches the payout history for the seller.
  Future<List<Payout>> getPayoutHistory() async {
    final response = await _apiClient.get(ApiEndpoints.sellerPayouts);
    final list = response.data as List<dynamic>;
    return list
        .map((e) => Payout.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// Fetches the current payout balance.
  Future<PayoutBalance> getCurrentBalance() async {
    final response =
        await _apiClient.get('${ApiEndpoints.sellerPayouts}/balance');
    return PayoutBalance.fromJson(response.data as Map<String, dynamic>);
  }

  /// Requests a payout of the given [amount] via the specified [method].
  Future<Payout> requestPayout(double amount, String method) async {
    final response = await _apiClient.post(
      '${ApiEndpoints.sellerPayouts}/request',
      data: {
        'amount': amount,
        'method': method,
      },
    );
    return Payout.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches a single payout by ID.
  Future<Payout> getPayoutById(String id) async {
    final response =
        await _apiClient.get('${ApiEndpoints.sellerPayouts}/$id');
    return Payout.fromJson(response.data as Map<String, dynamic>);
  }
}
