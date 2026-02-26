import 'dart:async';

import 'package:dio/dio.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Dio interceptor that handles JWT authentication.
///
/// - **onRequest**: Reads the access token from [SecureStorage] and attaches
///   it as a `Bearer` token in the `Authorization` header.
/// - **onError**: When a 401 response is received, attempts to refresh the
///   token using the stored refresh token. If the refresh succeeds, the
///   original request is retried with the new token. If the refresh fails,
///   the user is considered unauthenticated and all tokens are cleared.
class AuthInterceptor extends Interceptor {
  final Dio _dio;
  final SecureStorage _secureStorage;
  final AppLogger _logger;

  /// Whether a token refresh is currently in progress. Prevents concurrent
  /// refresh requests when multiple 401s arrive simultaneously.
  bool _isRefreshing = false;

  /// Completer that resolves when the ongoing token refresh finishes.
  /// Queued requests await this completer before retrying.
  Completer<String?>? _refreshCompleter;

  /// Paths that should never have an Authorization header attached
  /// (e.g., login and register endpoints).
  static const _publicPaths = <String>{
    ApiEndpoints.login,
    ApiEndpoints.register,
    ApiEndpoints.refreshToken,
    ApiEndpoints.forgotPassword,
  };

  AuthInterceptor({
    required Dio dio,
    required SecureStorage secureStorage,
    required AppLogger logger,
  })  : _dio = dio,
        _secureStorage = secureStorage,
        _logger = logger;

  // ---------------------------------------------------------------------------
  // onRequest — attach Bearer token
  // ---------------------------------------------------------------------------

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final path = options.path;

    // Skip auth header for public endpoints.
    if (_publicPaths.any((p) => path.contains(p))) {
      return handler.next(options);
    }

    final token = await _secureStorage.getAccessToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    handler.next(options);
  }

  // ---------------------------------------------------------------------------
  // onError — handle 401 with token refresh
  // ---------------------------------------------------------------------------

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode != 401) {
      return handler.next(err);
    }

    // Don't attempt to refresh if the failing request IS the refresh call.
    final requestPath = err.requestOptions.path;
    if (requestPath.contains(ApiEndpoints.refreshToken)) {
      _logger.warning('Refresh token request itself returned 401. Logging out.',
          tag: 'AuthInterceptor');
      await _clearTokens();
      return handler.next(err);
    }

    _logger.info('Received 401, attempting token refresh.',
        tag: 'AuthInterceptor');

    try {
      final newToken = await _refreshAccessToken();

      if (newToken != null) {
        // Retry the original request with the new token.
        final retryOptions = err.requestOptions;
        retryOptions.headers['Authorization'] = 'Bearer $newToken';

        final response = await _dio.fetch(retryOptions);
        return handler.resolve(response);
      } else {
        // Refresh failed — propagate the original 401.
        return handler.next(err);
      }
    } catch (e) {
      _logger.error('Token refresh failed.', tag: 'AuthInterceptor', error: e);
      return handler.next(err);
    }
  }

  // ---------------------------------------------------------------------------
  // Token refresh logic
  // ---------------------------------------------------------------------------

  /// Attempts to refresh the access token. If a refresh is already in progress,
  /// subsequent callers will wait for the same result instead of issuing
  /// duplicate refresh requests.
  Future<String?> _refreshAccessToken() async {
    if (_isRefreshing) {
      // Wait for the in-flight refresh to complete.
      return _refreshCompleter?.future;
    }

    _isRefreshing = true;
    _refreshCompleter = Completer<String?>();

    try {
      final refreshToken = await _secureStorage.getRefreshToken();
      if (refreshToken == null || refreshToken.isEmpty) {
        _logger.warning('No refresh token available.', tag: 'AuthInterceptor');
        await _clearTokens();
        _refreshCompleter?.complete(null);
        return null;
      }

      // Use a separate Dio instance to avoid interceptor recursion.
      final refreshDio = Dio(BaseOptions(baseUrl: _dio.options.baseUrl));
      final response = await refreshDio.post(
        ApiEndpoints.refreshToken,
        data: {'refreshToken': refreshToken},
      );

      if (response.statusCode == 200 && response.data != null) {
        final newAccessToken = response.data['accessToken'] as String?;
        final newRefreshToken = response.data['refreshToken'] as String?;
        final expiresIn = response.data['expiresIn'] as int?;

        if (newAccessToken != null) {
          await _secureStorage.setAccessToken(newAccessToken);

          if (newRefreshToken != null) {
            await _secureStorage.setRefreshToken(newRefreshToken);
          }

          if (expiresIn != null) {
            await _secureStorage.setTokenExpiry(
              DateTime.now().add(Duration(seconds: expiresIn)),
            );
          }

          _logger.info('Token refreshed successfully.',
              tag: 'AuthInterceptor');
          _refreshCompleter?.complete(newAccessToken);
          return newAccessToken;
        }
      }

      // Unexpected response shape.
      _logger.warning('Unexpected refresh response.', tag: 'AuthInterceptor');
      await _clearTokens();
      _refreshCompleter?.complete(null);
      return null;
    } catch (e) {
      _logger.error('Error during token refresh.',
          tag: 'AuthInterceptor', error: e);
      await _clearTokens();
      _refreshCompleter?.complete(null);
      return null;
    } finally {
      _isRefreshing = false;
      _refreshCompleter = null;
    }
  }

  /// Removes all stored tokens.
  Future<void> _clearTokens() async {
    await _secureStorage.deleteAccessToken();
    await _secureStorage.deleteRefreshToken();
  }
}
