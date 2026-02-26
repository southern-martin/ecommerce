import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import apiClient from '@/shared/lib/api-client';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Card } from '@/shared/components/ui/card';

export default function ReturnRequestPage() {
  const navigate = useNavigate();
  const [orderId, setOrderId] = useState('');
  const [reason, setReason] = useState('');

  const mutation = useMutation({
    mutationFn: (data: { order_id: string; reason: string }) =>
      apiClient.post('/returns', data),
    onSuccess: () => navigate('/account/returns'),
  });

  return (
    <div className="max-w-xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Request a Return</h1>
      <Card className="p-6">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            mutation.mutate({ order_id: orderId, reason });
          }}
          className="space-y-4"
        >
          <div className="space-y-2">
            <Label>Order ID</Label>
            <Input value={orderId} onChange={(e) => setOrderId(e.target.value)} placeholder="Enter your order ID" required />
          </div>
          <div className="space-y-2">
            <Label>Reason for Return</Label>
            <Input value={reason} onChange={(e) => setReason(e.target.value)} placeholder="Why are you returning?" required />
          </div>
          <Button type="submit" disabled={mutation.isPending}>Submit Return Request</Button>
        </form>
      </Card>
    </div>
  );
}
