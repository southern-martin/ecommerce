import { Check } from 'lucide-react';
import type { CheckoutStep } from '../hooks/useCheckout';

interface CheckoutStepperProps {
  steps: CheckoutStep[];
  currentStepIndex: number;
}

const stepLabels: Record<CheckoutStep, string> = {
  address: 'Shipping Address',
  payment: 'Payment',
  review: 'Review Order',
};

export function CheckoutStepper({ steps, currentStepIndex }: CheckoutStepperProps) {
  return (
    <nav className="mb-8">
      <ol className="flex items-center gap-2">
        {steps.map((step, index) => (
          <li key={step} className="flex items-center gap-2">
            <div
              className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium ${
                index < currentStepIndex
                  ? 'bg-primary text-primary-foreground'
                  : index === currentStepIndex
                    ? 'border-2 border-primary text-primary'
                    : 'border-2 border-muted text-muted-foreground'
              }`}
            >
              {index < currentStepIndex ? <Check className="h-4 w-4" /> : index + 1}
            </div>
            <span
              className={`text-sm ${
                index <= currentStepIndex ? 'font-medium' : 'text-muted-foreground'
              }`}
            >
              {stepLabels[step]}
            </span>
            {index < steps.length - 1 && (
              <div
                className={`mx-2 h-px w-12 ${
                  index < currentStepIndex ? 'bg-primary' : 'bg-muted'
                }`}
              />
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
}
