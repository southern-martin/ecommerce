import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import { Loader2 } from 'lucide-react';
import type { SellerCarrier } from '../services/seller-shipping.api';
import type { CreateShipmentInput, ShippingAddress, ShipmentItemInput } from '../services/seller-shipping.api';

const shipmentSchema = z.object({
  carrier_code: z.string().min(1, 'Select a carrier'),
  origin_street: z.string().min(1, 'Street is required'),
  origin_city: z.string().min(1, 'City is required'),
  origin_state: z.string().min(1, 'State is required'),
  origin_postal_code: z.string().min(1, 'Postal code is required'),
  origin_country: z.string().min(1, 'Country is required'),
  weight_grams: z.coerce.number().min(1, 'Weight must be at least 1g'),
});

type ShipmentFormValues = z.infer<typeof shipmentSchema>;

interface CreateShipmentFormProps {
  orderId: string;
  destinationAddress: ShippingAddress;
  items: ShipmentItemInput[];
  carriers: SellerCarrier[];
  onSubmit: (data: CreateShipmentInput) => void;
  isPending: boolean;
}

export function CreateShipmentForm({
  orderId,
  destinationAddress,
  items,
  carriers,
  onSubmit,
  isPending,
}: CreateShipmentFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<ShipmentFormValues>({
    resolver: zodResolver(shipmentSchema),
    defaultValues: {
      carrier_code: '',
      origin_street: '',
      origin_city: '',
      origin_state: '',
      origin_postal_code: '',
      origin_country: 'US',
      weight_grams: 500,
    },
  });

  const carrierValue = watch('carrier_code');

  const handleFormSubmit = (data: ShipmentFormValues) => {
    onSubmit({
      order_id: orderId,
      carrier_code: data.carrier_code,
      origin: {
        street: data.origin_street,
        city: data.origin_city,
        state: data.origin_state,
        postal_code: data.origin_postal_code,
        country: data.origin_country,
      },
      destination: destinationAddress,
      weight_grams: data.weight_grams,
      currency: 'USD',
      items,
    });
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
      <div className="space-y-2">
        <Label>Carrier</Label>
        {carriers.length > 0 ? (
          <Select
            value={carrierValue}
            onValueChange={(val) => setValue('carrier_code', val)}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select carrier" />
            </SelectTrigger>
            <SelectContent>
              {carriers.map((c) => (
                <SelectItem key={c.id} value={c.carrier_code}>
                  {c.carrier_name || c.carrier_code}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        ) : (
          <Input
            {...register('carrier_code')}
            placeholder="e.g. ups, fedex, usps"
          />
        )}
        {errors.carrier_code && (
          <p className="text-sm text-destructive">{errors.carrier_code.message}</p>
        )}
      </div>

      <div>
        <h4 className="mb-3 text-sm font-medium">Origin Address (Your Warehouse)</h4>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2 sm:col-span-2">
            <Label htmlFor="origin_street">Street</Label>
            <Input id="origin_street" {...register('origin_street')} placeholder="123 Warehouse St" />
            {errors.origin_street && (
              <p className="text-sm text-destructive">{errors.origin_street.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="origin_city">City</Label>
            <Input id="origin_city" {...register('origin_city')} placeholder="City" />
            {errors.origin_city && (
              <p className="text-sm text-destructive">{errors.origin_city.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="origin_state">State</Label>
            <Input id="origin_state" {...register('origin_state')} placeholder="State" />
            {errors.origin_state && (
              <p className="text-sm text-destructive">{errors.origin_state.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="origin_postal_code">Postal Code</Label>
            <Input id="origin_postal_code" {...register('origin_postal_code')} placeholder="12345" />
            {errors.origin_postal_code && (
              <p className="text-sm text-destructive">{errors.origin_postal_code.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="origin_country">Country</Label>
            <Input id="origin_country" {...register('origin_country')} placeholder="US" />
            {errors.origin_country && (
              <p className="text-sm text-destructive">{errors.origin_country.message}</p>
            )}
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="weight_grams">Package Weight (grams)</Label>
        <Input id="weight_grams" type="number" {...register('weight_grams')} />
        {errors.weight_grams && (
          <p className="text-sm text-destructive">{errors.weight_grams.message}</p>
        )}
      </div>

      <div className="rounded-lg border bg-muted/50 p-4">
        <h4 className="mb-2 text-sm font-medium">Destination</h4>
        <p className="text-sm text-muted-foreground">
          {destinationAddress.street}, {destinationAddress.city}, {destinationAddress.state}{' '}
          {destinationAddress.postal_code}, {destinationAddress.country}
        </p>
        <h4 className="mb-2 mt-3 text-sm font-medium">Items ({items.length})</h4>
        <ul className="text-sm text-muted-foreground">
          {items.map((item, i) => (
            <li key={i}>
              {item.product_name} x{item.quantity}
            </li>
          ))}
        </ul>
      </div>

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Create Shipment
      </Button>
    </form>
  );
}
