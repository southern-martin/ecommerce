import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    post: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { orderApi } from '../order.api';
import { paymentApi } from '../payment.api';

const mockApiClient = apiClient as unknown as {
  post: ReturnType<typeof vi.fn>;
};

describe('orderApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('createOrder', () => {
    it('sends POST to /orders with mapped payload', async () => {
      const mockResponse = {
        data: {
          data: {
            id: 'ord-1',
            order_number: 'ORD-100',
            status: 'pending',
            items: [
              { id: 'i1', product_id: 'p1', product_name: 'Widget', unit_price_cents: 1000, quantity: 2 },
            ],
            subtotal_cents: 2000,
            shipping_cents: 500,
            tax_cents: 180,
            discount_cents: 0,
            total_cents: 2680,
            created_at: '2025-06-01T00:00:00Z',
          },
        },
      };
      mockApiClient.post.mockResolvedValue(mockResponse);

      const shippingAddress = {
        first_name: 'John',
        last_name: 'Doe',
        address_line1: '123 Main St',
        city: 'NYC',
        state: 'NY',
        postal_code: '10001',
        country: 'US',
        phone: '555-1234',
      };

      const cartItems = [
        { product_id: 'p1', product_name: 'Widget', quantity: 2, price_cents: 1000, seller_id: 's1' },
      ];

      const result = await orderApi.createOrder(
        { shipping_address: shippingAddress, payment_method_id: 'pm_123' },
        cartItems,
        'user-1',
        'john@example.com'
      );

      expect(mockApiClient.post).toHaveBeenCalledWith('/orders', expect.objectContaining({
        buyer_id: 'user-1',
        buyer_email: 'john@example.com',
        currency: 'USD',
      }));
      expect(result.id).toBe('ord-1');
      expect(result.total).toBe(2680);
      expect(result.items).toHaveLength(1);
    });

    it('maps shipping address to backend format', async () => {
      mockApiClient.post.mockResolvedValue({ data: { data: { id: 'o1', status: 'pending', items: [], created_at: '' } } });

      await orderApi.createOrder(
        {
          shipping_address: {
            first_name: 'Jane',
            last_name: 'Smith',
            address_line1: '456 Oak Ave',
            address_line2: 'Apt 2',
            city: 'LA',
            state: 'CA',
            postal_code: '90001',
            country: 'US',
            phone: '555-9876',
          },
          payment_method_id: 'pm_456',
        },
        [],
        'user-2'
      );

      const payload = mockApiClient.post.mock.calls[0][1];
      expect(payload.shipping_address.full_name).toBe('Jane Smith');
      expect(payload.shipping_address.line1).toBe('456 Oak Ave');
      expect(payload.shipping_address.line2).toBe('Apt 2');
      expect(payload.shipping_address.country_code).toBe('US');
    });

    it('omits variant_id when not provided', async () => {
      mockApiClient.post.mockResolvedValue({ data: { data: { id: 'o1', status: 'pending', items: [], created_at: '' } } });

      await orderApi.createOrder(
        {
          shipping_address: { first_name: 'A', last_name: 'B', address_line1: '1', city: 'C', state: 'S', postal_code: '0', country: 'US', phone: '0' },
          payment_method_id: 'pm_1',
        },
        [{ product_id: 'p1', product_name: 'Test', quantity: 1, price_cents: 100 }],
        'u1'
      );

      const payload = mockApiClient.post.mock.calls[0][1];
      expect(payload.items[0]).not.toHaveProperty('variant_id');
    });
  });
});

describe('paymentApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('createPaymentIntent', () => {
    it('sends POST to /payments/create-intent', async () => {
      const mockIntent = { payment_id: 'pay-1', stripe_payment_id: 'pi_123', client_secret: 'cs_123', status: 'pending' };
      mockApiClient.post.mockResolvedValue({ data: mockIntent });

      const input = {
        order_id: 'ord-1',
        buyer_id: 'u1',
        amount_cents: 5000,
        currency: 'USD',
        seller_items: [{ seller_id: 's1', amount_cents: 5000 }],
      };
      const result = await paymentApi.createPaymentIntent(input);

      expect(mockApiClient.post).toHaveBeenCalledWith('/payments/create-intent', input);
      expect(result.stripe_payment_id).toBe('pi_123');
    });
  });

  describe('simulatePaymentSuccess', () => {
    it('sends POST to /payments/webhooks/stripe', async () => {
      mockApiClient.post.mockResolvedValue({});

      const sellerItems = [{ seller_id: 's1', amount_cents: 3000 }];
      await paymentApi.simulatePaymentSuccess('pi_abc', sellerItems);

      expect(mockApiClient.post).toHaveBeenCalledWith('/payments/webhooks/stripe', {
        type: 'payment_intent.succeeded',
        stripe_payment_id: 'pi_abc',
        seller_items: sellerItems,
      });
    });
  });
});
