import 'package:flutter_test/flutter_test.dart';
import 'package:ecommerce_api_client/ecommerce_api_client.dart';

void main() {
  group('ApiResponse', () {
    test('constructs with defaults', () {
      const response = ApiResponse<String>();
      expect(response.data, isNull);
      expect(response.message, isNull);
      expect(response.success, true);
      expect(response.meta, isNull);
    });

    test('fromJson deserializes data with factory function', () {
      final json = {
        'data': {'name': 'Widget', 'price': 29.99},
        'message': 'Product fetched',
        'success': true,
        'meta': {'responseTime': 42},
      };

      final response = ApiResponse<Map<String, dynamic>>.fromJson(
        json,
        (obj) => obj as Map<String, dynamic>,
      );

      expect(response.data, isNotNull);
      expect(response.data!['name'], 'Widget');
      expect(response.message, 'Product fetched');
      expect(response.success, true);
      expect(response.meta!['responseTime'], 42);
    });

    test('fromJson handles missing data key', () {
      final json = {
        'message': 'No data',
        'success': false,
      };

      final response = ApiResponse<String>.fromJson(
        json,
        (obj) => obj as String,
      );

      expect(response.data, isNull);
      expect(response.success, false);
    });

    test('fromJson defaults success to true when missing', () {
      final json = <String, dynamic>{
        'data': 'hello',
      };

      final response = ApiResponse<String>.fromJson(
        json,
        (obj) => obj as String,
      );

      expect(response.success, true);
    });

    test('toJson serializes with data factory', () {
      const response = ApiResponse<String>(
        data: 'test-data',
        message: 'ok',
        success: true,
      );

      final json = response.toJson((val) => val);
      expect(json['data'], 'test-data');
      expect(json['message'], 'ok');
      expect(json['success'], true);
    });

    test('toJson omits data when null', () {
      const response = ApiResponse<String>(
        message: 'empty',
        success: true,
      );

      final json = response.toJson((val) => val);
      expect(json.containsKey('data'), false);
    });

    test('toJson omits message when null', () {
      const response = ApiResponse<String>(data: 'test', success: true);
      final json = response.toJson((val) => val);
      expect(json.containsKey('message'), false);
    });

    test('toString includes key fields', () {
      const response = ApiResponse<String>(
        data: 'hello',
        message: 'ok',
        success: true,
      );
      final str = response.toString();
      expect(str, contains('success: true'));
      expect(str, contains('hello'));
    });
  });

  group('PaginatedResponse', () {
    test('constructs with required fields', () {
      const response = PaginatedResponse<String>(
        items: ['a', 'b', 'c'],
        total: 10,
        page: 1,
        pageSize: 3,
      );
      expect(response.items, ['a', 'b', 'c']);
      expect(response.total, 10);
      expect(response.page, 1);
      expect(response.pageSize, 3);
    });

    group('computed properties', () {
      const response = PaginatedResponse<String>(
        items: ['a', 'b'],
        total: 10,
        page: 2,
        pageSize: 3,
      );

      test('totalPages calculates correctly', () {
        // 10 / 3 = 3.33 → ceil → 4
        expect(response.totalPages, 4);
      });

      test('hasNextPage', () {
        expect(response.hasNextPage, true); // page 2 < 4

        const lastPage = PaginatedResponse<String>(
          items: ['x'],
          total: 10,
          page: 4,
          pageSize: 3,
        );
        expect(lastPage.hasNextPage, false);
      });

      test('hasPreviousPage', () {
        expect(response.hasPreviousPage, true); // page 2 > 1

        const firstPage = PaginatedResponse<String>(
          items: ['a'],
          total: 10,
          page: 1,
          pageSize: 3,
        );
        expect(firstPage.hasPreviousPage, false);
      });

      test('isFirstPage', () {
        const first = PaginatedResponse<String>(
          items: ['a'],
          total: 5,
          page: 1,
          pageSize: 2,
        );
        expect(first.isFirstPage, true);
        expect(response.isFirstPage, false);
      });

      test('isLastPage', () {
        const last = PaginatedResponse<String>(
          items: ['z'],
          total: 10,
          page: 4,
          pageSize: 3,
        );
        expect(last.isLastPage, true);
        expect(response.isLastPage, false);
      });
    });

    group('fromJson', () {
      test('parses with standard keys', () {
        final json = {
          'items': ['apple', 'banana'],
          'total': 50,
          'page': 3,
          'pageSize': 10,
        };

        final response = PaginatedResponse<String>.fromJson(
          json,
          (obj) => obj as String,
        );

        expect(response.items, ['apple', 'banana']);
        expect(response.total, 50);
        expect(response.page, 3);
        expect(response.pageSize, 10);
      });

      test('parses with alternative keys: data, totalCount, currentPage, perPage', () {
        final json = {
          'data': ['x', 'y'],
          'totalCount': 20,
          'currentPage': 2,
          'perPage': 5,
        };

        final response = PaginatedResponse<String>.fromJson(
          json,
          (obj) => obj as String,
        );

        expect(response.items, ['x', 'y']);
        expect(response.total, 20);
        expect(response.page, 2);
        expect(response.pageSize, 5);
      });

      test('parses with alternative key: limit', () {
        final json = {
          'items': ['a'],
          'total': 5,
          'page': 1,
          'limit': 25,
        };

        final response = PaginatedResponse<String>.fromJson(
          json,
          (obj) => obj as String,
        );

        expect(response.pageSize, 25);
      });

      test('handles missing list gracefully', () {
        final json = <String, dynamic>{
          'total': 0,
          'page': 1,
          'pageSize': 10,
        };

        final response = PaginatedResponse<String>.fromJson(
          json,
          (obj) => obj as String,
        );

        expect(response.items, isEmpty);
      });

      test('uses default values when keys are missing', () {
        final json = {
          'items': ['a'],
        };

        final response = PaginatedResponse<String>.fromJson(
          json,
          (obj) => obj as String,
        );

        expect(response.total, 0);
        expect(response.page, 1);
        expect(response.pageSize, 20);
      });
    });

    test('toJson serializes correctly', () {
      const response = PaginatedResponse<String>(
        items: ['a', 'b'],
        total: 10,
        page: 1,
        pageSize: 5,
      );

      final json = response.toJson((val) => val);
      expect(json['items'], ['a', 'b']);
      expect(json['total'], 10);
      expect(json['page'], 1);
      expect(json['pageSize'], 5);
    });

    test('toString includes pagination info', () {
      const response = PaginatedResponse<String>(
        items: ['a'],
        total: 10,
        page: 2,
        pageSize: 3,
      );
      final str = response.toString();
      expect(str, contains('2/4')); // page/totalPages
      expect(str, contains('total: 10'));
    });
  });
}
