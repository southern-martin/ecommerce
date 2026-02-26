import 'package:dio/dio.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

/// Dio interceptor that maps [DioException] instances to strongly-typed
/// [AppException] subtypes, providing a consistent error model for the
/// rest of the application.
class ErrorInterceptor extends Interceptor {
  final AppLogger _logger;

  ErrorInterceptor({required AppLogger logger}) : _logger = logger;

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final appException = _mapDioException(err);

    _logger.error(
      'API Error: ${appException.message}',
      tag: 'ErrorInterceptor',
      error: err,
    );

    handler.next(
      DioException(
        requestOptions: err.requestOptions,
        response: err.response,
        type: err.type,
        error: appException,
        message: appException.message,
      ),
    );
  }

  /// Converts a [DioException] into the most appropriate [AppException]
  /// subtype based on the error type and HTTP status code.
  AppException _mapDioException(DioException exception) {
    switch (exception.type) {
      // -----------------------------------------------------------------------
      // Timeout errors
      // -----------------------------------------------------------------------
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return NetworkException.timeout();

      // -----------------------------------------------------------------------
      // Connection error (no internet, DNS failure, etc.)
      // -----------------------------------------------------------------------
      case DioExceptionType.connectionError:
        return NetworkException.noConnection();

      // -----------------------------------------------------------------------
      // Server responded with an error status code
      // -----------------------------------------------------------------------
      case DioExceptionType.badResponse:
        return _mapStatusCode(exception);

      // -----------------------------------------------------------------------
      // Request was cancelled
      // -----------------------------------------------------------------------
      case DioExceptionType.cancel:
        return const AppException(
          message: 'Request was cancelled.',
          code: 'CANCELLED',
        );

      // -----------------------------------------------------------------------
      // Certificate / other errors
      // -----------------------------------------------------------------------
      case DioExceptionType.badCertificate:
        return const NetworkException(
          message: 'Could not verify the server certificate.',
          code: 'BAD_CERTIFICATE',
        );

      case DioExceptionType.unknown:
      default:
        // Check for common connection-related exceptions wrapped in unknown.
        if (exception.error != null &&
            exception.error.toString().contains('SocketException')) {
          return NetworkException.noConnection();
        }
        return NetworkException(
          message: exception.message ?? 'An unexpected error occurred.',
          code: 'UNKNOWN',
          originalError: exception,
        );
    }
  }

  /// Maps an HTTP status code to the appropriate [AppException] subtype.
  AppException _mapStatusCode(DioException exception) {
    final statusCode = exception.response?.statusCode;
    final responseData = exception.response?.data;

    // Attempt to extract a server-provided error message.
    String? serverMessage;
    Map<String, List<String>> fieldErrors = {};

    if (responseData is Map<String, dynamic>) {
      serverMessage = responseData['message'] as String? ??
          responseData['error'] as String?;

      // Parse field-level validation errors if present.
      final errors = responseData['errors'];
      if (errors is Map<String, dynamic>) {
        fieldErrors = errors.map((key, value) {
          if (value is List) {
            return MapEntry(key, value.cast<String>());
          }
          return MapEntry(key, [value.toString()]);
        });
      }
    }

    switch (statusCode) {
      case 400:
        return AppException(
          message: serverMessage ?? 'Bad request. Please check your input.',
          code: 'BAD_REQUEST',
          statusCode: 400,
          originalError: exception,
        );

      case 401:
        return UnauthorizedException(
          message: serverMessage ??
              'Your session has expired. Please sign in again.',
          originalError: exception,
        );

      case 403:
        return ForbiddenException(
          message: serverMessage ??
              'You do not have permission to perform this action.',
          originalError: exception,
        );

      case 404:
        return NotFoundException(
          message: serverMessage ?? 'The requested resource was not found.',
          originalError: exception,
        );

      case 409:
        return ConflictException(
          message: serverMessage ??
              'A conflict occurred. The resource may already exist.',
          originalError: exception,
        );

      case 422:
        return ValidationException(
          message: serverMessage ??
              'Please fix the errors below and try again.',
          fieldErrors: fieldErrors,
          originalError: exception,
        );

      case 429:
        return const AppException(
          message: 'Too many requests. Please wait and try again.',
          code: 'RATE_LIMITED',
          statusCode: 429,
        );

      case 500:
      case 502:
      case 503:
      case 504:
        return NetworkException.serverError(
          statusCode: statusCode,
          message: serverMessage,
        );

      default:
        return AppException(
          message: serverMessage ?? 'An unexpected error occurred.',
          code: 'HTTP_$statusCode',
          statusCode: statusCode,
          originalError: exception,
        );
    }
  }
}
