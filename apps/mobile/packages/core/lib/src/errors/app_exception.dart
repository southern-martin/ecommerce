/// Base exception class for all application-level errors.
///
/// Provides a consistent error structure with a human-readable [message],
/// an optional machine-readable [code], and an optional [statusCode].
class AppException implements Exception {
  /// Human-readable error message.
  final String message;

  /// Optional machine-readable error code (e.g., `'INVALID_CREDENTIALS'`).
  final String? code;

  /// Optional HTTP status code associated with this error.
  final int? statusCode;

  /// Optional original error / stacktrace for debugging.
  final dynamic originalError;

  const AppException({
    required this.message,
    this.code,
    this.statusCode,
    this.originalError,
  });

  @override
  String toString() =>
      'AppException(message: $message, code: $code, statusCode: $statusCode)';
}

/// Thrown when a network request fails due to connectivity issues,
/// timeouts, or server errors.
class NetworkException extends AppException {
  /// Whether the failure was caused by lack of internet connectivity.
  final bool isConnectionError;

  /// Whether the failure was caused by a request timeout.
  final bool isTimeout;

  const NetworkException({
    required super.message,
    super.code,
    super.statusCode,
    super.originalError,
    this.isConnectionError = false,
    this.isTimeout = false,
  });

  /// Factory for connection errors (no internet).
  factory NetworkException.noConnection() {
    return const NetworkException(
      message: 'No internet connection. Please check your network settings.',
      code: 'NO_CONNECTION',
      isConnectionError: true,
    );
  }

  /// Factory for timeout errors.
  factory NetworkException.timeout() {
    return const NetworkException(
      message: 'The request timed out. Please try again.',
      code: 'TIMEOUT',
      isTimeout: true,
    );
  }

  /// Factory for generic server errors.
  factory NetworkException.serverError({int? statusCode, String? message}) {
    return NetworkException(
      message: message ?? 'An unexpected server error occurred. Please try again later.',
      code: 'SERVER_ERROR',
      statusCode: statusCode ?? 500,
    );
  }

  @override
  String toString() =>
      'NetworkException(message: $message, code: $code, statusCode: $statusCode, '
      'isConnectionError: $isConnectionError, isTimeout: $isTimeout)';
}

/// Thrown when the user is not authenticated or the session has expired.
class UnauthorizedException extends AppException {
  const UnauthorizedException({
    super.message = 'Your session has expired. Please sign in again.',
    super.code = 'UNAUTHORIZED',
    super.statusCode = 401,
    super.originalError,
  });

  @override
  String toString() => 'UnauthorizedException(message: $message)';
}

/// Thrown when a requested resource cannot be found.
class NotFoundException extends AppException {
  /// The type of resource that was not found (e.g., `'Product'`, `'Order'`).
  final String? resourceType;

  /// The identifier of the resource that was not found.
  final String? resourceId;

  const NotFoundException({
    super.message = 'The requested resource was not found.',
    super.code = 'NOT_FOUND',
    super.statusCode = 404,
    super.originalError,
    this.resourceType,
    this.resourceId,
  });

  /// Factory for a specific missing resource.
  factory NotFoundException.forResource(String type, String id) {
    return NotFoundException(
      message: '$type with ID "$id" was not found.',
      resourceType: type,
      resourceId: id,
    );
  }

  @override
  String toString() =>
      'NotFoundException(message: $message, resourceType: $resourceType, '
      'resourceId: $resourceId)';
}

/// Thrown when user input or request payload fails validation.
class ValidationException extends AppException {
  /// Map of field names to their respective error messages.
  final Map<String, List<String>> fieldErrors;

  const ValidationException({
    super.message = 'Please fix the errors below and try again.',
    super.code = 'VALIDATION_ERROR',
    super.statusCode = 422,
    super.originalError,
    this.fieldErrors = const {},
  });

  /// Returns `true` if there are field-level errors.
  bool get hasFieldErrors => fieldErrors.isNotEmpty;

  /// Returns the first error message for a given [field], or `null`.
  String? errorForField(String field) {
    final errors = fieldErrors[field];
    return (errors != null && errors.isNotEmpty) ? errors.first : null;
  }

  @override
  String toString() =>
      'ValidationException(message: $message, fieldErrors: $fieldErrors)';
}

/// Thrown when the user does not have permission to perform an action.
class ForbiddenException extends AppException {
  const ForbiddenException({
    super.message = 'You do not have permission to perform this action.',
    super.code = 'FORBIDDEN',
    super.statusCode = 403,
    super.originalError,
  });

  @override
  String toString() => 'ForbiddenException(message: $message)';
}

/// Thrown when there is a conflict, such as duplicate resource creation.
class ConflictException extends AppException {
  const ConflictException({
    super.message = 'A conflict occurred. The resource may already exist.',
    super.code = 'CONFLICT',
    super.statusCode = 409,
    super.originalError,
  });

  @override
  String toString() => 'ConflictException(message: $message)';
}
