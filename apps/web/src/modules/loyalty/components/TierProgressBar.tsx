interface TierProgressBarProps {
  currentTier: string;
  nextTier: string | null;
  progress: number;
  pointsToNext: number;
}

export function TierProgressBar({
  currentTier,
  nextTier,
  progress,
  pointsToNext,
}: TierProgressBarProps) {
  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between text-sm">
        <span className="font-medium capitalize">{currentTier}</span>
        {nextTier && <span className="text-muted-foreground capitalize">{nextTier}</span>}
      </div>
      <div className="h-3 overflow-hidden rounded-full bg-muted">
        <div
          className="h-full rounded-full bg-gradient-to-r from-primary to-primary/80 transition-all"
          style={{ width: `${Math.min(progress, 100)}%` }}
        />
      </div>
      {nextTier ? (
        <p className="text-xs text-muted-foreground">
          {pointsToNext.toLocaleString()} points to reach {nextTier}
        </p>
      ) : (
        <p className="text-xs text-muted-foreground">
          You have reached the highest tier!
        </p>
      )}
    </div>
  );
}
