import 'package:flutter/material.dart';

import '../../../../core/di/injection.dart';
import '../../data/tracking_repository.dart';

class TrackingPage extends StatefulWidget {
  final String orderId;

  const TrackingPage({super.key, required this.orderId});

  @override
  State<TrackingPage> createState() => _TrackingPageState();
}

class _TrackingPageState extends State<TrackingPage> {
  final TrackingRepository _trackingRepo = getIt<TrackingRepository>();

  TrackingInfo? _tracking;
  bool _isLoading = true;
  String? _error;

  static const List<_TimelineStep> _steps = [
    _TimelineStep(key: 'ordered', label: 'Ordered', icon: Icons.receipt_long),
    _TimelineStep(key: 'shipped', label: 'Shipped', icon: Icons.local_shipping),
    _TimelineStep(key: 'in_transit', label: 'In Transit', icon: Icons.flight),
    _TimelineStep(key: 'out_for_delivery', label: 'Out for Delivery', icon: Icons.delivery_dining),
    _TimelineStep(key: 'delivered', label: 'Delivered', icon: Icons.check_circle),
  ];

  @override
  void initState() {
    super.initState();
    _loadTracking();
  }

  Future<void> _loadTracking() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final tracking = await _trackingRepo.getTrackingInfo(widget.orderId);
      if (mounted) {
        setState(() {
          _tracking = tracking;
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

  int _currentStepIndex() {
    if (_tracking == null) return 0;
    final status = _tracking!.currentStatus.toLowerCase().replaceAll(' ', '_');
    for (int i = _steps.length - 1; i >= 0; i--) {
      if (_steps[i].key == status) return i;
    }
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Order Tracking'),
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
                      Text('Failed to load tracking', style: theme.textTheme.titleMedium),
                      const SizedBox(height: 8),
                      FilledButton(onPressed: _loadTracking, child: const Text('Retry')),
                    ],
                  ),
                )
              : RefreshIndicator(
                  onRefresh: _loadTracking,
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      // Order info card
                      Card(
                        child: Padding(
                          padding: const EdgeInsets.all(16),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Order #${_tracking!.orderNumber}',
                                style: theme.textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                              const SizedBox(height: 12),
                              Row(
                                children: [
                                  Icon(Icons.local_shipping_outlined,
                                      size: 16, color: Colors.grey.shade600),
                                  const SizedBox(width: 8),
                                  Text(
                                    'Carrier: ${_tracking!.carrier}',
                                    style: theme.textTheme.bodyMedium,
                                  ),
                                ],
                              ),
                              const SizedBox(height: 8),
                              Row(
                                children: [
                                  Icon(Icons.tag,
                                      size: 16, color: Colors.grey.shade600),
                                  const SizedBox(width: 8),
                                  Expanded(
                                    child: Text(
                                      'Tracking #: ${_tracking!.trackingNumber}',
                                      style: theme.textTheme.bodyMedium,
                                    ),
                                  ),
                                ],
                              ),
                              if (_tracking!.estimatedDelivery != null) ...[
                                const SizedBox(height: 8),
                                Row(
                                  children: [
                                    Icon(Icons.calendar_today,
                                        size: 16, color: Colors.grey.shade600),
                                    const SizedBox(width: 8),
                                    Text(
                                      'Est. Delivery: ${_tracking!.estimatedDelivery!.day}/${_tracking!.estimatedDelivery!.month}/${_tracking!.estimatedDelivery!.year}',
                                      style: theme.textTheme.bodyMedium,
                                    ),
                                  ],
                                ),
                              ],
                            ],
                          ),
                        ),
                      ),
                      const SizedBox(height: 24),

                      // Timeline
                      Text(
                        'Delivery Progress',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),
                      _buildTimeline(theme),

                      // Events detail
                      if (_tracking!.events.isNotEmpty) ...[
                        const SizedBox(height: 24),
                        Text(
                          'Tracking Events',
                          style: theme.textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 12),
                        ...(_tracking!.events.map((event) {
                          return Card(
                            margin: const EdgeInsets.only(bottom: 8),
                            child: ListTile(
                              title: Text(event.description),
                              subtitle: Text(
                                event.location.isNotEmpty
                                    ? '${event.location} - ${event.timestamp.day}/${event.timestamp.month}/${event.timestamp.year} ${event.timestamp.hour}:${event.timestamp.minute.toString().padLeft(2, '0')}'
                                    : '${event.timestamp.day}/${event.timestamp.month}/${event.timestamp.year} ${event.timestamp.hour}:${event.timestamp.minute.toString().padLeft(2, '0')}',
                              ),
                              leading: Icon(
                                Icons.circle,
                                size: 12,
                                color: theme.colorScheme.primary,
                              ),
                            ),
                          );
                        })),
                      ],
                    ],
                  ),
                ),
    );
  }

  Widget _buildTimeline(ThemeData theme) {
    final currentStep = _currentStepIndex();

    return Column(
      children: List.generate(_steps.length, (index) {
        final step = _steps[index];
        final isCompleted = index <= currentStep;
        final isLast = index == _steps.length - 1;
        final color = isCompleted ? theme.colorScheme.primary : Colors.grey.shade300;

        return Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Timeline line and icon
            SizedBox(
              width: 40,
              child: Column(
                children: [
                  Container(
                    width: 36,
                    height: 36,
                    decoration: BoxDecoration(
                      color: isCompleted
                          ? theme.colorScheme.primary
                          : Colors.grey.shade200,
                      shape: BoxShape.circle,
                    ),
                    child: Icon(
                      step.icon,
                      size: 18,
                      color: isCompleted ? Colors.white : Colors.grey,
                    ),
                  ),
                  if (!isLast)
                    Container(
                      width: 2,
                      height: 40,
                      color: color,
                    ),
                ],
              ),
            ),
            const SizedBox(width: 12),

            // Label
            Expanded(
              child: Padding(
                padding: const EdgeInsets.only(top: 6),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      step.label,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight:
                            isCompleted ? FontWeight.bold : FontWeight.normal,
                        color: isCompleted ? null : Colors.grey,
                      ),
                    ),
                    SizedBox(height: isLast ? 0 : 24),
                  ],
                ),
              ),
            ),
          ],
        );
      }),
    );
  }
}

class _TimelineStep {
  final String key;
  final String label;
  final IconData icon;

  const _TimelineStep({
    required this.key,
    required this.label,
    required this.icon,
  });
}
