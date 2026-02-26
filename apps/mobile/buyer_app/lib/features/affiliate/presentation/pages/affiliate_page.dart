import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../../../core/di/injection.dart';
import '../../data/affiliate_repository.dart';

class AffiliatePage extends StatefulWidget {
  const AffiliatePage({super.key});

  @override
  State<AffiliatePage> createState() => _AffiliatePageState();
}

class _AffiliatePageState extends State<AffiliatePage> {
  final AffiliateRepository _affiliateRepo = getIt<AffiliateRepository>();

  AffiliateStats? _stats;
  bool _isLoading = true;
  bool _isRequestingPayout = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadStats();
  }

  Future<void> _loadStats() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final stats = await _affiliateRepo.getStats();
      if (mounted) {
        setState(() {
          _stats = stats;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _error = e.toString();
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _copyReferralLink() async {
    if (_stats == null) return;
    await Clipboard.setData(ClipboardData(text: _stats!.referralLink));
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Referral link copied to clipboard')),
      );
    }
  }

  Future<void> _requestPayout() async {
    if (_stats == null || _stats!.pendingPayout <= 0) return;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Request Payout'),
        content: Text(
          'Request a payout of \$${_stats!.pendingPayout.toStringAsFixed(2)}?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Request'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      setState(() => _isRequestingPayout = true);
      try {
        await _affiliateRepo.requestPayout();
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Payout requested successfully')),
          );
          _loadStats();
        }
      } catch (e) {
        if (mounted) {
          setState(() => _isRequestingPayout = false);
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to request payout: $e')),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Affiliate Program'),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(Icons.error_outline, size: 64, color: Colors.grey.shade400),
                      const SizedBox(height: 16),
                      Text('Failed to load affiliate data',
                          style: theme.textTheme.titleMedium),
                      const SizedBox(height: 8),
                      FilledButton(onPressed: _loadStats, child: const Text('Retry')),
                    ],
                  ),
                )
              : RefreshIndicator(
                  onRefresh: _loadStats,
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      // Stats cards grid
                      GridView.count(
                        crossAxisCount: 2,
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        mainAxisSpacing: 12,
                        crossAxisSpacing: 12,
                        childAspectRatio: 1.4,
                        children: [
                          _buildStatCard(
                            theme,
                            icon: Icons.mouse,
                            label: 'Total Clicks',
                            value: '${_stats!.totalClicks}',
                            color: Colors.blue,
                          ),
                          _buildStatCard(
                            theme,
                            icon: Icons.shopping_bag,
                            label: 'Conversions',
                            value: '${_stats!.totalConversions}',
                            color: Colors.green,
                          ),
                          _buildStatCard(
                            theme,
                            icon: Icons.attach_money,
                            label: 'Total Earnings',
                            value: '\$${_stats!.totalEarnings.toStringAsFixed(2)}',
                            color: Colors.orange,
                          ),
                          _buildStatCard(
                            theme,
                            icon: Icons.account_balance_wallet,
                            label: 'Pending Payout',
                            value: '\$${_stats!.pendingPayout.toStringAsFixed(2)}',
                            color: Colors.purple,
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Referral link
                      Text(
                        'Your Referral Link',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Card(
                        child: Padding(
                          padding: const EdgeInsets.all(16),
                          child: Row(
                            children: [
                              Expanded(
                                child: Text(
                                  _stats!.referralLink,
                                  style: theme.textTheme.bodyMedium?.copyWith(
                                    color: theme.colorScheme.primary,
                                  ),
                                  overflow: TextOverflow.ellipsis,
                                  maxLines: 2,
                                ),
                              ),
                              const SizedBox(width: 12),
                              IconButton(
                                onPressed: _copyReferralLink,
                                icon: const Icon(Icons.copy),
                                tooltip: 'Copy link',
                              ),
                            ],
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Referral Code: ${_stats!.referralCode}',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: Colors.grey,
                        ),
                      ),
                      const SizedBox(height: 24),

                      // Payout button
                      SizedBox(
                        width: double.infinity,
                        height: 52,
                        child: ElevatedButton.icon(
                          onPressed: _stats!.pendingPayout > 0 && !_isRequestingPayout
                              ? _requestPayout
                              : null,
                          icon: _isRequestingPayout
                              ? const SizedBox(
                                  width: 20,
                                  height: 20,
                                  child: CircularProgressIndicator(strokeWidth: 2),
                                )
                              : const Icon(Icons.payments_outlined),
                          label: Text(
                            _stats!.pendingPayout > 0
                                ? 'Request Payout (\$${_stats!.pendingPayout.toStringAsFixed(2)})'
                                : 'No Pending Payout',
                          ),
                          style: ElevatedButton.styleFrom(
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
    );
  }

  Widget _buildStatCard(
    ThemeData theme, {
    required IconData icon,
    required String label,
    required String value,
    required Color color,
  }) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Icon(icon, color: color, size: 28),
            const Spacer(),
            Text(
              value,
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 2),
            Text(
              label,
              style: theme.textTheme.bodySmall?.copyWith(
                color: Colors.grey,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
