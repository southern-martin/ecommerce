import 'package:test/test.dart';
import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';

void main() {
  group('CartItem', () {
    const item = CartItem(
      id: 'ci-1',
      productId: 'prod-1',
      name: 'Widget',
      slug: 'widget',
      priceCents: 2500,
      quantity: 3,
    );

    test('constructs with required fields and correct defaults', () {
      expect(item.id, 'ci-1');
      expect(item.productId, 'prod-1');
      expect(item.name, 'Widget');
      expect(item.slug, 'widget');
      expect(item.priceCents, 2500);
      expect(item.quantity, 3);
      expect(item.imageUrl, isNull);
      expect(item.variantId, isNull);
      expect(item.variantName, isNull);
    });

    test('totalCents multiplies price by quantity', () {
      expect(item.totalCents, 7500); // 2500 * 3
    });

    test('totalCents returns priceCents when quantity is 1', () {
      final single = item.copyWith(quantity: 1);
      expect(single.totalCents, 2500);
    });

    group('serialization', () {
      test('round-trips through fromJson/toJson', () {
        final json = {
          'id': 'ci-2',
          'productId': 'prod-2',
          'name': 'Gadget',
          'slug': 'gadget',
          'imageUrl': 'https://example.com/gadget.jpg',
          'priceCents': 1999,
          'quantity': 2,
          'variantId': 'var-1',
          'variantName': 'Large',
        };

        final fromJson = CartItem.fromJson(json);
        expect(fromJson.id, 'ci-2');
        expect(fromJson.imageUrl, 'https://example.com/gadget.jpg');
        expect(fromJson.variantId, 'var-1');
        expect(fromJson.variantName, 'Large');

        final backToJson = fromJson.toJson();
        final roundTrip = CartItem.fromJson(backToJson);
        expect(roundTrip.id, fromJson.id);
        expect(roundTrip.priceCents, fromJson.priceCents);
        expect(roundTrip.variantName, fromJson.variantName);
      });

      test('handles null optional fields', () {
        final json = {
          'id': 'ci-3',
          'productId': 'prod-3',
          'name': 'Basic',
          'slug': 'basic',
          'priceCents': 500,
          'quantity': 1,
        };

        final fromJson = CartItem.fromJson(json);
        expect(fromJson.imageUrl, isNull);
        expect(fromJson.variantId, isNull);
        expect(fromJson.variantName, isNull);
      });
    });

    test('copyWith overrides selected fields', () {
      final updated = item.copyWith(quantity: 5, priceCents: 3000);
      expect(updated.quantity, 5);
      expect(updated.priceCents, 3000);
      expect(updated.name, 'Widget'); // unchanged
    });

    group('equality', () {
      test('equal when same id', () {
        final a = item;
        final b = item.copyWith(quantity: 10);
        expect(a, equals(b));
      });

      test('not equal when different id', () {
        final other = item.copyWith(id: 'ci-other');
        expect(item, isNot(equals(other)));
      });

      test('hashCode based on id', () {
        final a = item;
        final b = item.copyWith(quantity: 99);
        expect(a.hashCode, b.hashCode);
      });
    });

    test('toString includes key fields', () {
      final s = item.toString();
      expect(s, contains('ci-1'));
      expect(s, contains('Widget'));
      expect(s, contains('3'));
      expect(s, contains('2500'));
    });
  });

  group('Cart', () {
    const emptyCart = Cart(id: 'cart-1');

    final cartWithItems = Cart(
      id: 'cart-2',
      items: const [
        CartItem(
          id: 'ci-1',
          productId: 'prod-1',
          name: 'Widget',
          slug: 'widget',
          priceCents: 2000,
          quantity: 2,
        ),
        CartItem(
          id: 'ci-2',
          productId: 'prod-2',
          name: 'Gadget',
          slug: 'gadget',
          priceCents: 3000,
          quantity: 1,
        ),
      ],
    );

    test('constructs with defaults', () {
      expect(emptyCart.id, 'cart-1');
      expect(emptyCart.items, isEmpty);
    });

    group('subtotalCents', () {
      test('returns 0 for empty cart', () {
        expect(emptyCart.subtotalCents, 0);
      });

      test('sums all item totals', () {
        // (2000 * 2) + (3000 * 1) = 7000
        expect(cartWithItems.subtotalCents, 7000);
      });
    });

    group('itemCount', () {
      test('returns 0 for empty cart', () {
        expect(emptyCart.itemCount, 0);
      });

      test('sums all item quantities', () {
        expect(cartWithItems.itemCount, 3); // 2 + 1
      });
    });

    group('isEmpty / isNotEmpty', () {
      test('empty cart', () {
        expect(emptyCart.isEmpty, true);
        expect(emptyCart.isNotEmpty, false);
      });

      test('cart with items', () {
        expect(cartWithItems.isEmpty, false);
        expect(cartWithItems.isNotEmpty, true);
      });
    });

    group('serialization', () {
      test('round-trips through fromJson/toJson', () {
        final json = {
          'id': 'cart-3',
          'items': [
            {
              'id': 'ci-1',
              'productId': 'prod-1',
              'name': 'Item',
              'slug': 'item',
              'priceCents': 1500,
              'quantity': 4,
            },
          ],
        };

        final cart = Cart.fromJson(json);
        expect(cart.id, 'cart-3');
        expect(cart.items.length, 1);
        expect(cart.items.first.quantity, 4);

        final backToJson = cart.toJson();
        final roundTrip = Cart.fromJson(backToJson);
        expect(roundTrip.id, cart.id);
        expect(roundTrip.items.length, cart.items.length);
      });

      test('handles empty items list', () {
        final json = {'id': 'cart-4', 'items': <Map<String, dynamic>>[]};
        final cart = Cart.fromJson(json);
        expect(cart.isEmpty, true);
      });
    });

    test('copyWith overrides selected fields', () {
      final updated = emptyCart.copyWith(
        items: [
          const CartItem(
            id: 'ci-new',
            productId: 'p-new',
            name: 'New',
            slug: 'new',
            priceCents: 100,
            quantity: 1,
          ),
        ],
      );
      expect(updated.id, 'cart-1'); // unchanged
      expect(updated.items.length, 1);
    });

    test('toString includes id, count, and subtotal', () {
      final s = cartWithItems.toString();
      expect(s, contains('cart-2'));
      expect(s, contains('2')); // item count
      expect(s, contains('7000')); // subtotal
    });
  });
}
