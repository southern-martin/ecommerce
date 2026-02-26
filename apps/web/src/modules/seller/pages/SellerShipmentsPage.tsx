import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { ChevronLeft, ChevronRight, Loader2 } from 'lucide-react';
import { ShipmentTable } from '../components/ShipmentTable';
import {
  useSellerShipments,
  useSellerCarriers,
  useSetupCarrier,
} from '../hooks/useSellerShipping';
import type { SellerCarrier } from '../services/seller-shipping.api';

function CarrierSetupForm() {
  const setupCarrier = useSetupCarrier();
  const { register, handleSubmit, reset } = useForm<{
    carrier_code: string;
    account_number: string;
  }>();

  const onSubmit = (data: { carrier_code: string; account_number: string }) => {
    setupCarrier.mutate(data, { onSuccess: () => reset() });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="carrier_code">Carrier Code</Label>
          <Input
            id="carrier_code"
            {...register('carrier_code')}
            placeholder="e.g. ups, fedex, usps"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="account_number">Account Number</Label>
          <Input
            id="account_number"
            {...register('account_number')}
            placeholder="Your carrier account number"
          />
        </div>
      </div>
      <Button type="submit" disabled={setupCarrier.isPending}>
        {setupCarrier.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Add Carrier
      </Button>
    </form>
  );
}

function CarriersList() {
  const { data: carriers, isLoading } = useSellerCarriers();

  if (isLoading) {
    return <Skeleton className="h-32 w-full" />;
  }

  const carrierList: SellerCarrier[] = carriers ?? [];

  return (
    <div className="space-y-4">
      {carrierList.length > 0 ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {carrierList.map((carrier) => (
            <Card key={carrier.id}>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">{carrier.carrier_name}</p>
                    <p className="text-sm text-muted-foreground">{carrier.carrier_code}</p>
                  </div>
                  <StatusBadge status={carrier.is_active ? 'active' : 'inactive'} />
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <p className="py-4 text-center text-muted-foreground">No carriers configured.</p>
      )}
    </div>
  );
}

export default function SellerShipmentsPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useSellerShipments(page);

  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;
  const shipments = data?.data ?? [];

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Shipping</h1>

      <Tabs defaultValue="shipments">
        <TabsList>
          <TabsTrigger value="shipments">Shipments</TabsTrigger>
          <TabsTrigger value="carriers">Carriers</TabsTrigger>
        </TabsList>

        <TabsContent value="shipments">
          {isLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : shipments.length > 0 ? (
            <>
              <ShipmentTable shipments={shipments} />
              {totalPages > 1 && (
                <div className="mt-6 flex items-center justify-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={page === 1}
                    onClick={() => setPage((p) => p - 1)}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {page} of {totalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={page === totalPages}
                    onClick={() => setPage((p) => p + 1)}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </>
          ) : (
            <p className="py-8 text-center text-muted-foreground">No shipments yet.</p>
          )}
        </TabsContent>

        <TabsContent value="carriers" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Add Carrier</CardTitle>
            </CardHeader>
            <CardContent>
              <CarrierSetupForm />
            </CardContent>
          </Card>

          <div>
            <h2 className="mb-4 text-lg font-semibold">Your Carriers</h2>
            <CarriersList />
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
