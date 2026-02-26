import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class OrderItem {
  final String id;
  final String productId;
  final String name;
  final String imageUrl;
  final double price;
  final int quantity;
  final String? variantLabel;

  const OrderItem({
    required this.id,
    required this.productId,
    required this.name,
    required this.imageUrl,
    required this.price,
    required this.quantity,
    this.variantLabel,
  });

  factory OrderItem.fromJson(Map<String, dynamic> json) {
    return OrderItem(
      id: json['id'] as String,
      productId: json['productId'] as String,
      name: json['name'] as String,
      imageUrl: json['imageUrl'] as String,
      price: (json['price'] as num).toDouble(),
      quantity: json['quantity'] as int,
      variantLabel: json['variantLabel'] as String?,
    );
  }
}

class ShippingAddress {
  final String name;
  final String street;
  final String city;
  final String state;
  final String zip;
  final String phone;

  const ShippingAddress({
    required this.name,
    required this.street,
    required this.city,
    required this.state,
    required this.zip,
    required this.phone,
  });

  String get fullAddress => '$street, $city, $state $zip';

  factory ShippingAddress.fromJson(Map<String, dynamic> json) {
    return ShippingAddress(
      name: json['name'] as String,
      street: json['street'] as String,
      city: json['city'] as String,
      state: json['state'] as String,
      zip: json['zip'] as String,
      phone: json['phone'] as String? ?? '',
    );
  }
}

class OrderStatusStep {
  final String status;
  final String label;
  final DateTime? date;
  final bool isCompleted;

  const OrderStatusStep({
    required this.status,
    required this.label,
    this.date,
    required this.isCompleted,
  });

  factory OrderStatusStep.fromJson(Map<String, dynamic> json) {
    return OrderStatusStep(
      status: json['status'] as String,
      label: json['label'] as String,
      date: json['date'] != null ? DateTime.parse(json['date'] as String) : null,
      isCompleted: json['isCompleted'] as bool? ?? false,
    );
  }
}

class Order {
  final String id;
  final String orderNumber;
  final String status;
  final DateTime createdAt;
  final List<OrderItem> items;
  final double subtotal;
  final double shipping;
  final double tax;
  final double discount;
  final double total;
  final ShippingAddress? shippingAddress;
  final List<OrderStatusStep> timeline;
  final String? trackingNumber;
  final String? carrier;

  const Order({
    required this.id,
    required this.orderNumber,
    required this.status,
    required this.createdAt,
    required this.items,
    required this.subtotal,
    required this.shipping,
    required this.tax,
    required this.discount,
    required this.total,
    this.shippingAddress,
    this.timeline = const [],
    this.trackingNumber,
    this.carrier,
  });

  int get itemCount => items.fold(0, (sum, item) => sum + item.quantity);

  factory Order.fromJson(Map<String, dynamic> json) {
    return Order(
      id: json['id'] as String,
      orderNumber: json['orderNumber'] as String,
      status: json['status'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => OrderItem.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      subtotal: (json['subtotal'] as num? ?? 0).toDouble(),
      shipping: (json['shipping'] as num? ?? 0).toDouble(),
      tax: (json['tax'] as num? ?? 0).toDouble(),
      discount: (json['discount'] as num? ?? 0).toDouble(),
      total: (json['total'] as num).toDouble(),
      shippingAddress: json['shippingAddress'] != null
          ? ShippingAddress.fromJson(json['shippingAddress'] as Map<String, dynamic>)
          : null,
      timeline: (json['timeline'] as List<dynamic>?)
              ?.map((e) => OrderStatusStep.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      trackingNumber: json['trackingNumber'] as String?,
      carrier: json['carrier'] as String?,
    );
  }
}

class PaginatedOrders {
  final List<Order> orders;
  final int total;
  final int page;
  final int pageSize;

  const PaginatedOrders({
    required this.orders,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  bool get hasMore => page * pageSize < total;

  factory PaginatedOrders.fromJson(Map<String, dynamic> json) {
    return PaginatedOrders(
      orders: (json['data'] as List<dynamic>)
          .map((e) => Order.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      page: json['page'] as int,
      pageSize: json['pageSize'] as int,
    );
  }
}

class OrderRepository {
  final ApiClient _apiClient;

  OrderRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<PaginatedOrders> getOrders({int page = 1, String? status}) async {
    final queryParams = <String, dynamic>{'page': page};
    if (status != null) queryParams['status'] = status;

    final response = await _apiClient.get('/orders', queryParameters: queryParams);
    return PaginatedOrders.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Order> getOrderById(String id) async {
    final response = await _apiClient.get('/orders/$id');
    return Order.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> cancelOrder(String id) async {
    await _apiClient.post('/orders/$id/cancel');
  }
}
