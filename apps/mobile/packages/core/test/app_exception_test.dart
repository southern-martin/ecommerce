import 'package:flutter_test/flutter_test.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

void main() {
  group('AppException', () {
    test('constructs with required message', () {
      const ex = AppException(message: 'Something went wrong');
      expect(ex.message, 'Something went wrong');
      expect(ex.code, isNull);
      expect(ex.statusCode, isNull);
      expect(ex.originalError, isNull);
    });

    test('constructs with all fields', () {
      const ex = AppException(
        message: 'Bad input',
        code: 'BAD_INPUT',
        statusCode: 400,
        originalError: 'raw error',
      );
      expect(ex.message, 'Bad input');
      expect(ex.code, 'BAD_INPUT');
      expect(ex.statusCode, 400);
      expect(ex.originalError, 'raw error');
    });

    test('toString includes all fields', () {
      const ex = AppException(message: 'Test', code: 'CODE', statusCode: 500);
      final str = ex.toString();
      expect(str, contains('Test'));
      expect(str, contains('CODE'));
      expect(str, contains('500'));
    });

    test('implements Exception', () {
      const ex = AppException(message: 'error');
      expect(ex, isA<Exception>());
    });
  });

  group('NetworkException', () {
    test('noConnection factory', () {
      final ex = NetworkException.noConnection();
      expect(ex.isConnectionError, true);
      expect(ex.isTimeout, false);
      expect(ex.code, 'NO_CONNECTION');
      expect(ex.message, contains('No internet'));
    });

    test('timeout factory', () {
      final ex = NetworkException.timeout();
      expect(ex.isTimeout, true);
      expect(ex.isConnectionError, false);
      expect(ex.code, 'TIMEOUT');
      expect(ex.message, contains('timed out'));
    });

    test('serverError factory with defaults', () {
      final ex = NetworkException.serverError();
      expect(ex.code, 'SERVER_ERROR');
      expect(ex.statusCode, 500);
      expect(ex.message, contains('unexpected server error'));
    });

    test('serverError factory with custom values', () {
      final ex = NetworkException.serverError(
        statusCode: 503,
        message: 'Service unavailable',
      );
      expect(ex.statusCode, 503);
      expect(ex.message, 'Service unavailable');
    });

    test('extends AppException', () {
      final ex = NetworkException.noConnection();
      expect(ex, isA<AppException>());
    });

    test('toString includes connection and timeout flags', () {
      final ex = NetworkException.noConnection();
      final str = ex.toString();
      expect(str, contains('isConnectionError: true'));
      expect(str, contains('isTimeout: false'));
    });
  });

  group('UnauthorizedException', () {
    test('default message', () {
      const ex = UnauthorizedException();
      expect(ex.message, contains('session has expired'));
      expect(ex.code, 'UNAUTHORIZED');
      expect(ex.statusCode, 401);
    });

    test('custom message', () {
      const ex = UnauthorizedException(message: 'Token invalid');
      expect(ex.message, 'Token invalid');
    });

    test('extends AppException', () {
      const ex = UnauthorizedException();
      expect(ex, isA<AppException>());
    });
  });

  group('NotFoundException', () {
    test('default message', () {
      const ex = NotFoundException();
      expect(ex.message, contains('not found'));
      expect(ex.code, 'NOT_FOUND');
      expect(ex.statusCode, 404);
      expect(ex.resourceType, isNull);
      expect(ex.resourceId, isNull);
    });

    test('forResource factory', () {
      final ex = NotFoundException.forResource('Product', 'prod-123');
      expect(ex.message, 'Product with ID "prod-123" was not found.');
      expect(ex.resourceType, 'Product');
      expect(ex.resourceId, 'prod-123');
    });

    test('toString includes resource info', () {
      final ex = NotFoundException.forResource('Order', 'ord-1');
      final str = ex.toString();
      expect(str, contains('Order'));
      expect(str, contains('ord-1'));
    });
  });

  group('ValidationException', () {
    test('default message and empty fieldErrors', () {
      const ex = ValidationException();
      expect(ex.message, contains('fix the errors'));
      expect(ex.code, 'VALIDATION_ERROR');
      expect(ex.statusCode, 422);
      expect(ex.fieldErrors, isEmpty);
      expect(ex.hasFieldErrors, false);
    });

    test('hasFieldErrors returns true when errors present', () {
      const ex = ValidationException(
        fieldErrors: {
          'email': ['Email is required'],
          'password': ['Too short', 'Missing digit'],
        },
      );
      expect(ex.hasFieldErrors, true);
    });

    test('errorForField returns first error message', () {
      const ex = ValidationException(
        fieldErrors: {
          'email': ['Email is required', 'Invalid format'],
        },
      );
      expect(ex.errorForField('email'), 'Email is required');
    });

    test('errorForField returns null for unknown field', () {
      const ex = ValidationException(
        fieldErrors: {
          'email': ['Required'],
        },
      );
      expect(ex.errorForField('name'), isNull);
    });

    test('errorForField returns null for empty error list', () {
      const ex = ValidationException(
        fieldErrors: {
          'email': [],
        },
      );
      expect(ex.errorForField('email'), isNull);
    });
  });

  group('ForbiddenException', () {
    test('default values', () {
      const ex = ForbiddenException();
      expect(ex.message, contains('permission'));
      expect(ex.code, 'FORBIDDEN');
      expect(ex.statusCode, 403);
    });

    test('extends AppException', () {
      const ex = ForbiddenException();
      expect(ex, isA<AppException>());
    });
  });

  group('ConflictException', () {
    test('default values', () {
      const ex = ConflictException();
      expect(ex.message, contains('conflict'));
      expect(ex.code, 'CONFLICT');
      expect(ex.statusCode, 409);
    });

    test('extends AppException', () {
      const ex = ConflictException();
      expect(ex, isA<AppException>());
    });
  });
}
