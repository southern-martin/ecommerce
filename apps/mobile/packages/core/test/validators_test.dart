import 'package:flutter_test/flutter_test.dart';
import 'package:ecommerce_core/ecommerce_core.dart';

void main() {
  group('Validators.email', () {
    test('returns null for valid email', () {
      expect(Validators.email('user@example.com'), isNull);
      expect(Validators.email('name+tag@domain.co.uk'), isNull);
      expect(Validators.email('  user@example.com  '), isNull); // trimmed
    });

    test('returns error for null', () {
      expect(Validators.email(null), isNotNull);
    });

    test('returns error for empty string', () {
      expect(Validators.email(''), isNotNull);
      expect(Validators.email('   '), isNotNull);
    });

    test('returns error for invalid emails', () {
      expect(Validators.email('not-an-email'), isNotNull);
      expect(Validators.email('@domain.com'), isNotNull);
      expect(Validators.email('user@'), isNotNull);
    });
  });

  group('Validators.password', () {
    test('returns null for valid password', () {
      expect(Validators.password('Str0ng!Pass'), isNull);
      expect(Validators.password('Ab1!defgh'), isNull);
    });

    test('returns error for null', () {
      expect(Validators.password(null), isNotNull);
    });

    test('returns error for empty string', () {
      expect(Validators.password(''), isNotNull);
    });

    test('returns error for short password', () {
      final result = Validators.password('Ab1!xyz');
      expect(result, isNotNull);
      expect(result, contains('8 characters'));
    });

    test('returns error for missing uppercase', () {
      final result = Validators.password('abcdefg1!');
      expect(result, isNotNull);
      expect(result, contains('uppercase'));
    });

    test('returns error for missing lowercase', () {
      final result = Validators.password('ABCDEFG1!');
      expect(result, isNotNull);
      expect(result, contains('lowercase'));
    });

    test('returns error for missing digit', () {
      final result = Validators.password('Abcdefgh!');
      expect(result, isNotNull);
      expect(result, contains('number'));
    });

    test('returns error for missing special character', () {
      final result = Validators.password('Abcdefg1');
      expect(result, isNotNull);
      expect(result, contains('special'));
    });
  });

  group('Validators.confirmPassword', () {
    test('returns null when passwords match', () {
      expect(Validators.confirmPassword('abc123', 'abc123'), isNull);
    });

    test('returns error when null', () {
      expect(Validators.confirmPassword(null, 'abc'), isNotNull);
    });

    test('returns error when empty', () {
      expect(Validators.confirmPassword('', 'abc'), isNotNull);
    });

    test('returns error when passwords differ', () {
      final result = Validators.confirmPassword('abc', 'xyz');
      expect(result, isNotNull);
      expect(result, contains('do not match'));
    });
  });

  group('Validators.phone', () {
    test('returns null for valid phone numbers', () {
      expect(Validators.phone('+15551234567'), isNull);
      expect(Validators.phone('555-123-4567'), isNull);
      expect(Validators.phone('(555) 123 4567'), isNull);
      expect(Validators.phone('1234567'), isNull); // 7 digits min
    });

    test('returns error for null', () {
      expect(Validators.phone(null), isNotNull);
    });

    test('returns error for empty', () {
      expect(Validators.phone(''), isNotNull);
      expect(Validators.phone('   '), isNotNull);
    });

    test('returns error for too few digits', () {
      expect(Validators.phone('123456'), isNotNull); // 6 digits
    });

    test('returns error for non-digit characters', () {
      expect(Validators.phone('abc-def-ghij'), isNotNull);
    });
  });

  group('Validators.required', () {
    test('returns null for non-empty string', () {
      expect(Validators.required('hello'), isNull);
    });

    test('returns error for null', () {
      expect(Validators.required(null), isNotNull);
    });

    test('returns error for empty or whitespace', () {
      expect(Validators.required(''), isNotNull);
      expect(Validators.required('   '), isNotNull);
    });

    test('uses custom field name in message', () {
      final result = Validators.required(null, 'Email');
      expect(result, 'Email is required');
    });

    test('uses default field name when not specified', () {
      final result = Validators.required(null);
      expect(result, 'This field is required');
    });
  });

  group('Validators.minLength', () {
    test('returns null when meets minimum', () {
      expect(Validators.minLength('hello', 3), isNull);
      expect(Validators.minLength('abc', 3), isNull);
    });

    test('returns error when too short', () {
      final result = Validators.minLength('ab', 3);
      expect(result, isNotNull);
      expect(result, contains('3 characters'));
    });

    test('returns error for null', () {
      expect(Validators.minLength(null, 1), isNotNull);
    });

    test('uses custom field name', () {
      final result = Validators.minLength('a', 5, 'Username');
      expect(result, contains('Username'));
    });
  });

  group('Validators.maxLength', () {
    test('returns null when within limit', () {
      expect(Validators.maxLength('abc', 5), isNull);
      expect(Validators.maxLength('hello', 5), isNull);
    });

    test('returns null for null value', () {
      expect(Validators.maxLength(null, 5), isNull);
    });

    test('returns error when too long', () {
      final result = Validators.maxLength('toolongtext', 5);
      expect(result, isNotNull);
      expect(result, contains('5 characters'));
    });

    test('uses custom field name', () {
      final result = Validators.maxLength('toolong', 3, 'Name');
      expect(result, contains('Name'));
    });
  });

  group('Validators.url', () {
    test('returns null for valid URLs', () {
      expect(Validators.url('https://example.com'), isNull);
      expect(Validators.url('http://localhost:3000'), isNull);
      expect(Validators.url('ftp://files.example.com/doc.pdf'), isNull);
    });

    test('returns error for null or empty', () {
      expect(Validators.url(null), isNotNull);
      expect(Validators.url(''), isNotNull);
      expect(Validators.url('   '), isNotNull);
    });

    test('returns error for invalid URL', () {
      expect(Validators.url('not-a-url'), isNotNull);
      expect(Validators.url('example.com'), isNotNull); // no scheme
    });
  });

  group('Validators.numericOnly', () {
    test('returns null for digits-only string', () {
      expect(Validators.numericOnly('12345'), isNull);
      expect(Validators.numericOnly('0'), isNull);
    });

    test('returns error for null or empty', () {
      expect(Validators.numericOnly(null), isNotNull);
      expect(Validators.numericOnly(''), isNotNull);
    });

    test('returns error for non-digit characters', () {
      expect(Validators.numericOnly('12.5'), isNotNull);
      expect(Validators.numericOnly('abc'), isNotNull);
      expect(Validators.numericOnly('-10'), isNotNull);
    });

    test('uses custom field name', () {
      final result = Validators.numericOnly('abc', 'Zip Code');
      expect(result, contains('Zip Code'));
    });
  });
}
