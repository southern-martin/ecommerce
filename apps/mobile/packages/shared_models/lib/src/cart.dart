import 'package:json_annotation/json_annotation.dart';

part 'cart.g.dart';

/// Represents an item in the shopping cart.
@JsonSerializable()
class CartItem {
  final String id;
  final String productId;
  final String name;
  final String slug;
  final String? imageUrl;
  final int priceCents;
  final int quantity;
  final String? variantId;
  final String? variantName;

  const CartItem({
    required this.id,
    required this.productId,
    required this.name,
    required this.slug,
    this.imageUrl,
    required this.priceCents,
    required this.quantity,
    this.variantId,
    this.variantName,
  });

  /// Total price for this line item (price * quantity) in cents.
  int get totalCents => priceCents * quantity;

  factory CartItem.fromJson(Map<String, dynamic> json) =>
      _$CartItemFromJson(json);

  Map<String, dynamic> toJson() => _$CartItemToJson(this);

  CartItem copyWith({
    String? id,
    String? productId,
    String? name,
    String? slug,
    String? imageUrl,
    int? priceCents,
    int? quantity,
    String? variantId,
    String? variantName,
  }) {
    return CartItem(
      id: id ?? this.id,
      productId: productId ?? this.productId,
      name: name ?? this.name,
      slug: slug ?? this.slug,
      imageUrl: imageUrl ?? this.imageUrl,
      priceCents: priceCents ?? this.priceCents,
      quantity: quantity ?? this.quantity,
      variantId: variantId ?? this.variantId,
      variantName: variantName ?? this.variantName,
    );
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is CartItem && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() =>
      'CartItem(id: $id, product: $name, qty: $quantity, price: $priceCents)';
}

/// Represents the entire shopping cart.
@JsonSerializable()
class Cart {
  final String id;
  final List<CartItem> items;

  const Cart({
    required this.id,
    this.items = const [],
  });

  /// Subtotal of all items in cents.
  int get subtotalCents =>
      items.fold(0, (sum, item) => sum + item.totalCents);

  /// Total number of individual items (sum of quantities).
  int get itemCount =>
      items.fold(0, (sum, item) => sum + item.quantity);

  /// Whether the cart is empty.
  bool get isEmpty => items.isEmpty;

  /// Whether the cart has items.
  bool get isNotEmpty => items.isNotEmpty;

  factory Cart.fromJson(Map<String, dynamic> json) => _$CartFromJson(json);

  Map<String, dynamic> toJson() => _$CartToJson(this);

  Cart copyWith({
    String? id,
    List<CartItem>? items,
  }) {
    return Cart(
      id: id ?? this.id,
      items: items ?? this.items,
    );
  }

  @override
  String toString() =>
      'Cart(id: $id, items: ${items.length}, subtotal: $subtotalCents)';
}
