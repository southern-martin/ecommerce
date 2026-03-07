import 'package:test/test.dart';
import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';

void main() {
  group('AuthTokens', () {
    test('constructs with all required fields', () {
      final tokens = AuthTokens(
        accessToken: 'access-abc',
        refreshToken: 'refresh-xyz',
        expiresAt: DateTime(2026, 12, 31),
      );
      expect(tokens.accessToken, 'access-abc');
      expect(tokens.refreshToken, 'refresh-xyz');
      expect(tokens.expiresAt, DateTime(2026, 12, 31));
    });

    group('isExpired', () {
      test('returns true when expiresAt is in the past', () {
        final tokens = AuthTokens(
          accessToken: 'a',
          refreshToken: 'r',
          expiresAt: DateTime(2020, 1, 1),
        );
        expect(tokens.isExpired, true);
      });

      test('returns false when expiresAt is in the future', () {
        final tokens = AuthTokens(
          accessToken: 'a',
          refreshToken: 'r',
          expiresAt: DateTime.now().add(const Duration(hours: 1)),
        );
        expect(tokens.isExpired, false);
      });
    });

    group('expiresWithin', () {
      test('returns true when token expires within given duration', () {
        final tokens = AuthTokens(
          accessToken: 'a',
          refreshToken: 'r',
          expiresAt: DateTime.now().add(const Duration(minutes: 3)),
        );
        expect(tokens.expiresWithin(const Duration(minutes: 5)), true);
      });

      test('returns false when token does not expire within given duration', () {
        final tokens = AuthTokens(
          accessToken: 'a',
          refreshToken: 'r',
          expiresAt: DateTime.now().add(const Duration(hours: 2)),
        );
        expect(tokens.expiresWithin(const Duration(minutes: 5)), false);
      });
    });

    group('serialization', () {
      test('round-trips through fromJson/toJson', () {
        final json = {
          'accessToken': 'tok-123',
          'refreshToken': 'ref-456',
          'expiresAt': '2026-06-15T12:00:00.000',
        };

        final tokens = AuthTokens.fromJson(json);
        expect(tokens.accessToken, 'tok-123');
        expect(tokens.refreshToken, 'ref-456');
        expect(tokens.expiresAt, DateTime(2026, 6, 15, 12, 0));

        final backToJson = tokens.toJson();
        final roundTrip = AuthTokens.fromJson(backToJson);
        expect(roundTrip.accessToken, tokens.accessToken);
        expect(roundTrip.refreshToken, tokens.refreshToken);
        expect(roundTrip.expiresAt, tokens.expiresAt);
      });
    });

    test('copyWith overrides selected fields', () {
      final tokens = AuthTokens(
        accessToken: 'old-access',
        refreshToken: 'old-refresh',
        expiresAt: DateTime(2026, 1, 1),
      );
      final updated = tokens.copyWith(accessToken: 'new-access');
      expect(updated.accessToken, 'new-access');
      expect(updated.refreshToken, 'old-refresh');
    });

    test('toString does not expose tokens', () {
      final tokens = AuthTokens(
        accessToken: 'secret-token',
        refreshToken: 'secret-refresh',
        expiresAt: DateTime(2026, 1, 1),
      );
      final str = tokens.toString();
      expect(str, contains('expiresAt'));
      expect(str, isNot(contains('secret-token')));
    });
  });

  group('LoginRequest', () {
    test('constructs with email and password', () {
      const req = LoginRequest(email: 'user@test.com', password: 'pass123');
      expect(req.email, 'user@test.com');
      expect(req.password, 'pass123');
    });

    test('serialization round-trip', () {
      final json = {'email': 'a@b.com', 'password': 'secret'};
      final req = LoginRequest.fromJson(json);
      expect(req.email, 'a@b.com');

      final backToJson = req.toJson();
      expect(backToJson['email'], 'a@b.com');
      expect(backToJson['password'], 'secret');
    });

    test('toString shows email only', () {
      const req = LoginRequest(email: 'user@test.com', password: 'secret');
      expect(req.toString(), 'LoginRequest(email: user@test.com)');
    });
  });

  group('RegisterRequest', () {
    test('constructs with all fields', () {
      const req = RegisterRequest(
        email: 'new@test.com',
        password: 'P@ssw0rd!',
        firstName: 'Alice',
        lastName: 'Wonder',
      );
      expect(req.email, 'new@test.com');
      expect(req.firstName, 'Alice');
      expect(req.lastName, 'Wonder');
    });

    test('serialization round-trip', () {
      final json = {
        'email': 'bob@test.com',
        'password': 'Str0ng!Pass',
        'firstName': 'Bob',
        'lastName': 'Builder',
      };
      final req = RegisterRequest.fromJson(json);
      expect(req.firstName, 'Bob');
      expect(req.lastName, 'Builder');

      final backToJson = req.toJson();
      expect(backToJson['firstName'], 'Bob');
      expect(backToJson['lastName'], 'Builder');
    });

    test('toString shows email and name', () {
      const req = RegisterRequest(
        email: 'alice@test.com',
        password: 'secret',
        firstName: 'Alice',
        lastName: 'Wonder',
      );
      final str = req.toString();
      expect(str, contains('alice@test.com'));
      expect(str, contains('Alice Wonder'));
    });
  });
}
