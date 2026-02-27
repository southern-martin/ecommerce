import 'package:ecommerce_api_client/ecommerce_api_client.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Repository handling seller authentication operations.
///
/// Uses [ApiClient] for network requests and [SecureStorage] to persist
/// tokens and credentials securely on the device keychain/keystore.
class SellerAuthRepository {
  final ApiClient _apiClient;
  final SecureStorage _secureStorage;

  SellerAuthRepository({
    required ApiClient apiClient,
    required SecureStorage secureStorage,
  })  : _apiClient = apiClient,
        _secureStorage = secureStorage;

  /// Authenticates a seller with [email] and [password].
  ///
  /// On success, persists the access and refresh tokens to secure storage.
  /// Returns `true` if login succeeded, `false` otherwise.
  Future<bool> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await _apiClient.post(
        ApiEndpoints.login,
        data: {
          'email': email,
          'password': password,
        },
      );

      final data = response.data as Map<String, dynamic>;
      final accessToken = data['access_token'] as String;
      final refreshToken = data['refresh_token'] as String;
      final expiresIn = data['expires_in'] as int?;
      final expiry = DateTime.now().add(
        Duration(seconds: expiresIn ?? 86400),
      );

      await _secureStorage.setAccessToken(accessToken);
      await _secureStorage.setRefreshToken(refreshToken);
      await _secureStorage.setTokenExpiry(expiry);
      await _secureStorage.write(key: 'seller_email', value: email);

      return true;
    } catch (e) {
      return false;
    }
  }

  /// Logs out the current seller by notifying the server and clearing
  /// all stored credentials.
  Future<void> logout() async {
    try {
      await _apiClient.post(ApiEndpoints.logout);
    } catch (_) {
      // Best-effort logout on server; always clear local tokens.
    }
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

      final response = await _apiClient.post(
        ApiEndpoints.refreshToken,
        data: {
          'refresh_token': currentRefreshToken,
        },
      );

      final data = response.data as Map<String, dynamic>;
      final newAccessToken = data['access_token'] as String;
      final expiresIn = data['expires_in'] as int?;
      final newExpiry = DateTime.now().add(
        Duration(seconds: expiresIn ?? 86400),
      );

      await _secureStorage.setAccessToken(newAccessToken);
      await _secureStorage.setTokenExpiry(newExpiry);

      // Update refresh token if the server rotated it.
      if (data.containsKey('refresh_token')) {
        await _secureStorage.setRefreshToken(
          data['refresh_token'] as String,
        );
      }

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
