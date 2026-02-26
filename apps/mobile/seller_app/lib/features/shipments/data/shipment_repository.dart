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
}

/// Repository for managing shipments.
class ShipmentRepository {
  /// Fetches paginated list of shipments.
  Future<PaginatedShipments> getShipments({int page = 1}) async {
    // TODO: Replace with actual API call to /seller/shipments
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    final carriers = ['FedEx', 'UPS', 'USPS', 'DHL'];
    final statuses = ['pending', 'in_transit', 'delivered'];

    final allShipments = List.generate(
      12,
      (i) => Shipment(
        id: 'ship_${i + 1}',
        orderId: 'order_${i + 1}',
        orderNumber: 'ORD-${1001 + i}',
        carrier: carriers[i % carriers.length],
        trackingNumber: 'TRK${900000 + i * 111}',
        status: statuses[i % statuses.length],
        weight: i % 2 == 0 ? 2.5 + (i * 0.5) : null,
        dimensions: i % 3 == 0
            ? ShipmentDimensions(
                length: 10.0 + i,
                width: 8.0 + i,
                height: 4.0 + i,
              )
            : null,
        createdAt: now.subtract(Duration(days: i * 2)),
        updatedAt: now.subtract(Duration(days: i)),
      ),
    );

    return PaginatedShipments(
      shipments: allShipments,
      totalCount: allShipments.length,
      currentPage: page,
      totalPages: 1,
    );
  }

  /// Creates a new shipment for an order.
  Future<Shipment> createShipment({
    required String orderId,
    required String carrier,
    required String trackingNumber,
    double? weight,
    ShipmentDimensions? dimensions,
  }) async {
    // TODO: Replace with actual API call to /seller/shipments
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return Shipment(
      id: 'ship_new_${now.millisecondsSinceEpoch}',
      orderId: orderId,
      orderNumber: 'ORD-$orderId',
      carrier: carrier,
      trackingNumber: trackingNumber,
      status: 'pending',
      weight: weight,
      dimensions: dimensions,
      createdAt: now,
      updatedAt: now,
    );
  }

  /// Fetches a single shipment by ID.
  Future<Shipment> getShipmentById(String id) async {
    // TODO: Replace with actual API call to /seller/shipments/:id
    await Future.delayed(const Duration(seconds: 1));

    final now = DateTime.now();
    return Shipment(
      id: id,
      orderId: 'order_1',
      orderNumber: 'ORD-1001',
      carrier: 'FedEx',
      trackingNumber: 'TRK900000',
      status: 'in_transit',
      weight: 3.5,
      dimensions: const ShipmentDimensions(length: 12, width: 10, height: 6),
      createdAt: now.subtract(const Duration(days: 3)),
      updatedAt: now.subtract(const Duration(days: 1)),
    );
  }
}
