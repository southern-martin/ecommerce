import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';

class WishlistRepository {
  final ApiClient _apiClient;

  WishlistRepository(this._apiClient);

  Future<List<Product>> getWishlist() async {
    final response = await _apiClient.get('/wishlist');
    final List<dynamic> data = response.data['data'] ?? [];
    return data.map((json) => Product.fromJson(json)).toList();
  }

  Future<void> addToWishlist(String productId) async {
    await _apiClient.post('/wishlist', data: {'product_id': productId});
  }

  Future<void> removeFromWishlist(String productId) async {
    await _apiClient.delete('/wishlist/$productId');
  }
}
