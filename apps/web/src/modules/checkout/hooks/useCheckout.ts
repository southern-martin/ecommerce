import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { orderApi } from '../services/order.api';
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
    mutationFn: (data: CreateOrderData) =>
      orderApi.createOrder(
        data,
        cartItems.map((i) => ({
          product_id: i.product_id,
          product_name: i.product_name,
          quantity: i.quantity,
          price_cents: i.price_cents,
          image_url: i.image_url,
          variant_id: i.variant_id,
          seller_id: i.seller_id,
        })),
        user?.id || ''
      ),
    onSuccess: (order) => {
      clearCart();
      navigate(`/order-confirmation/${order.order_number}`);
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
