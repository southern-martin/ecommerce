import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class TrackingEvent {
  final String status;
  final String description;
  final String location;
  final DateTime timestamp;

  const TrackingEvent({
    required this.status,
    required this.description,
    required this.location,
    required this.timestamp,
  });

  factory TrackingEvent.fromJson(Map<String, dynamic> json) {
    return TrackingEvent(
      status: json['status'] as String,
      description: json['description'] as String,
      location: json['location'] as String? ?? '',
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }
}

class TrackingInfo {
  final String orderId;
  final String orderNumber;
  final String carrier;
  final String trackingNumber;
  final String currentStatus;
  final DateTime? estimatedDelivery;
  final List<TrackingEvent> events;

  const TrackingInfo({
    required this.orderId,
    required this.orderNumber,
    required this.carrier,
    required this.trackingNumber,
    required this.currentStatus,
    this.estimatedDelivery,
    this.events = const [],
  });

  factory TrackingInfo.fromJson(Map<String, dynamic> json) {
    return TrackingInfo(
      orderId: json['orderId'] as String,
      orderNumber: json['orderNumber'] as String,
      carrier: json['carrier'] as String,
      trackingNumber: json['trackingNumber'] as String,
      currentStatus: json['currentStatus'] as String,
      estimatedDelivery: json['estimatedDelivery'] != null
          ? DateTime.parse(json['estimatedDelivery'] as String)
          : null,
      events: (json['events'] as List<dynamic>?)
              ?.map((e) => TrackingEvent.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
    );
  }
}

class TrackingRepository {
  final ApiClient _apiClient;

  TrackingRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<TrackingInfo> getTrackingInfo(String orderId) async {
    final response = await _apiClient.get('/orders/$orderId/tracking');
    return TrackingInfo.fromJson(response.data as Map<String, dynamic>);
  }
}
