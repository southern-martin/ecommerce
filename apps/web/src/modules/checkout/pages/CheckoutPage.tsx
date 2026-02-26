import { Button } from '@/shared/components/ui/button';
import { Loader2 } from 'lucide-react';
import { CheckoutStepper } from '../components/CheckoutStepper';
import { ShippingForm } from '../components/ShippingForm';
import { PaymentForm } from '../components/PaymentForm';
import { useCheckout } from '../hooks/useCheckout';
import { useCart } from '@/modules/cart/hooks/useCart';

export default function CheckoutPage() {
  const {
    currentStep,
    currentStepIndex,
    steps,
    shippingAddress,
    setShippingAddress,
    setPaymentMethodId,
    couponCode,
    setCouponCode,
    goToNextStep,
    goToPreviousStep,
    submitOrder,
    isSubmitting,
  } = useCheckout();

  const { cart } = useCart();

  return (
    <div className="mx-auto max-w-3xl">
      <h1 className="mb-6 text-3xl font-bold">Checkout</h1>

      <CheckoutStepper steps={steps} currentStepIndex={currentStepIndex} />

      {currentStep === 'address' && (
        <div className="space-y-6">
          <h2 className="text-xl font-semibold">Shipping Address</h2>
          <ShippingForm
            defaultValues={shippingAddress}
            onSubmit={(address) => {
              setShippingAddress(address);
              goToNextStep();
            }}
          />
        </div>
      )}

      {currentStep === 'payment' && (
        <div className="space-y-6">
          <h2 className="text-xl font-semibold">Payment Method</h2>
          <PaymentForm
            onBack={goToPreviousStep}
            onContinue={(method) => {
              setPaymentMethodId(method);
              goToNextStep();
            }}
            couponCode={couponCode}
            onCouponChange={setCouponCode}
          />
        </div>
      )}

      {currentStep === 'review' && (
        <div className="space-y-6">
          <h2 className="text-xl font-semibold">Review Your Order</h2>

          {shippingAddress && (
            <div className="rounded-lg border p-4">
              <h3 className="text-sm font-medium">Shipping Address</h3>
              <p className="mt-1 text-sm text-muted-foreground">
                {shippingAddress.first_name} {shippingAddress.last_name}
                <br />
                {shippingAddress.address_line1}
                {shippingAddress.address_line2 && <>, {shippingAddress.address_line2}</>}
                <br />
                {shippingAddress.city}, {shippingAddress.state} {shippingAddress.postal_code}
                <br />
                {shippingAddress.country}
              </p>
            </div>
          )}

          {cart && (
            <div className="rounded-lg border p-4">
              <h3 className="text-sm font-medium">Items ({cart.item_count})</h3>
              <div className="mt-2 space-y-2">
                {cart.items.map((item) => (
                  <div key={item.id} className="flex items-center gap-3 text-sm">
                    <img src={item.image_url} alt={item.name} className="h-10 w-10 rounded object-cover" />
                    <span className="flex-1">{item.name} x{item.quantity}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="flex gap-4">
            <Button variant="outline" onClick={goToPreviousStep}>
              Back
            </Button>
            <Button className="flex-1" onClick={submitOrder} disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Place Order
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
