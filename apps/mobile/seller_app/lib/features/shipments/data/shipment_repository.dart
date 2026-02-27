import 'package:ecommerce_api_client/ecommerce_api_client.dart';

/// Represents dimensions of a shipment package.
class ShipmentDimensions {
  final double length;
  final double width;
  final double height;

  const ShipmentDimensions({
    required this.length,
    required this.width,
    required this.height,
  });

  @override
  String toString() => '${length}L x ${width}W x ${height}H';

  factory ShipmentDimensions.fromJson(Map<String, dynamic> json) {
    return ShipmentDimensions(
      length: (json['length'] as num).toDouble(),
      width: (json['width'] as num).toDouble(),
      height: (json['height'] as num).toDouble(),
    );
  }
}

/// Represents a shipment record.
class Shipment {
  final String id;
  final String orderId;
  final String orderNumber;
  final String carrier;
  final String trackingNumber;
  final String status; // pending, in_transit, delivered
  final double? weight;
  final ShipmentDimensions? dimensions;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Shipment({
    required this.id,
    required this.orderId,
    required this.orderNumber,
    required this.carrier,
    required this.trackingNumber,
    required this.status,
    this.weight,
    this.dimensions,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Shipment.fromJson(Map<String, dynamic> json) {
    return Shipment(
      id: json['id'] as String,
      orderId: json['orderId'] as String,
      orderNumber: json['orderNumber'] as String,
      carrier: json['carrier'] as String,
      trackingNumber: json['trackingNumber'] as String,
      status: json['status'] as String,
      weight: (json['weight'] as num?)?.toDouble(),
      dimensions: json['dimensions'] != null
          ? ShipmentDimensions.fromJson(
              json['dimensions'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}

/// Paginated result wrapper for shipments.
class PaginatedShipments {
  final List<Shipment> shipments;
  final int totalCount;
  final int currentPage;
  final int totalPages;

  const PaginatedShipments({
    required this.shipments,
    required this.totalCount,
    required this.currentPage,
    required this.totalPages,
  });

  factory PaginatedShipments.fromJson(Map<String, dynamic> json) {
    return PaginatedShipments(
      shipments: (json['data'] as List<dynamic>)
          .map((e) => Shipment.fromJson(e as Map<String, dynamic>))
          .toList(),
      totalCount: json['totalCount'] as int,
      currentPage: json['currentPage'] as int,
      totalPages: json['totalPages'] as int,
    );
  }
}

const String _sellerShipments = '/api/v1/seller/shipments';

/// Repository for managing shipments.
class ShipmentRepository {
  final ApiClient _apiClient;

  ShipmentRepository({required ApiClient apiClient})
      : _apiClient = apiClient;

  /// Fetches paginated list of shipments.
  Future<PaginatedShipments> getShipments({int page = 1}) async {
    final queryParams = <String, dynamic>{'page': page};

    final response = await _apiClient.get(
      _sellerShipments,
      queryParameters: queryParams,
    );
    return PaginatedShipments.fromJson(response.data as Map<String, dynamic>);
  }

  /// Creates a new shipment for an order.
  Future<Shipment> createShipment({
    required String orderId,
    required String carrier,
    required String trackingNumber,
    double? weight,
    ShipmentDimensions? dimensions,
  }) async {
    final body = <String, dynamic>{
      'orderId': orderId,
      'carrier': carrier,
      'trackingNumber': trackingNumber,
    };
    if (weight != null) body['weight'] = weight;
    if (dimensions != null) {
      body['dimensions'] = {
        'length': dimensions.length,
        'width': dimensions.width,
        'height': dimensions.height,
      };
    }

    final response = await _apiClient.post(
      _sellerShipments,
      data: body,
    );
    return Shipment.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches a single shipment by ID.
  Future<Shipment> getShipmentById(String id) async {
    final response = await _apiClient.get(
      '$_sellerShipments/$id',
    );
    return Shipment.fromJson(response.data as Map<String, dynamic>);
  }
}
