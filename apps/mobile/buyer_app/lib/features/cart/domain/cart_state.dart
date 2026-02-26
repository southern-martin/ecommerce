import '../data/cart_repository.dart';

sealed class CartState {
  const CartState();
}

class CartInitial extends CartState {
  const CartInitial();
}

class CartLoading extends CartState {
  const CartLoading();
}

class CartLoaded extends CartState {
  final Cart cart;

  const CartLoaded({required this.cart});
}

class CartError extends CartState {
  final String message;

  const CartError({required this.message});
}
