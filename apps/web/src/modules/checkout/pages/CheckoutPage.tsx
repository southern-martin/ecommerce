import { Button } from '@/shared/components/ui/button';
import { Loader2 } from 'lucide-react';
import { CheckoutStepper } from '../components/CheckoutStepper';
import { CheckoutSummary } from '../components/CheckoutSummary';
import { ShippingForm } from '../components/ShippingForm';
import { PaymentForm } from '../components/PaymentForm';
import { useCheckout } from '../hooks/useCheckout';
import { useCart } from '@/modules/cart/hooks/useCart';
import { PageLayout } from '@/shared/components/layout/PageLayout';

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
    <PageLayout
      title="Checkout"
      breadcrumbs={[{ label: 'Cart', href: '/cart' }, { label: 'Checkout' }]}
    >
      <div className="mx-auto max-w-3xl">
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
              <div className="rounded-2xl border bg-card p-5">
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
              <CheckoutSummary
                items={cart.items}
                subtotal={cart.subtotal}
                shipping={0}
                tax={0}
                discount={0}
                total={cart.subtotal}
              />
            )}

            <div className="flex gap-4">
              <Button variant="outline" className="rounded-xl" onClick={goToPreviousStep}>
                Back
              </Button>
              <Button
                className="flex-1 rounded-xl font-semibold"
                onClick={submitOrder}
                disabled={isSubmitting}
              >
                {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Place Order
              </Button>
            </div>
          </div>
        )}
      </div>
    </PageLayout>
  );
}
