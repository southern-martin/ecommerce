import 'package:json_annotation/json_annotation.dart';

part 'order.g.dart';

/// Represents a placed order.
@JsonSerializable()
class Order {
  final String id;
  final String orderNumber;
  final OrderStatus status;
  final List<OrderItem> items;
  final int subtotalCents;
  final int taxCents;
  final int shippingCents;
  final int totalCents;
  final Address? shippingAddress;
  final DateTime createdAt;

  const Order({
    required this.id,
    required this.orderNumber,
    required this.status,
    this.items = const [],
    required this.subtotalCents,
    this.taxCents = 0,
    this.shippingCents = 0,
    required this.totalCents,
    this.shippingAddress,
    required this.createdAt,
  });

  /// Total number of individual items in the order.
  int get itemCount => items.fold(0, (sum, item) => sum + item.quantity);

  /// Whether the order can be cancelled.
  bool get isCancellable =>
      status == OrderStatus.pending || status == OrderStatus.confirmed;

  /// Whether the order can be returned.
  bool get isReturnable => status == OrderStatus.delivered;

  factory Order.fromJson(Map<String, dynamic> json) => _$OrderFromJson(json);

  Map<String, dynamic> toJson() => _$OrderToJson(this);

  Order copyWith({
    String? id,
    String? orderNumber,
    OrderStatus? status,
    List<OrderItem>? items,
    int? subtotalCents,
    int? taxCents,
    int? shippingCents,
    int? totalCents,
    Address? shippingAddress,
    DateTime? createdAt,
  }) {
    return Order(
      id: id ?? this.id,
      orderNumber: orderNumber ?? this.orderNumber,
      status: status ?? this.status,
      items: items ?? this.items,
      subtotalCents: subtotalCents ?? this.subtotalCents,
      taxCents: taxCents ?? this.taxCents,
      shippingCents: shippingCents ?? this.shippingCents,
      totalCents: totalCents ?? this.totalCents,
      shippingAddress: shippingAddress ?? this.shippingAddress,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Order && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() =>
      'Order(id: $id, number: $orderNumber, status: $status, total: $totalCents)';
}

/// Represents a single line item within an order.
@JsonSerializable()
class OrderItem {
  final String id;
  final String productId;
  final String name;
  final int quantity;
  final int priceCents;
  final String? imageUrl;

  const OrderItem({
    required this.id,
    required this.productId,
    required this.name,
    required this.quantity,
    required this.priceCents,
    this.imageUrl,
  });

  /// Total price for this line item in cents.
  int get totalCents => priceCents * quantity;

  factory OrderItem.fromJson(Map<String, dynamic> json) =>
      _$OrderItemFromJson(json);

  Map<String, dynamic> toJson() => _$OrderItemToJson(this);

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is OrderItem && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// Represents a physical address.
@JsonSerializable()
class Address {
  final String street;
  final String city;
  final String state;
  final String zip;
  final String country;
  final String firstName;
  final String lastName;
  final String? phone;

  const Address({
    required this.street,
    required this.city,
    required this.state,
    required this.zip,
    required this.country,
    required this.firstName,
    required this.lastName,
    this.phone,
  });

  /// Full name for the address.
  String get fullName => '$firstName $lastName';

  /// Formatted single-line address.
  String get formatted => '$street, $city, $state $zip, $country';

  factory Address.fromJson(Map<String, dynamic> json) =>
      _$AddressFromJson(json);

  Map<String, dynamic> toJson() => _$AddressToJson(this);

  Address copyWith({
    String? street,
    String? city,
    String? state,
    String? zip,
    String? country,
    String? firstName,
    String? lastName,
    String? phone,
  }) {
    return Address(
      street: street ?? this.street,
      city: city ?? this.city,
      state: state ?? this.state,
      zip: zip ?? this.zip,
      country: country ?? this.country,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      phone: phone ?? this.phone,
    );
  }

  @override
  String toString() => 'Address($formatted)';
}

/// Status of an order throughout its lifecycle.
@JsonEnum(valueField: 'value')
enum OrderStatus {
  pending('pending'),
  confirmed('confirmed'),
  processing('processing'),
  shipped('shipped'),
  outForDelivery('out_for_delivery'),
  delivered('delivered'),
  cancelled('cancelled'),
  returned('returned'),
  refunded('refunded');

  final String value;
  const OrderStatus(this.value);

  /// Human-readable display label.
  String get displayLabel {
    switch (this) {
      case OrderStatus.pending:
        return 'Pending';
      case OrderStatus.confirmed:
        return 'Confirmed';
      case OrderStatus.processing:
        return 'Processing';
      case OrderStatus.shipped:
        return 'Shipped';
      case OrderStatus.outForDelivery:
        return 'Out for Delivery';
      case OrderStatus.delivered:
        return 'Delivered';
      case OrderStatus.cancelled:
        return 'Cancelled';
      case OrderStatus.returned:
        return 'Returned';
      case OrderStatus.refunded:
        return 'Refunded';
    }
  }
}
