import { Check, MapPin } from 'lucide-react';
import { formatDateTime } from '@/shared/lib/utils';
import type { TrackingEvent } from '../services/shipping.api';

interface TrackingTimelineProps {
  events: TrackingEvent[];
}

export function TrackingTimeline({ events }: TrackingTimelineProps) {
  return (
    <div className="space-y-0">
      {events.map((event, index) => (
        <div key={index} className="flex gap-4">
          <div className="flex flex-col items-center">
            <div
              className={`flex h-8 w-8 items-center justify-center rounded-full ${
                index === 0 ? 'bg-primary text-primary-foreground' : 'bg-muted'
              }`}
            >
              {index === 0 ? (
                <Check className="h-4 w-4" />
              ) : (
                <MapPin className="h-4 w-4 text-muted-foreground" />
              )}
            </div>
            {index < events.length - 1 && (
              <div className="h-12 w-px bg-border" />
            )}
          </div>
          <div className="pb-8">
            <p className="text-sm font-medium">{event.description}</p>
            <p className="text-xs text-muted-foreground">{event.location}</p>
            <p className="text-xs text-muted-foreground">
              {formatDateTime(event.timestamp)}
            </p>
          </div>
        </div>
      ))}
    </div>
  );
}
