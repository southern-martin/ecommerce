import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Tag, X } from 'lucide-react';

interface CouponInputProps {
  value: string;
  onChange: (code: string) => void;
  appliedDiscount?: number;
}

export function CouponInput({ value, onChange, appliedDiscount }: CouponInputProps) {
  const [inputValue, setInputValue] = useState(value);

  const handleApply = () => {
    onChange(inputValue.trim());
  };

  const handleRemove = () => {
    setInputValue('');
    onChange('');
  };

  if (value && appliedDiscount !== undefined) {
    return (
      <div className="flex items-center justify-between rounded-md border bg-muted/50 px-3 py-2">
        <div className="flex items-center gap-2">
          <Tag className="h-4 w-4 text-green-600" />
          <span className="text-sm font-medium">{value}</span>
        </div>
        <Button variant="ghost" size="icon" className="h-6 w-6" onClick={handleRemove}>
          <X className="h-3 w-3" />
        </Button>
      </div>
    );
  }

  return (
    <div className="flex gap-2">
      <Input
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        placeholder="Enter coupon code"
        className="flex-1"
      />
      <Button variant="outline" onClick={handleApply} disabled={!inputValue.trim()}>
        Apply
      </Button>
    </div>
  );
}
