import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Wrapper around [FlutterSecureStorage] for persisting sensitive data
/// such as authentication tokens and credentials.
///
/// All values are stored in the platform keychain / keystore and are
/// encrypted at rest.
class SecureStorage {
  static const String _accessTokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _tokenExpiryKey = 'token_expiry';
  static const String _userIdKey = 'user_id';

  final FlutterSecureStorage _storage;

  /// Creates a [SecureStorage] with an optional custom storage instance.
  SecureStorage({FlutterSecureStorage? storage})
      : _storage = storage ??
            const FlutterSecureStorage(
              aOptions: AndroidOptions(encryptedSharedPreferences: true),
              iOptions: IOSOptions(
                accessibility: KeychainAccessibility.first_unlock_this_device,
              ),
            );

  // ---------------------------------------------------------------------------
  // Access Token
  // ---------------------------------------------------------------------------

  /// Persists the JWT access token.
  Future<void> setAccessToken(String token) async {
    await _storage.write(key: _accessTokenKey, value: token);
  }

  /// Retrieves the stored access token, or `null` if none exists.
  Future<String?> getAccessToken() async {
    return _storage.read(key: _accessTokenKey);
  }

  /// Removes the stored access token.
  Future<void> deleteAccessToken() async {
    await _storage.delete(key: _accessTokenKey);
  }

  // ---------------------------------------------------------------------------
  // Refresh Token
  // ---------------------------------------------------------------------------

  /// Persists the refresh token.
  Future<void> setRefreshToken(String token) async {
    await _storage.write(key: _refreshTokenKey, value: token);
  }

  /// Retrieves the stored refresh token, or `null` if none exists.
  Future<String?> getRefreshToken() async {
    return _storage.read(key: _refreshTokenKey);
  }

  /// Removes the stored refresh token.
  Future<void> deleteRefreshToken() async {
    await _storage.delete(key: _refreshTokenKey);
  }

  // ---------------------------------------------------------------------------
  // Token Expiry
  // ---------------------------------------------------------------------------

  /// Persists the token expiry as an ISO-8601 string.
  Future<void> setTokenExpiry(DateTime expiry) async {
    await _storage.write(
      key: _tokenExpiryKey,
      value: expiry.toIso8601String(),
    );
  }

  /// Retrieves the stored token expiry, or `null` if none exists.
  Future<DateTime?> getTokenExpiry() async {
    final value = await _storage.read(key: _tokenExpiryKey);
    if (value == null) return null;
    return DateTime.tryParse(value);
  }

  /// Returns `true` if the stored token has expired.
  Future<bool> isTokenExpired() async {
    final expiry = await getTokenExpiry();
    if (expiry == null) return true;
    return DateTime.now().isAfter(expiry);
  }

  // ---------------------------------------------------------------------------
  // User ID
  // ---------------------------------------------------------------------------

  /// Persists the current user ID.
  Future<void> setUserId(String userId) async {
    await _storage.write(key: _userIdKey, value: userId);
  }

  /// Retrieves the stored user ID, or `null` if none exists.
  Future<String?> getUserId() async {
    return _storage.read(key: _userIdKey);
  }

  // ---------------------------------------------------------------------------
  // Generic helpers
  // ---------------------------------------------------------------------------

  /// Writes an arbitrary key-value pair to secure storage.
  Future<void> write({required String key, required String value}) async {
    await _storage.write(key: key, value: value);
  }

  /// Reads an arbitrary value from secure storage.
  Future<String?> read({required String key}) async {
    return _storage.read(key: key);
  }

  /// Deletes an arbitrary key from secure storage.
  Future<void> delete({required String key}) async {
    await _storage.delete(key: key);
  }

  /// Clears **all** values from secure storage. Use with caution.
  Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
