import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { orderApi } from '../services/order.api';
import { paymentApi } from '../services/payment.api';
import type { CreateOrderData, ShippingAddress } from '../services/order.api';
import { useCartStore } from '@/shared/stores/cart.store';
import { useAuthStore } from '@/shared/stores/auth.store';

export type CheckoutStep = 'address' | 'payment' | 'review';

const STEPS: CheckoutStep[] = ['address', 'payment', 'review'];

export function useCheckout() {
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState<CheckoutStep>('address');
  const [shippingAddress, setShippingAddress] = useState<ShippingAddress | null>(null);
  const [paymentMethodId, setPaymentMethodId] = useState<string | null>(null);
  const [couponCode, setCouponCode] = useState<string>('');

  const cartItems = useCartStore((s) => s.items);
  const clearCart = useCartStore((s) => s.clearCart);
  const subtotalFn = useCartStore((s) => s.subtotal);
  const user = useAuthStore((s) => s.user);

  const currentStepIndex = STEPS.indexOf(currentStep);

  const goToNextStep = () => {
    if (currentStepIndex < STEPS.length - 1) {
      setCurrentStep(STEPS[currentStepIndex + 1]);
    }
  };

  const goToPreviousStep = () => {
    if (currentStepIndex > 0) {
      setCurrentStep(STEPS[currentStepIndex - 1]);
    }
  };

  const placeOrder = useMutation({
    mutationFn: async (data: CreateOrderData) => {
      const items = cartItems.map((i) => ({
        product_id: i.product_id,
        product_name: i.product_name,
        quantity: i.quantity,
        price_cents: i.price_cents,
        image_url: i.image_url,
        variant_id: i.variant_id,
        seller_id: i.seller_id,
      }));

      // 1. Create the order
      const order = await orderApi.createOrder(
        data,
        items,
        user?.id || '',
        user?.email
      );

      // 2. Trigger payment (demo mode — auto-confirm)
      try {
        // Aggregate seller items for commission splitting
        const sellerMap = new Map<string, number>();
        for (const item of items) {
          const sid = item.seller_id || user?.id || '';
          sellerMap.set(sid, (sellerMap.get(sid) || 0) + item.price_cents * item.quantity);
        }
        const sellerItems = Array.from(sellerMap.entries()).map(([seller_id, amount_cents]) => ({
          seller_id,
          amount_cents,
        }));

        const totalCents = subtotalFn();

        const paymentIntent = await paymentApi.createPaymentIntent({
          order_id: order.id,
          buyer_id: user?.id || '',
          amount_cents: totalCents,
          currency: 'usd',
          seller_items: sellerItems,
        });

        // Simulate Stripe webhook confirmation
        await paymentApi.simulatePaymentSuccess(
          paymentIntent.stripe_payment_id,
          sellerItems
        );
      } catch {
        // Payment failures are non-blocking for demo — order is still placed
        console.warn('Demo payment simulation failed (non-blocking)');
      }

      return order;
    },
    onSuccess: (order) => {
      clearCart();
      navigate(`/order-confirmation/${order.order_number}`, {
        state: { order },
      });
    },
  });

  const submitOrder = () => {
    if (!shippingAddress) return;
    placeOrder.mutate({
      shipping_address: shippingAddress,
      payment_method_id: paymentMethodId || 'cod',
      coupon_code: couponCode || undefined,
    });
  };

  return {
    currentStep,
    currentStepIndex,
    steps: STEPS,
    shippingAddress,
    setShippingAddress,
    paymentMethodId,
    setPaymentMethodId,
    couponCode,
    setCouponCode,
    goToNextStep,
    goToPreviousStep,
    submitOrder,
    isSubmitting: placeOrder.isPending,
    error: placeOrder.error,
  };
}
