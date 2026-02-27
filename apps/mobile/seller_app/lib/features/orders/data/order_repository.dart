import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Represents an item within a seller order.
class OrderItem {
  final String id;
  final String productId;
  final String productName;
  final String? imageUrl;
  final int quantity;
  final double price;
  final double total;

  const OrderItem({
    required this.id,
    required this.productId,
    required this.productName,
    this.imageUrl,
    required this.quantity,
    required this.price,
    required this.total,
  });

  factory OrderItem.fromJson(Map<String, dynamic> json) {
    return OrderItem(
      id: json['id'] as String,
      productId: json['productId'] as String,
      productName: json['productName'] as String,
      imageUrl: json['imageUrl'] as String?,
      quantity: json['quantity'] as int,
      price: (json['price'] as num).toDouble(),
      total: (json['total'] as num).toDouble(),
    );
  }
}

/// Represents a customer's shipping address.
class CustomerInfo {
  final String name;
  final String email;
  final String phone;
  final String address;
  final String city;
  final String state;
  final String zipCode;

  const CustomerInfo({
    required this.name,
    required this.email,
    required this.phone,
    required this.address,
    required this.city,
    required this.state,
    required this.zipCode,
  });

  String get fullAddress => '$address, $city, $state $zipCode';

  factory CustomerInfo.fromJson(Map<String, dynamic> json) {
    return CustomerInfo(
      name: json['name'] as String,
      email: json['email'] as String,
      phone: json['phone'] as String,
      address: json['address'] as String,
      city: json['city'] as String,
      state: json['state'] as String,
      zipCode: json['zipCode'] as String,
    );
  }
}

/// Represents a seller order.
class SellerOrder {
  final String id;
  final String orderNumber;
  final String customerName;
  final CustomerInfo? customerInfo;
  final List<OrderItem> items;
  final double total;
  final String status;
  final String? trackingNumber;
  final DateTime createdAt;
  final DateTime updatedAt;

  const SellerOrder({
    required this.id,
    required this.orderNumber,
    required this.customerName,
    this.customerInfo,
    required this.items,
    required this.total,
    required this.status,
    this.trackingNumber,
    required this.createdAt,
    required this.updatedAt,
  });

  int get itemCount => items.fold(0, (sum, item) => sum + item.quantity);

  factory SellerOrder.fromJson(Map<String, dynamic> json) {
    return SellerOrder(
      id: json['id'] as String,
      orderNumber: json['orderNumber'] as String,
      customerName: json['customerName'] as String,
      customerInfo: json['customerInfo'] != null
          ? CustomerInfo.fromJson(json['customerInfo'] as Map<String, dynamic>)
          : null,
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => OrderItem.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      total: (json['total'] as num).toDouble(),
      status: json['status'] as String,
      trackingNumber: json['trackingNumber'] as String?,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}

/// Paginated result wrapper for orders.
class PaginatedOrders {
  final List<SellerOrder> orders;
  final int totalCount;
  final int currentPage;
  final int totalPages;

  const PaginatedOrders({
    required this.orders,
    required this.totalCount,
    required this.currentPage,
    required this.totalPages,
  });

  factory PaginatedOrders.fromJson(Map<String, dynamic> json) {
    return PaginatedOrders(
      orders: (json['data'] as List<dynamic>)
          .map((e) => SellerOrder.fromJson(e as Map<String, dynamic>))
          .toList(),
      totalCount: json['totalCount'] as int,
      currentPage: json['currentPage'] as int,
      totalPages: json['totalPages'] as int,
    );
  }
}

/// Order statistics summary.
class OrderStats {
  final int totalOrders;
  final int pendingOrders;
  final int processingOrders;
  final int shippedOrders;
  final int deliveredOrders;
  final int cancelledOrders;

  const OrderStats({
    required this.totalOrders,
    required this.pendingOrders,
    required this.processingOrders,
    required this.shippedOrders,
    required this.deliveredOrders,
    required this.cancelledOrders,
  });

  factory OrderStats.fromJson(Map<String, dynamic> json) {
    return OrderStats(
      totalOrders: json['totalOrders'] as int,
      pendingOrders: json['pendingOrders'] as int,
      processingOrders: json['processingOrders'] as int,
      shippedOrders: json['shippedOrders'] as int,
      deliveredOrders: json['deliveredOrders'] as int,
      cancelledOrders: json['cancelledOrders'] as int,
    );
  }
}

/// Repository for managing seller orders.
class SellerOrderRepository {
  final ApiClient _apiClient;

  SellerOrderRepository({required ApiClient apiClient})
      : _apiClient = apiClient;

  /// Fetches paginated list of seller orders with optional status filter.
  Future<PaginatedOrders> getOrders({
    int page = 1,
    String? status,
  }) async {
    final queryParams = <String, dynamic>{'page': page};
    if (status != null) queryParams['status'] = status;

    final response = await _apiClient.get(
      ApiEndpoints.sellerOrders,
      queryParameters: queryParams,
    );
    return PaginatedOrders.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches a single order by ID.
  Future<SellerOrder> getOrderById(String id) async {
    final response = await _apiClient.get(
      '${ApiEndpoints.sellerOrders}/$id',
    );
    return SellerOrder.fromJson(response.data as Map<String, dynamic>);
  }

  /// Updates the status of an order.
  Future<SellerOrder> updateOrderStatus(
    String id,
    String status, {
    String? trackingNumber,
  }) async {
    final body = <String, dynamic>{'status': status};
    if (trackingNumber != null) body['trackingNumber'] = trackingNumber;

    final response = await _apiClient.put(
      '${ApiEndpoints.sellerOrders}/$id/status',
      data: body,
    );
    return SellerOrder.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches order statistics summary.
  Future<OrderStats> getOrderStats() async {
    final response = await _apiClient.get(
      '${ApiEndpoints.sellerOrders}/stats',
    );
    return OrderStats.fromJson(response.data as Map<String, dynamic>);
  }
}
