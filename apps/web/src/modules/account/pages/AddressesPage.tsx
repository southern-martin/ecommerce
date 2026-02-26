import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/shared/lib/api-client';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Card } from '@/shared/components/ui/card';

interface Address {
  id: string;
  label: string;
  first_name: string;
  last_name: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  phone: string;
  is_default: boolean;
}

export default function AddressesPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);

  const { data: addresses = [], isLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: async () => {
      const res = await apiClient.get<{ data: Address[] } | Address[]>('/users/me/addresses');
      const d = res.data;
      return Array.isArray(d) ? d : (d as any).data || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: (address: Partial<Address>) => apiClient.post('/users/me/addresses', address),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] });
      setShowForm(false);
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    createMutation.mutate({
      label: form.get('label') as string,
      first_name: form.get('first_name') as string,
      last_name: form.get('last_name') as string,
      address_line1: form.get('address_line1') as string,
      city: form.get('city') as string,
      state: form.get('state') as string,
      postal_code: form.get('postal_code') as string,
      country: form.get('country') as string || 'US',
      phone: form.get('phone') as string,
    });
  };

  if (isLoading) return <div className="p-6">Loading addresses...</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">My Addresses</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancel' : 'Add Address'}</Button>
      </div>

      {showForm && (
        <Card className="p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div><Label>Label</Label><Input name="label" placeholder="Home" required /></div>
              <div><Label>Phone</Label><Input name="phone" placeholder="+1..." /></div>
              <div><Label>First Name</Label><Input name="first_name" required /></div>
              <div><Label>Last Name</Label><Input name="last_name" required /></div>
            </div>
            <div><Label>Address Line 1</Label><Input name="address_line1" required /></div>
            <div className="grid grid-cols-3 gap-4">
              <div><Label>City</Label><Input name="city" required /></div>
              <div><Label>State</Label><Input name="state" required /></div>
              <div><Label>Postal Code</Label><Input name="postal_code" required /></div>
            </div>
            <div><Label>Country</Label><Input name="country" defaultValue="US" /></div>
            <Button type="submit" disabled={createMutation.isPending}>Save Address</Button>
          </form>
        </Card>
      )}

      <div className="grid gap-4 md:grid-cols-2">
        {addresses.map((addr: Address) => (
          <Card key={addr.id} className="p-4">
            <div className="flex justify-between items-start">
              <div>
                <div className="font-medium">{addr.label} {addr.is_default && <span className="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded ml-2">Default</span>}</div>
                <div className="text-sm text-muted-foreground mt-1">
                  {addr.first_name} {addr.last_name}<br/>
                  {addr.address_line1}<br/>
                  {addr.city}, {addr.state} {addr.postal_code}<br/>
                  {addr.country}
                </div>
              </div>
            </div>
          </Card>
        ))}
        {addresses.length === 0 && !showForm && (
          <p className="text-muted-foreground col-span-2 text-center py-8">No addresses yet. Add your first address.</p>
        )}
      </div>
    </div>
  );
}
