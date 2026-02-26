import 'package:dio/dio.dart';
import 'package:ecommerce_core/ecommerce_core.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';

import 'interceptors/auth_interceptor.dart';
import 'interceptors/error_interceptor.dart';

/// Pre-configured Dio HTTP client with authentication, error handling,
/// logging, and retry capabilities.
///
/// Usage:
/// ```dart
/// final client = ApiClient(
///   baseUrl: 'https://api.example.com',
///   secureStorage: secureStorage,
/// );
///
/// final response = await client.get('/api/v1/products');
/// ```
class ApiClient {
  late final Dio _dio;
  final SecureStorage _secureStorage;
  final AppLogger _logger;

  /// The base URL for all API requests.
  final String baseUrl;

  /// Creates an [ApiClient] configured with interceptors for auth, error
  /// mapping, retry, and logging.
  ///
  /// [baseUrl] is the root URL of the API server.
  /// [secureStorage] is used by the auth interceptor to read/write tokens.
  /// [logger] is used for debug logging; defaults to a new [AppLogger].
  /// [enableLogging] controls whether HTTP request/response logging is active.
  ApiClient({
    required this.baseUrl,
    required SecureStorage secureStorage,
    AppLogger? logger,
    bool enableLogging = true,
  })  : _secureStorage = secureStorage,
        _logger = logger ?? AppLogger() {
    _dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(
          milliseconds: AppConstants.connectionTimeout,
        ),
        receiveTimeout: const Duration(
          milliseconds: AppConstants.receiveTimeout,
        ),
        sendTimeout: const Duration(
          milliseconds: AppConstants.sendTimeout,
        ),
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );

    // Order matters: auth first, then error mapping, then logging.
    _dio.interceptors.add(
      AuthInterceptor(
        dio: _dio,
        secureStorage: _secureStorage,
        logger: _logger,
      ),
    );

    _dio.interceptors.add(
      ErrorInterceptor(logger: _logger),
    );

    if (enableLogging) {
      _dio.interceptors.add(
        PrettyDioLogger(
          requestHeader: true,
          requestBody: true,
          responseHeader: false,
          responseBody: true,
          error: true,
          compact: true,
        ),
      );
    }
  }

  /// The underlying [Dio] instance for advanced use cases.
  Dio get dio => _dio;

  // ---------------------------------------------------------------------------
  // HTTP methods
  // ---------------------------------------------------------------------------

  /// Sends a GET request to [path] with optional [queryParameters].
  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) {
    return _dio.get<T>(
      path,
      queryParameters: queryParameters,
      options: options,
      cancelToken: cancelToken,
    );
  }

  /// Sends a POST request to [path] with optional [data] body.
  Future<Response<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) {
    return _dio.post<T>(
      path,
      data: data,
      queryParameters: queryParameters,
      options: options,
      cancelToken: cancelToken,
    );
  }

  /// Sends a PUT request to [path] with optional [data] body.
  Future<Response<T>> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) {
    return _dio.put<T>(
      path,
      data: data,
      queryParameters: queryParameters,
      options: options,
      cancelToken: cancelToken,
    );
  }

  /// Sends a PATCH request to [path] with optional [data] body.
  Future<Response<T>> patch<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) {
    return _dio.patch<T>(
      path,
      data: data,
      queryParameters: queryParameters,
      options: options,
      cancelToken: cancelToken,
    );
  }

  /// Sends a DELETE request to [path].
  Future<Response<T>> delete<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
  }) {
    return _dio.delete<T>(
      path,
      data: data,
      queryParameters: queryParameters,
      options: options,
      cancelToken: cancelToken,
    );
  }

  /// Uploads a file using multipart form data.
  Future<Response<T>> upload<T>(
    String path, {
    required FormData formData,
    Options? options,
    CancelToken? cancelToken,
    void Function(int, int)? onSendProgress,
  }) {
    return _dio.post<T>(
      path,
      data: formData,
      options: options,
      cancelToken: cancelToken,
      onSendProgress: onSendProgress,
    );
  }

  /// Downloads a file from [urlPath] and saves it to [savePath].
  Future<Response> download(
    String urlPath,
    String savePath, {
    Map<String, dynamic>? queryParameters,
    CancelToken? cancelToken,
    void Function(int, int)? onReceiveProgress,
  }) {
    return _dio.download(
      urlPath,
      savePath,
      queryParameters: queryParameters,
      cancelToken: cancelToken,
      onReceiveProgress: onReceiveProgress,
    );
  }
}
