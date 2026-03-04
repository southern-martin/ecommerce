import apiClient from '@/shared/lib/api-client';

export interface PaymentIntent {
  payment_id: string;
  stripe_payment_id: string;
  client_secret: string;
  status: string;
}

export const paymentApi = {
  /**
   * Create a payment intent for an order.
   * Backend route: POST /api/v1/payments/create-intent
   */
  createPaymentIntent: async (input: {
    order_id: string;
    buyer_id: string;
    amount_cents: number;
    currency: string;
    seller_items: Array<{ seller_id: string; amount_cents: number }>;
  }): Promise<PaymentIntent> => {
    const response = await apiClient.post('/payments/create-intent', input);
    return response.data;
  },

  /**
   * Simulate a Stripe webhook to confirm payment (demo mode).
   * Backend route: POST /api/v1/payments/webhooks/stripe
   */
  simulatePaymentSuccess: async (
    stripePaymentId: string,
    sellerItems: Array<{ seller_id: string; amount_cents: number }>
  ): Promise<void> => {
    await apiClient.post('/payments/webhooks/stripe', {
      type: 'payment_intent.succeeded',
      stripe_payment_id: stripePaymentId,
      seller_items: sellerItems,
    });
  },
};
