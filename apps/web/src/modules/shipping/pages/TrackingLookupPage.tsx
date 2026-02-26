import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Badge } from '@/shared/components/ui/badge';
import { Search, Package } from 'lucide-react';
import { TrackingTimeline } from '../components/TrackingTimeline';
import { useTracking } from '../hooks/useTracking';

export default function TrackingLookupPage() {
  const [trackingNumber, setTrackingNumber] = useState('');
  const [searchedNumber, setSearchedNumber] = useState('');
  const { data: tracking, isLoading } = useTracking(searchedNumber);

  const handleSearch = () => {
    if (trackingNumber.trim()) {
      setSearchedNumber(trackingNumber.trim());
    }
  };

  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="mb-6 text-2xl font-bold">Track Your Order</h1>

      <div className="mb-8 flex gap-2">
        <div className="flex-1 space-y-2">
          <Label htmlFor="tracking">Tracking Number</Label>
          <Input
            id="tracking"
            value={trackingNumber}
            onChange={(e) => setTrackingNumber(e.target.value)}
            placeholder="Enter tracking number"
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
        </div>
        <div className="flex items-end">
          <Button onClick={handleSearch} disabled={!trackingNumber.trim()}>
            <Search className="mr-2 h-4 w-4" />
            Track
          </Button>
        </div>
      </div>

      {isLoading && (
        <div className="flex items-center justify-center py-16">
          <p className="text-muted-foreground">Looking up tracking information...</p>
        </div>
      )}

      {tracking && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2 text-base">
                <Package className="h-5 w-5" />
                {tracking.carrier}
              </CardTitle>
              <Badge>{tracking.status}</Badge>
            </div>
            <p className="text-sm text-muted-foreground">
              Tracking: {tracking.tracking_number}
            </p>
            {tracking.estimated_delivery && (
              <p className="text-sm text-muted-foreground">
                Estimated delivery: {tracking.estimated_delivery}
              </p>
            )}
          </CardHeader>
          <CardContent>
            <TrackingTimeline events={tracking.events} />
          </CardContent>
        </Card>
      )}

      {searchedNumber && !tracking && !isLoading && (
        <div className="flex flex-col items-center py-16">
          <Package className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-muted-foreground">No tracking information found.</p>
        </div>
      )}
    </div>
  );
}
