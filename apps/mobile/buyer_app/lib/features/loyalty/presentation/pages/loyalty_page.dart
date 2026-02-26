import 'package:flutter/material.dart';

import '../../../../core/di/injection.dart';
import '../../data/loyalty_repository.dart';

class LoyaltyPage extends StatefulWidget {
  const LoyaltyPage({super.key});

  @override
  State<LoyaltyPage> createState() => _LoyaltyPageState();
}

class _LoyaltyPageState extends State<LoyaltyPage> {
  final LoyaltyRepository _loyaltyRepo = getIt<LoyaltyRepository>();

  LoyaltyMembership? _membership;
  List<LoyaltyTier> _tiers = [];
  List<LoyaltyTransaction> _transactions = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final results = await Future.wait([
        _loyaltyRepo.getMembership(),
        _loyaltyRepo.getTiers(),
        _loyaltyRepo.getTransactions(),
      ]);

      if (mounted) {
        setState(() {
          _membership = results[0] as LoyaltyMembership;
          _tiers = results[1] as List<LoyaltyTier>;
          _transactions = results[2] as List<LoyaltyTransaction>;
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

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Loyalty Program'),
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
                      Text('Failed to load loyalty data', style: theme.textTheme.titleMedium),
                      const SizedBox(height: 8),
                      FilledButton(onPressed: _loadData, child: const Text('Retry')),
                    ],
                  ),
                )
              : RefreshIndicator(
                  onRefresh: _loadData,
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      // Tier card
                      Card(
                        color: theme.colorScheme.primary,
                        child: Padding(
                          padding: const EdgeInsets.all(20),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    _membership!.tierName,
                                    style: theme.textTheme.headlineSmall?.copyWith(
                                      fontWeight: FontWeight.bold,
                                      color: Colors.white,
                                    ),
                                  ),
                                  const Icon(
                                    Icons.stars,
                                    color: Colors.white,
                                    size: 32,
                                  ),
                                ],
                              ),
                              const SizedBox(height: 8),
                              Text(
                                'Member',
                                style: theme.textTheme.bodyMedium?.copyWith(
                                  color: Colors.white70,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Points
                      Text(
                        '${_membership!.currentPoints}',
                        style: theme.textTheme.displaySmall?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: theme.colorScheme.primary,
                        ),
                        textAlign: TextAlign.center,
                      ),
                      Text(
                        'Available Points',
                        style: theme.textTheme.bodyLarge?.copyWith(
                          color: Colors.grey,
                        ),
                        textAlign: TextAlign.center,
                      ),
                      const SizedBox(height: 16),

                      // Progress to next tier
                      if (_membership!.nextTierName != null) ...[
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              _membership!.tierName,
                              style: theme.textTheme.bodySmall?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            Text(
                              _membership!.nextTierName!,
                              style: theme.textTheme.bodySmall?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        LinearProgressIndicator(
                          value: _membership!.pointsToNextTier > 0
                              ? _membership!.lifetimePoints /
                                  (_membership!.lifetimePoints +
                                      _membership!.pointsToNextTier)
                              : 1.0,
                          minHeight: 8,
                          borderRadius: BorderRadius.circular(4),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${_membership!.pointsToNextTier} points to ${_membership!.nextTierName}',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: Colors.grey,
                          ),
                        ),
                      ],
                      const SizedBox(height: 24),

                      // Tier benefits
                      if (_tiers.isNotEmpty)
                        ...(_tiers.map((tier) {
                          final isCurrentTier = tier.id == _membership!.tierId;
                          return ExpansionTile(
                            initiallyExpanded: isCurrentTier,
                            leading: Icon(
                              Icons.workspace_premium,
                              color: isCurrentTier
                                  ? theme.colorScheme.primary
                                  : Colors.grey,
                            ),
                            title: Text(
                              tier.name,
                              style: TextStyle(
                                fontWeight: isCurrentTier
                                    ? FontWeight.bold
                                    : FontWeight.normal,
                              ),
                            ),
                            subtitle: Text(
                              '${tier.multiplier}x points multiplier',
                              style: theme.textTheme.bodySmall,
                            ),
                            children: tier.benefits.map((benefit) {
                              return ListTile(
                                leading: Icon(
                                  Icons.check_circle_outline,
                                  size: 20,
                                  color: theme.colorScheme.primary,
                                ),
                                title: Text(
                                  benefit,
                                  style: theme.textTheme.bodyMedium,
                                ),
                                dense: true,
                              );
                            }).toList(),
                          );
                        })),
                      const SizedBox(height: 24),

                      // Transaction history
                      Text(
                        'Transaction History',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 12),
                      if (_transactions.isEmpty)
                        Center(
                          child: Padding(
                            padding: const EdgeInsets.all(32),
                            child: Text(
                              'No transactions yet',
                              style: theme.textTheme.bodyMedium?.copyWith(
                                color: Colors.grey,
                              ),
                            ),
                          ),
                        )
                      else
                        ListView.builder(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          itemCount: _transactions.length,
                          itemBuilder: (context, index) {
                            final tx = _transactions[index];
                            final isEarned = tx.points > 0;
                            return ListTile(
                              leading: CircleAvatar(
                                backgroundColor: isEarned
                                    ? Colors.green.withOpacity(0.1)
                                    : Colors.red.withOpacity(0.1),
                                child: Icon(
                                  isEarned ? Icons.add : Icons.remove,
                                  color: isEarned ? Colors.green : Colors.red,
                                  size: 20,
                                ),
                              ),
                              title: Text(tx.description),
                              subtitle: Text(
                                '${tx.createdAt.day}/${tx.createdAt.month}/${tx.createdAt.year}',
                                style: theme.textTheme.bodySmall?.copyWith(
                                  color: Colors.grey,
                                ),
                              ),
                              trailing: Text(
                                '${isEarned ? '+' : ''}${tx.points}',
                                style: TextStyle(
                                  fontWeight: FontWeight.bold,
                                  color: isEarned ? Colors.green : Colors.red,
                                ),
                              ),
                              contentPadding: EdgeInsets.zero,
                            );
                          },
                        ),
                    ],
                  ),
                ),
    );
  }
}
