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
  /// Fetches the payout history for the seller.
  Future<List<Payout>> getPayoutHistory() async {
    // TODO: Replace with actual API call to /seller/payouts
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return [
      Payout(
        id: 'pay_1',
        amount: 1250.00,
        status: 'completed',
        method: 'bank_transfer',
        createdAt: now.subtract(const Duration(days: 3)),
        completedAt: now.subtract(const Duration(days: 1)),
        reference: 'TXN-20260223-001',
      ),
      Payout(
        id: 'pay_2',
        amount: 890.50,
        status: 'processing',
        method: 'paypal',
        createdAt: now.subtract(const Duration(days: 1)),
      ),
      Payout(
        id: 'pay_3',
        amount: 2100.00,
        status: 'completed',
        method: 'bank_transfer',
        createdAt: now.subtract(const Duration(days: 10)),
        completedAt: now.subtract(const Duration(days: 8)),
        reference: 'TXN-20260216-003',
      ),
      Payout(
        id: 'pay_4',
        amount: 450.00,
        status: 'failed',
        method: 'stripe',
        createdAt: now.subtract(const Duration(days: 15)),
      ),
      Payout(
        id: 'pay_5',
        amount: 3200.00,
        status: 'completed',
        method: 'bank_transfer',
        createdAt: now.subtract(const Duration(days: 20)),
        completedAt: now.subtract(const Duration(days: 18)),
        reference: 'TXN-20260206-002',
      ),
      Payout(
        id: 'pay_6',
        amount: 675.25,
        status: 'pending',
        method: 'paypal',
        createdAt: now.subtract(const Duration(hours: 6)),
      ),
      Payout(
        id: 'pay_7',
        amount: 1500.00,
        status: 'completed',
        method: 'stripe',
        createdAt: now.subtract(const Duration(days: 30)),
        completedAt: now.subtract(const Duration(days: 28)),
        reference: 'TXN-20260127-001',
      ),
    ];
  }

  /// Fetches the current payout balance.
  Future<PayoutBalance> getCurrentBalance() async {
    // TODO: Replace with actual API call to /seller/payouts/balance
    await Future.delayed(const Duration(milliseconds: 800));

    return const PayoutBalance(
      available: 4325.75,
      pending: 890.50,
      totalEarned: 34250.00,
    );
  }

  /// Requests a payout of the given [amount] via the specified [method].
  Future<Payout> requestPayout(double amount, String method) async {
    // TODO: Replace with actual API call to /seller/payouts/request
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return Payout(
      id: 'pay_new_${now.millisecondsSinceEpoch}',
      amount: amount,
      status: 'pending',
      method: method,
      createdAt: now,
    );
  }

  /// Fetches a single payout by ID.
  Future<Payout> getPayoutById(String id) async {
    // TODO: Replace with actual API call to /seller/payouts/:id
    await Future.delayed(const Duration(milliseconds: 800));

    final now = DateTime.now();
    return Payout(
      id: id,
      amount: 1250.00,
      status: 'completed',
      method: 'bank_transfer',
      createdAt: now.subtract(const Duration(days: 3)),
      completedAt: now.subtract(const Duration(days: 1)),
      reference: 'TXN-20260223-001',
    );
  }
}
