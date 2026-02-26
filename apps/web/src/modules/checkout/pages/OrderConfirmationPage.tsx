import { useParams, Link } from 'react-router-dom';

export default function OrderConfirmationPage() {
  const { orderId } = useParams<{ orderId: string }>();
  return (
    <div className="mx-auto max-w-2xl text-center py-12">
      <div className="mb-6 text-6xl">🎉</div>
      <h1 className="text-3xl font-bold mb-4">Order Confirmed!</h1>
      <p className="text-muted-foreground mb-2">Your order has been placed successfully.</p>
      <p className="text-sm text-muted-foreground mb-8">Order ID: <span className="font-mono">{orderId}</span></p>
      <div className="flex gap-4 justify-center">
        <Link to="/account/orders" className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
          View Orders
        </Link>
        <Link to="/products" className="inline-flex items-center justify-center rounded-md border px-4 py-2 text-sm font-medium hover:bg-accent">
          Continue Shopping
        </Link>
      </div>
    </div>
  );
}
