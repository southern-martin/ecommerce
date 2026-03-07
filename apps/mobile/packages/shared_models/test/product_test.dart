import 'package:test/test.dart';
import 'package:ecommerce_shared_models/ecommerce_shared_models.dart';

void main() {
  group('Product', () {
    late Product product;

    setUp(() {
      product = const Product(
        id: 'prod-1',
        name: 'Test Widget',
        slug: 'test-widget',
        description: 'A high-quality widget.',
        priceCents: 2999,
      );
    });

    test('constructs with required fields and correct defaults', () {
      expect(product.id, 'prod-1');
      expect(product.name, 'Test Widget');
      expect(product.slug, 'test-widget');
      expect(product.description, 'A high-quality widget.');
      expect(product.priceCents, 2999);
      expect(product.compareAtPriceCents, isNull);
      expect(product.images, isEmpty);
      expect(product.category, isNull);
      expect(product.rating, 0.0);
      expect(product.reviewCount, 0);
      expect(product.inStock, true);
      expect(product.stockQuantity, 0);
      expect(product.seller, isNull);
      expect(product.variants, isEmpty);
    });

    group('primaryImageUrl', () {
      test('returns null when no images', () {
        expect(product.primaryImageUrl, isNull);
      });

      test('returns first image url when no primary flagged', () {
        final p = product.copyWith(
          images: [
            const ProductImage(id: 'img-1', url: 'https://example.com/a.jpg'),
            const ProductImage(id: 'img-2', url: 'https://example.com/b.jpg'),
          ],
        );
        expect(p.primaryImageUrl, 'https://example.com/a.jpg');
      });

      test('returns primary image url when flagged', () {
        final p = product.copyWith(
          images: [
            const ProductImage(id: 'img-1', url: 'https://example.com/a.jpg'),
            const ProductImage(
              id: 'img-2',
              url: 'https://example.com/b.jpg',
              isPrimary: true,
            ),
          ],
        );
        expect(p.primaryImageUrl, 'https://example.com/b.jpg');
      });
    });

    group('isOnSale', () {
      test('returns false when no compareAtPriceCents', () {
        expect(product.isOnSale, false);
      });

      test('returns false when compareAtPriceCents <= priceCents', () {
        final p = product.copyWith(compareAtPriceCents: 2999);
        expect(p.isOnSale, false);
      });

      test('returns true when compareAtPriceCents > priceCents', () {
        final p = product.copyWith(compareAtPriceCents: 3999);
        expect(p.isOnSale, true);
      });
    });

    group('discountPercentage', () {
      test('returns 0 when not on sale', () {
        expect(product.discountPercentage, 0);
      });

      test('calculates correct percentage', () {
        final p = product.copyWith(
          priceCents: 7500,
          compareAtPriceCents: 10000,
        );
        expect(p.discountPercentage, 25);
      });

      test('rounds to nearest integer', () {
        final p = product.copyWith(
          priceCents: 6666,
          compareAtPriceCents: 10000,
        );
        // (10000 - 6666) / 10000 * 100 = 33.34 → rounds to 33
        expect(p.discountPercentage, 33);
      });
    });

    group('hasVariants', () {
      test('returns false when no variants', () {
        expect(product.hasVariants, false);
      });

      test('returns true when variants present', () {
        final p = product.copyWith(
          variants: [
            const ProductVariant(
              id: 'v-1',
              name: 'Size',
              value: 'Large',
            ),
          ],
        );
        expect(p.hasVariants, true);
      });
    });

    group('serialization', () {
      test('round-trips through fromJson/toJson', () {
        final json = {
          'id': 'prod-2',
          'name': 'Gizmo',
          'slug': 'gizmo',
          'description': 'A fine gizmo.',
          'priceCents': 4500,
          'compareAtPriceCents': 5000,
          'images': [
            {
              'id': 'img-1',
              'url': 'https://example.com/gizmo.jpg',
              'alt': 'Gizmo photo',
              'isPrimary': true,
            },
          ],
          'category': 'electronics',
          'rating': 4.5,
          'reviewCount': 12,
          'inStock': true,
          'stockQuantity': 50,
          'seller': 'seller-1',
          'variants': [
            {
              'id': 'v-1',
              'name': 'Color',
              'value': 'Red',
              'priceModifier': 100,
              'stockQuantity': 10,
            },
          ],
        };

        final fromJson = Product.fromJson(json);
        expect(fromJson.id, 'prod-2');
        expect(fromJson.name, 'Gizmo');
        expect(fromJson.priceCents, 4500);
        expect(fromJson.compareAtPriceCents, 5000);
        expect(fromJson.images.length, 1);
        expect(fromJson.images.first.isPrimary, true);
        expect(fromJson.category, 'electronics');
        expect(fromJson.rating, 4.5);
        expect(fromJson.reviewCount, 12);
        expect(fromJson.variants.length, 1);
        expect(fromJson.variants.first.name, 'Color');

        // Round-trip
        final backToJson = fromJson.toJson();
        final roundTrip = Product.fromJson(backToJson);
        expect(roundTrip.id, fromJson.id);
        expect(roundTrip.name, fromJson.name);
        expect(roundTrip.priceCents, fromJson.priceCents);
        expect(roundTrip.images.length, fromJson.images.length);
        expect(roundTrip.variants.length, fromJson.variants.length);
      });

      test('handles minimal JSON with defaults', () {
        final json = {
          'id': 'prod-3',
          'name': 'Basic',
          'slug': 'basic',
          'description': 'Simple.',
          'priceCents': 100,
        };

        final p = Product.fromJson(json);
        expect(p.images, isEmpty);
        expect(p.variants, isEmpty);
        expect(p.rating, 0.0);
        expect(p.inStock, true);
      });
    });

    group('equality', () {
      test('equal when same id', () {
        final a = product;
        final b = product.copyWith(name: 'Different Name');
        expect(a, equals(b));
      });

      test('not equal when different id', () {
        final other = Product(
          id: 'prod-other',
          name: product.name,
          slug: product.slug,
          description: product.description,
          priceCents: product.priceCents,
        );
        expect(product, isNot(equals(other)));
      });

      test('consistent hashCode based on id', () {
        final a = product;
        final b = product.copyWith(name: 'Another');
        expect(a.hashCode, b.hashCode);
      });
    });

    test('toString includes id, name, and price', () {
      expect(
        product.toString(),
        'Product(id: prod-1, name: Test Widget, price: 2999)',
      );
    });
  });

  group('ProductImage', () {
    test('constructs with defaults', () {
      const img = ProductImage(id: 'img-1', url: 'https://example.com/a.jpg');
      expect(img.alt, isNull);
      expect(img.isPrimary, false);
    });

    test('serialization round-trip', () {
      final json = {
        'id': 'img-1',
        'url': 'https://example.com/photo.jpg',
        'alt': 'Photo',
        'isPrimary': true,
      };
      final img = ProductImage.fromJson(json);
      expect(img.id, 'img-1');
      expect(img.isPrimary, true);

      final backToJson = img.toJson();
      final roundTrip = ProductImage.fromJson(backToJson);
      expect(roundTrip.id, img.id);
      expect(roundTrip.url, img.url);
    });

    test('equality by id', () {
      const a = ProductImage(id: 'img-1', url: 'https://a.com');
      const b = ProductImage(id: 'img-1', url: 'https://b.com');
      expect(a, equals(b));
    });
  });

  group('ProductVariant', () {
    test('constructs with defaults', () {
      const v = ProductVariant(id: 'v-1', name: 'Size', value: 'M');
      expect(v.priceModifier, 0);
      expect(v.stockQuantity, 0);
    });

    test('serialization round-trip', () {
      final json = {
        'id': 'v-1',
        'name': 'Color',
        'value': 'Blue',
        'priceModifier': 200,
        'stockQuantity': 5,
      };
      final v = ProductVariant.fromJson(json);
      expect(v.priceModifier, 200);
      expect(v.stockQuantity, 5);

      final roundTrip = ProductVariant.fromJson(v.toJson());
      expect(roundTrip.id, v.id);
      expect(roundTrip.value, v.value);
    });

    test('equality by id', () {
      const a = ProductVariant(id: 'v-1', name: 'Size', value: 'S');
      const b = ProductVariant(id: 'v-1', name: 'Size', value: 'L');
      expect(a, equals(b));
    });
  });
}
