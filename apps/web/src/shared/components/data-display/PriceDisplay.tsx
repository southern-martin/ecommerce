import { cn } from "@/shared/lib/utils";
import { formatPrice } from "@/shared/lib/utils";

interface PriceDisplayProps {
  priceInCents: number;
  compareAtPriceInCents?: number;
  currency?: string;
  className?: string;
  size?: "sm" | "md" | "lg";
}

const sizeClasses = {
  sm: "text-sm",
  md: "text-base",
  lg: "text-2xl",
};

export function PriceDisplay({
  priceInCents,
  compareAtPriceInCents,
  currency = "USD",
  className,
  size = "md",
}: PriceDisplayProps) {
  const hasDiscount =
    compareAtPriceInCents != null && compareAtPriceInCents > priceInCents;

  const discountPercent = hasDiscount
    ? Math.round(
        ((compareAtPriceInCents - priceInCents) / compareAtPriceInCents) * 100
      )
    : 0;

  return (
    <div className={cn("flex items-center gap-2", className)}>
      <span
        className={cn(
          "font-semibold",
          sizeClasses[size],
          hasDiscount && "text-destructive"
        )}
      >
        {formatPrice(priceInCents, currency)}
      </span>

      {hasDiscount && (
        <>
          <span
            className={cn(
              "text-muted-foreground line-through",
              size === "lg" ? "text-base" : "text-sm"
            )}
          >
            {formatPrice(compareAtPriceInCents, currency)}
          </span>
          <span className="rounded bg-destructive/10 px-1.5 py-0.5 text-xs font-medium text-destructive">
            -{discountPercent}%
          </span>
        </>
      )}
    </div>
  );
}
