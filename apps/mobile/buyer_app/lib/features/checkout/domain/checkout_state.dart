import '../data/checkout_repository.dart';

sealed class CheckoutState {
  const CheckoutState();
}

class CheckoutInitial extends CheckoutState {
  const CheckoutInitial();
}

class CheckoutLoading extends CheckoutState {
  const CheckoutLoading();
}

class CheckoutAddressStep extends CheckoutState {
  const CheckoutAddressStep();
}

class CheckoutShippingStep extends CheckoutState {
  final List<ShippingRate> rates;

  const CheckoutShippingStep({required this.rates});
}

class CheckoutPaymentStep extends CheckoutState {
  final List<PaymentMethod> methods;
  final OrderSummary summary;

  const CheckoutPaymentStep({
    required this.methods,
    required this.summary,
  });
}

class CheckoutPlacing extends CheckoutState {
  const CheckoutPlacing();
}

class CheckoutConfirmed extends CheckoutState {
  final String orderId;
  final String orderNumber;

  const CheckoutConfirmed({
    required this.orderId,
    required this.orderNumber,
  });
}

class CheckoutError extends CheckoutState {
  final String message;

  const CheckoutError({required this.message});
}

class OrderSummary {
  final double subtotal;
  final double shipping;
  final double tax;
  final double discount;
  final double total;

  const OrderSummary({
    required this.subtotal,
    required this.shipping,
    required this.tax,
    required this.discount,
    required this.total,
  });
}
