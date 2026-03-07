import 'package:test/test.dart';
import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';

void main() {
  group('Order', () {
    late Order order;

    setUp(() {
      order = Order(
        id: 'order-1',
        orderNumber: 'ORD-2026-0001',
        status: OrderStatus.pending,
        items: const [
          OrderItem(
            id: 'item-1',
            productId: 'prod-1',
            name: 'Widget',
            quantity: 2,
            priceCents: 1500,
          ),
          OrderItem(
            id: 'item-2',
            productId: 'prod-2',
            name: 'Gadget',
            quantity: 1,
            priceCents: 3000,
          ),
        ],
        subtotalCents: 6000,
        taxCents: 600,
        shippingCents: 500,
        totalCents: 7100,
        createdAt: DateTime(2026, 3, 1),
      );
    });

    test('constructs with required fields and correct defaults', () {
      final minimal = Order(
        id: 'order-2',
        orderNumber: 'ORD-2026-0002',
        status: OrderStatus.pending,
        subtotalCents: 1000,
        totalCents: 1000,
        createdAt: DateTime(2026, 3, 1),
      );
      expect(minimal.items, isEmpty);
      expect(minimal.taxCents, 0);
      expect(minimal.shippingCents, 0);
      expect(minimal.shippingAddress, isNull);
    });

    group('itemCount', () {
      test('sums all item quantities', () {
        expect(order.itemCount, 3); // 2 + 1
      });

      test('returns 0 for empty order', () {
        final empty = order.copyWith(items: []);
        expect(empty.itemCount, 0);
      });
    });

    group('isCancellable', () {
      test('true when pending', () {
        expect(order.isCancellable, true);
      });

      test('true when confirmed', () {
        final confirmed = order.copyWith(status: OrderStatus.confirmed);
        expect(confirmed.isCancellable, true);
      });

      test('false when processing', () {
        final processing = order.copyWith(status: OrderStatus.processing);
        expect(processing.isCancellable, false);
      });

      test('false when shipped', () {
        final shipped = order.copyWith(status: OrderStatus.shipped);
        expect(shipped.isCancellable, false);
      });

      test('false when delivered', () {
        final delivered = order.copyWith(status: OrderStatus.delivered);
        expect(delivered.isCancellable, false);
      });

      test('false when cancelled', () {
        final cancelled = order.copyWith(status: OrderStatus.cancelled);
        expect(cancelled.isCancellable, false);
      });
    });

    group('isReturnable', () {
      test('true only when delivered', () {
        final delivered = order.copyWith(status: OrderStatus.delivered);
        expect(delivered.isReturnable, true);
      });

      test('false for all other statuses', () {
        for (final status in OrderStatus.values) {
          if (status == OrderStatus.delivered) continue;
          final o = order.copyWith(status: status);
          expect(o.isReturnable, false, reason: '$status should not be returnable');
        }
      });
    });

    group('serialization', () {
      test('round-trips through fromJson/toJson', () {
        final json = {
          'id': 'order-3',
          'orderNumber': 'ORD-2026-0003',
          'status': 'confirmed',
          'items': [
            {
              'id': 'item-1',
              'productId': 'prod-1',
              'name': 'Thing',
              'quantity': 3,
              'priceCents': 999,
              'imageUrl': 'https://example.com/thing.jpg',
            },
          ],
          'subtotalCents': 2997,
          'taxCents': 300,
          'shippingCents': 500,
          'totalCents': 3797,
          'shippingAddress': {
            'street': '123 Main St',
            'city': 'Springfield',
            'state': 'IL',
            'zip': '62701',
            'country': 'US',
            'firstName': 'John',
            'lastName': 'Doe',
            'phone': '+15551234567',
          },
          'createdAt': '2026-03-01T00:00:00.000',
        };

        final fromJson = Order.fromJson(json);
        expect(fromJson.id, 'order-3');
        expect(fromJson.status, OrderStatus.confirmed);
        expect(fromJson.items.length, 1);
        expect(fromJson.items.first.name, 'Thing');
        expect(fromJson.shippingAddress, isNotNull);
        expect(fromJson.shippingAddress!.city, 'Springfield');

        // Round-trip
        final backToJson = fromJson.toJson();
        final roundTrip = Order.fromJson(backToJson);
        expect(roundTrip.id, fromJson.id);
        expect(roundTrip.status, fromJson.status);
        expect(roundTrip.totalCents, fromJson.totalCents);
      });
    });

    group('equality', () {
      test('equal when same id', () {
        final a = order;
        final b = order.copyWith(orderNumber: 'DIFFERENT');
        expect(a, equals(b));
      });

      test('not equal when different id', () {
        final other = order.copyWith(id: 'order-other');
        expect(order, isNot(equals(other)));
      });
    });

    test('toString includes key fields', () {
      final s = order.toString();
      expect(s, contains('order-1'));
      expect(s, contains('ORD-2026-0001'));
      expect(s, contains('pending'));
    });
  });

  group('OrderItem', () {
    test('totalCents multiplies price by quantity', () {
      const item = OrderItem(
        id: 'item-1',
        productId: 'prod-1',
        name: 'Widget',
        quantity: 3,
        priceCents: 1500,
      );
      expect(item.totalCents, 4500);
    });

    test('serialization round-trip', () {
      final json = {
        'id': 'item-2',
        'productId': 'prod-2',
        'name': 'Gadget',
        'quantity': 1,
        'priceCents': 2000,
        'imageUrl': 'https://example.com/gadget.jpg',
      };
      final item = OrderItem.fromJson(json);
      expect(item.imageUrl, 'https://example.com/gadget.jpg');

      final roundTrip = OrderItem.fromJson(item.toJson());
      expect(roundTrip.id, item.id);
      expect(roundTrip.priceCents, item.priceCents);
    });

    test('equality by id', () {
      const a = OrderItem(
        id: 'item-1',
        productId: 'prod-1',
        name: 'A',
        quantity: 1,
        priceCents: 100,
      );
      const b = OrderItem(
        id: 'item-1',
        productId: 'prod-2',
        name: 'B',
        quantity: 5,
        priceCents: 999,
      );
      expect(a, equals(b));
    });
  });

  group('Address', () {
    const address = Address(
      street: '456 Oak Ave',
      city: 'Portland',
      state: 'OR',
      zip: '97201',
      country: 'US',
      firstName: 'Jane',
      lastName: 'Smith',
      phone: '+15559876543',
    );

    test('fullName concatenates first and last', () {
      expect(address.fullName, 'Jane Smith');
    });

    test('formatted returns single-line address', () {
      expect(address.formatted, '456 Oak Ave, Portland, OR 97201, US');
    });

    test('serialization round-trip', () {
      final json = address.toJson();
      final roundTrip = Address.fromJson(json);
      expect(roundTrip.street, address.street);
      expect(roundTrip.city, address.city);
      expect(roundTrip.phone, address.phone);
    });

    test('copyWith overrides selected fields', () {
      final updated = address.copyWith(city: 'Seattle', state: 'WA');
      expect(updated.city, 'Seattle');
      expect(updated.state, 'WA');
      expect(updated.street, address.street); // unchanged
    });

    test('toString uses formatted', () {
      expect(address.toString(), 'Address(456 Oak Ave, Portland, OR 97201, US)');
    });
  });

  group('OrderStatus', () {
    test('all enum values have correct string values', () {
      expect(OrderStatus.pending.value, 'pending');
      expect(OrderStatus.confirmed.value, 'confirmed');
      expect(OrderStatus.processing.value, 'processing');
      expect(OrderStatus.shipped.value, 'shipped');
      expect(OrderStatus.outForDelivery.value, 'out_for_delivery');
      expect(OrderStatus.delivered.value, 'delivered');
      expect(OrderStatus.cancelled.value, 'cancelled');
      expect(OrderStatus.returned.value, 'returned');
      expect(OrderStatus.refunded.value, 'refunded');
    });

    test('displayLabel returns human-readable labels', () {
      expect(OrderStatus.pending.displayLabel, 'Pending');
      expect(OrderStatus.outForDelivery.displayLabel, 'Out for Delivery');
      expect(OrderStatus.delivered.displayLabel, 'Delivered');
      expect(OrderStatus.cancelled.displayLabel, 'Cancelled');
    });

    test('all 9 statuses defined', () {
      expect(OrderStatus.values.length, 9);
    });
  });
}
