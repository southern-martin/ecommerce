import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:ecommerce_core/ecommerce_core.dart';
import 'package:ecommerce_api_client/ecommerce_api_client.dart';

void main() {
  late ErrorInterceptor interceptor;
  late AppLogger logger;

  setUp(() {
    logger = AppLogger();
    interceptor = ErrorInterceptor(logger: logger);
  });

  /// Helper to invoke the interceptor's onError and capture the result.
  AppException? invokeOnError(DioException dioException) {
    AppException? captured;
    final handler = _MockErrorInterceptorHandler((err) {
      captured = err.error as AppException?;
    });
    interceptor.onError(dioException, handler);
    return captured;
  }

  DioException _makeDioException({
    required DioExceptionType type,
    int? statusCode,
    Map<String, dynamic>? responseData,
    String? message,
    dynamic error,
  }) {
    final options = RequestOptions(path: '/test');
    Response<dynamic>? response;
    if (statusCode != null) {
      response = Response(
        requestOptions: options,
        statusCode: statusCode,
        data: responseData,
      );
    }
    return DioException(
      requestOptions: options,
      type: type,
      response: response,
      message: message,
      error: error,
    );
  }

  group('Timeout errors', () {
    test('connectionTimeout maps to NetworkException.timeout', () {
      final result = invokeOnError(
        _makeDioException(type: DioExceptionType.connectionTimeout),
      );
      expect(result, isA<NetworkException>());
      final netEx = result as NetworkException;
      expect(netEx.isTimeout, true);
      expect(netEx.code, 'TIMEOUT');
    });

    test('sendTimeout maps to NetworkException.timeout', () {
      final result = invokeOnError(
        _makeDioException(type: DioExceptionType.sendTimeout),
      );
      expect(result, isA<NetworkException>());
      expect((result as NetworkException).isTimeout, true);
    });

    test('receiveTimeout maps to NetworkException.timeout', () {
      final result = invokeOnError(
        _makeDioException(type: DioExceptionType.receiveTimeout),
      );
      expect(result, isA<NetworkException>());
      expect((result as NetworkException).isTimeout, true);
    });
  });

  group('Connection errors', () {
    test('connectionError maps to NetworkException.noConnection', () {
      final result = invokeOnError(
        _makeDioException(type: DioExceptionType.connectionError),
      );
      expect(result, isA<NetworkException>());
      final netEx = result as NetworkException;
      expect(netEx.isConnectionError, true);
      expect(netEx.code, 'NO_CONNECTION');
    });
  });

  group('Bad response status codes', () {
    test('400 maps to AppException with BAD_REQUEST', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 400,
          responseData: {'message': 'Invalid input'},
        ),
      );
      expect(result, isA<AppException>());
      expect(result!.code, 'BAD_REQUEST');
      expect(result.statusCode, 400);
      expect(result.message, 'Invalid input');
    });

    test('401 maps to UnauthorizedException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 401,
          responseData: {'message': 'Token expired'},
        ),
      );
      expect(result, isA<UnauthorizedException>());
      expect(result!.message, 'Token expired');
    });

    test('401 uses default message when no server message', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 401,
        ),
      );
      expect(result, isA<UnauthorizedException>());
      expect(result!.message, contains('session has expired'));
    });

    test('403 maps to ForbiddenException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 403,
        ),
      );
      expect(result, isA<ForbiddenException>());
    });

    test('404 maps to NotFoundException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 404,
          responseData: {'message': 'Product not found'},
        ),
      );
      expect(result, isA<NotFoundException>());
      expect(result!.message, 'Product not found');
    });

    test('409 maps to ConflictException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 409,
        ),
      );
      expect(result, isA<ConflictException>());
    });

    test('422 maps to ValidationException with field errors', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 422,
          responseData: {
            'message': 'Validation failed',
            'errors': {
              'email': ['Email is required', 'Invalid format'],
              'password': ['Too short'],
            },
          },
        ),
      );
      expect(result, isA<ValidationException>());
      final valEx = result as ValidationException;
      expect(valEx.message, 'Validation failed');
      expect(valEx.hasFieldErrors, true);
      expect(valEx.errorForField('email'), 'Email is required');
      expect(valEx.fieldErrors['password'], ['Too short']);
    });

    test('422 handles non-list error values', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 422,
          responseData: {
            'errors': {
              'name': 'Required',
            },
          },
        ),
      );
      expect(result, isA<ValidationException>());
      final valEx = result as ValidationException;
      expect(valEx.fieldErrors['name'], ['Required']);
    });

    test('429 maps to rate-limited AppException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 429,
        ),
      );
      expect(result, isA<AppException>());
      expect(result!.code, 'RATE_LIMITED');
      expect(result.statusCode, 429);
    });

    test('500 maps to NetworkException.serverError', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 500,
          responseData: {'message': 'Internal error'},
        ),
      );
      expect(result, isA<NetworkException>());
      expect(result!.code, 'SERVER_ERROR');
      expect(result.message, 'Internal error');
    });

    test('502 maps to NetworkException.serverError', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 502,
        ),
      );
      expect(result, isA<NetworkException>());
      expect(result!.statusCode, 502);
    });

    test('503 maps to NetworkException.serverError', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 503,
        ),
      );
      expect(result, isA<NetworkException>());
      expect(result!.statusCode, 503);
    });

    test('504 maps to NetworkException.serverError', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 504,
        ),
      );
      expect(result, isA<NetworkException>());
      expect(result!.statusCode, 504);
    });

    test('unknown status code maps to generic AppException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 418,
          responseData: {'message': 'I am a teapot'},
        ),
      );
      expect(result, isA<AppException>());
      expect(result!.code, 'HTTP_418');
      expect(result.statusCode, 418);
      expect(result.message, 'I am a teapot');
    });

    test('extracts error key as fallback message', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.badResponse,
          statusCode: 400,
          responseData: {'error': 'Some error detail'},
        ),
      );
      expect(result!.message, 'Some error detail');
    });
  });

  group('Cancelled requests', () {
    test('cancel maps to AppException with CANCELLED code', () {
      final result = invokeOnError(
        _makeDioException(type: DioExceptionType.cancel),
      );
      expect(result, isA<AppException>());
      expect(result!.code, 'CANCELLED');
      expect(result.message, contains('cancelled'));
    });
  });

  group('Certificate errors', () {
    test('badCertificate maps to NetworkException', () {
      final result = invokeOnError(
        _makeDioException(type: DioExceptionType.badCertificate),
      );
      expect(result, isA<NetworkException>());
      expect(result!.code, 'BAD_CERTIFICATE');
    });
  });

  group('Unknown errors', () {
    test('SocketException maps to noConnection', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.unknown,
          error: Exception('SocketException: Connection refused'),
        ),
      );
      expect(result, isA<NetworkException>());
      expect((result as NetworkException).isConnectionError, true);
    });

    test('generic unknown error maps to NetworkException', () {
      final result = invokeOnError(
        _makeDioException(
          type: DioExceptionType.unknown,
          message: 'Something failed',
        ),
      );
      expect(result, isA<NetworkException>());
      expect(result!.code, 'UNKNOWN');
    });
  });
}

/// Minimal mock for ErrorInterceptorHandler to capture the error passed to next().
class _MockErrorInterceptorHandler extends ErrorInterceptorHandler {
  final void Function(DioException err) _onNext;

  _MockErrorInterceptorHandler(this._onNext);

  @override
  void next(DioException err) {
    _onNext(err);
  }

  @override
  void resolve(Response response) {}

  @override
  void reject(DioException err) {
    _onNext(err);
  }
}
