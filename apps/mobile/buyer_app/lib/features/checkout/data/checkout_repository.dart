import 'package:ecommerce_api_client/ecommerce_api_client.dart';

class ShippingRate {
  final String id;
  final String name;
  final String carrier;
  final double price;
  final String estimatedDays;

  const ShippingRate({
    required this.id,
    required this.name,
    required this.carrier,
    required this.price,
    required this.estimatedDays,
  });

  factory ShippingRate.fromJson(Map<String, dynamic> json) {
    return ShippingRate(
      id: json['id'] as String,
      name: json['name'] as String,
      carrier: json['carrier'] as String,
      price: (json['price'] as num).toDouble(),
      estimatedDays: json['estimatedDays'] as String,
    );
  }
}

class PaymentMethod {
  final String id;
  final String type;
  final String label;
  final String? last4;
  final String? brand;

  const PaymentMethod({
    required this.id,
    required this.type,
    required this.label,
    this.last4,
    this.brand,
  });

  factory PaymentMethod.fromJson(Map<String, dynamic> json) {
    return PaymentMethod(
      id: json['id'] as String,
      type: json['type'] as String,
      label: json['label'] as String,
      last4: json['last4'] as String?,
      brand: json['brand'] as String?,
    );
  }
}

class TaxResult {
  final double taxAmount;
  final double subtotal;
  final double total;

  const TaxResult({
    required this.taxAmount,
    required this.subtotal,
    required this.total,
  });

  factory TaxResult.fromJson(Map<String, dynamic> json) {
    return TaxResult(
      taxAmount: (json['taxAmount'] as num).toDouble(),
      subtotal: (json['subtotal'] as num).toDouble(),
      total: (json['total'] as num).toDouble(),
    );
  }
}

class OrderResult {
  final String orderId;
  final String orderNumber;

  const OrderResult({required this.orderId, required this.orderNumber});

  factory OrderResult.fromJson(Map<String, dynamic> json) {
    return OrderResult(
      orderId: json['orderId'] as String,
      orderNumber: json['orderNumber'] as String,
    );
  }
}

class CheckoutRepository {
  final ApiClient _apiClient;

  CheckoutRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  Future<List<ShippingRate>> getShippingRates(String addressId) async {
    final response = await _apiClient.get('/checkout/shipping-rates', queryParameters: {
      'addressId': addressId,
    });
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => ShippingRate.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<TaxResult> calculateTax({
    required String cartId,
    required String addressId,
  }) async {
    final response = await _apiClient.post('/checkout/calculate-tax', data: {
      'cartId': cartId,
      'addressId': addressId,
    });
    return TaxResult.fromJson(response.data as Map<String, dynamic>);
  }

  Future<OrderResult> placeOrder({
    required String shippingMethodId,
    required String paymentMethodId,
    required String addressId,
  }) async {
    final response = await _apiClient.post('/checkout/place-order', data: {
      'shippingMethodId': shippingMethodId,
      'paymentMethodId': paymentMethodId,
      'addressId': addressId,
    });
    return OrderResult.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<PaymentMethod>> getPaymentMethods() async {
    final response = await _apiClient.get('/checkout/payment-methods');
    final List<dynamic> data = response.data as List<dynamic>;
    return data.map((e) => PaymentMethod.fromJson(e as Map<String, dynamic>)).toList();
  }
}
