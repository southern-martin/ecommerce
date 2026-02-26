import 'package:ecommerce_core/ecommerce_core.dart';

/// Repository handling seller authentication operations.
///
/// Uses [SecureStorage] from ecommerce_core to persist tokens and credentials
/// securely on the device keychain/keystore.
class SellerAuthRepository {
  final SecureStorage _secureStorage;

  SellerAuthRepository({required SecureStorage secureStorage})
      : _secureStorage = secureStorage;

  /// Authenticates a seller with [email] and [password].
  ///
  /// On success, persists the access and refresh tokens to secure storage.
  /// Returns `true` if login succeeded, `false` otherwise.
  Future<bool> login({
    required String email,
    required String password,
  }) async {
    try {
      // TODO: Replace with actual ApiClient call once ecommerce_api_client is wired
      // final response = await _apiClient.post('/seller/auth/login', body: {
      //   'email': email,
      //   'password': password,
      // });

      // Simulate API call
      await Future.delayed(const Duration(seconds: 1));

      // Simulated token response
      const accessToken = 'seller_access_token_placeholder';
      const refreshToken = 'seller_refresh_token_placeholder';
      final expiry = DateTime.now().add(const Duration(hours: 24));

      await _secureStorage.setAccessToken(accessToken);
      await _secureStorage.setRefreshToken(refreshToken);
      await _secureStorage.setTokenExpiry(expiry);
      await _secureStorage.write(key: 'seller_email', value: email);

      return true;
    } catch (e) {
      return false;
    }
  }

  /// Logs out the current seller by clearing all stored credentials.
  Future<void> logout() async {
    await _secureStorage.deleteAccessToken();
    await _secureStorage.deleteRefreshToken();
    await _secureStorage.delete(key: 'seller_email');
    await _secureStorage.clearAll();
  }

  /// Attempts to refresh the access token using the stored refresh token.
  ///
  /// Returns `true` if the token was successfully refreshed.
  Future<bool> refreshToken() async {
    try {
      final currentRefreshToken = await _secureStorage.getRefreshToken();
      if (currentRefreshToken == null) return false;

      // TODO: Replace with actual ApiClient call
      // final response = await _apiClient.post('/seller/auth/refresh', body: {
      //   'refresh_token': currentRefreshToken,
      // });

      await Future.delayed(const Duration(milliseconds: 500));

      const newAccessToken = 'refreshed_seller_access_token';
      final newExpiry = DateTime.now().add(const Duration(hours: 24));

      await _secureStorage.setAccessToken(newAccessToken);
      await _secureStorage.setTokenExpiry(newExpiry);

      return true;
    } catch (e) {
      return false;
    }
  }

  /// Checks whether the seller is currently authenticated.
  ///
  /// Returns `true` if a valid, non-expired access token exists.
  Future<bool> isAuthenticated() async {
    final token = await _secureStorage.getAccessToken();
    if (token == null) return false;

    final isExpired = await _secureStorage.isTokenExpired();
    if (isExpired) {
      return await refreshToken();
    }

    return true;
  }

  /// Returns the currently stored seller email, or null.
  Future<String?> getCurrentSellerEmail() async {
    return _secureStorage.read(key: 'seller_email');
  }
}
