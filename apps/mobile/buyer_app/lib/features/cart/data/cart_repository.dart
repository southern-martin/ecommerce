import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class CartItem {
  final String id;
  final String productId;
  final String name;
  final String imageUrl;
  final double price;
  final int quantity;
  final String? variantId;
  final String? variantLabel;

  const CartItem({
    required this.id,
    required this.productId,
    required this.name,
    required this.imageUrl,
    required this.price,
    required this.quantity,
    this.variantId,
    this.variantLabel,
  });

  double get total => price * quantity;

  factory CartItem.fromJson(Map<String, dynamic> json) {
    return CartItem(
      id: json['id'] as String,
      productId: json['productId'] as String,
      name: json['name'] as String,
      imageUrl: json['imageUrl'] as String,
      price: (json['price'] as num).toDouble(),
      quantity: json['quantity'] as int,
      variantId: json['variantId'] as String?,
      variantLabel: json['variantLabel'] as String?,
    );
  }
}

class Cart {
  final String id;
  final List<CartItem> items;
  final double subtotal;
  final double shippingEstimate;
  final double tax;
  final double discount;
  final double total;
  final String? couponCode;
  final int loyaltyPointsApplied;

  const Cart({
    required this.id,
    required this.items,
    required this.subtotal,
    required this.shippingEstimate,
    required this.tax,
    required this.discount,
    required this.total,
    this.couponCode,
    this.loyaltyPointsApplied = 0,
  });

  int get itemCount => items.fold(0, (sum, item) => sum + item.quantity);

  factory Cart.fromJson(Map<String, dynamic> json) {
    return Cart(
      id: json['id'] as String,
      items: (json['items'] as List<dynamic>)
          .map((e) => CartItem.fromJson(e as Map<String, dynamic>))
          .toList(),
      subtotal: (json['subtotal'] as num).toDouble(),
      shippingEstimate: (json['shippingEstimate'] as num? ?? 0).toDouble(),
      tax: (json['tax'] as num? ?? 0).toDouble(),
      discount: (json['discount'] as num? ?? 0).toDouble(),
      total: (json['total'] as num).toDouble(),
      couponCode: json['couponCode'] as String?,
      loyaltyPointsApplied: json['loyaltyPointsApplied'] as int? ?? 0,
    );
  }
}

class CartRepository {
  final ApiClient _apiClient;

  CartRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<Cart> getCart() async {
    final response = await _apiClient.get('/cart');
    return Cart.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Cart> addToCart({
    required String productId,
    required int quantity,
    String? variantId,
  }) async {
    final response = await _apiClient.post('/cart/items', data: {
      'productId': productId,
      'quantity': quantity,
      if (variantId != null) 'variantId': variantId,
    });
    return Cart.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Cart> updateQuantity({
    required String itemId,
    required int quantity,
  }) async {
    final response = await _apiClient.put('/cart/items/$itemId', data: {
      'quantity': quantity,
    });
    return Cart.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Cart> removeFromCart(String itemId) async {
    final response = await _apiClient.delete('/cart/items/$itemId');
    return Cart.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Cart> applyCoupon(String code) async {
    final response = await _apiClient.post('/cart/coupon', data: {
      'code': code,
    });
    return Cart.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Cart> applyPoints(int points) async {
    final response = await _apiClient.post('/cart/loyalty-points', data: {
      'points': points,
    });
    return Cart.fromJson(response.data as Map<String, dynamic>);
  }
}
