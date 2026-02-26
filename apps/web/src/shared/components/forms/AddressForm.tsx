import * as React from "react";
import { Input } from "@/shared/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select";
import { FormField } from "./FormField";
import { cn } from "@/shared/lib/utils";

export interface AddressData {
  street: string;
  city: string;
  state: string;
  zip: string;
  country: string;
}

interface AddressFormProps {
  value: AddressData;
  onChange: (address: AddressData) => void;
  errors?: Partial<Record<keyof AddressData, string>>;
  className?: string;
  disabled?: boolean;
}

const COUNTRIES = [
  { value: "US", label: "United States" },
  { value: "CA", label: "Canada" },
  { value: "GB", label: "United Kingdom" },
  { value: "AU", label: "Australia" },
  { value: "DE", label: "Germany" },
  { value: "FR", label: "France" },
  { value: "JP", label: "Japan" },
  { value: "IN", label: "India" },
];

export function AddressForm({
  value,
  onChange,
  errors,
  className,
  disabled,
}: AddressFormProps) {
  const handleChange = (field: keyof AddressData, fieldValue: string) => {
    onChange({ ...value, [field]: fieldValue });
  };

  return (
    <div className={cn("space-y-4", className)}>
      <FormField
        label="Street Address"
        htmlFor="street"
        error={errors?.street}
        required
      >
        <Input
          id="street"
          placeholder="123 Main St"
          value={value.street}
          onChange={(e) => handleChange("street", e.target.value)}
          error={!!errors?.street}
          disabled={disabled}
        />
      </FormField>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField
          label="City"
          htmlFor="city"
          error={errors?.city}
          required
        >
          <Input
            id="city"
            placeholder="City"
            value={value.city}
            onChange={(e) => handleChange("city", e.target.value)}
            error={!!errors?.city}
            disabled={disabled}
          />
        </FormField>

        <FormField
          label="State / Province"
          htmlFor="state"
          error={errors?.state}
          required
        >
          <Input
            id="state"
            placeholder="State"
            value={value.state}
            onChange={(e) => handleChange("state", e.target.value)}
            error={!!errors?.state}
            disabled={disabled}
          />
        </FormField>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField
          label="ZIP / Postal Code"
          htmlFor="zip"
          error={errors?.zip}
          required
        >
          <Input
            id="zip"
            placeholder="12345"
            value={value.zip}
            onChange={(e) => handleChange("zip", e.target.value)}
            error={!!errors?.zip}
            disabled={disabled}
          />
        </FormField>

        <FormField
          label="Country"
          htmlFor="country"
          error={errors?.country}
          required
        >
          <Select
            value={value.country}
            onValueChange={(val) => handleChange("country", val)}
            disabled={disabled}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select country" />
            </SelectTrigger>
            <SelectContent>
              {COUNTRIES.map((country) => (
                <SelectItem key={country.value} value={country.value}>
                  {country.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </FormField>
      </div>
    </div>
  );
}
