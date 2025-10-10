import React from "react";

export function StickyActions({
  primaryLabel,
  onPrimary,
  secondaryLabel,
  onSecondary,
  disabled,
  busy,
}: {
  primaryLabel: string;
  onPrimary: () => void;
  secondaryLabel?: string;
  onSecondary?: () => void;
  disabled?: boolean;
  busy?: boolean;
}) {
  return (
    <div className="fixed bottom-0 left-0 right-0 z-40 border-t bg-background/90 backdrop-blur supports-[backdrop-filter]:bg-background/60 md:static md:bg-transparent md:backdrop-blur-0 md:border-0">
      <div className="mx-auto max-w-6xl p-3 flex gap-2 justify-end">
        {secondaryLabel && onSecondary ? (
          <button
            type="button"
            className="inline-flex items-center justify-center h-11 rounded-md border px-4 text-sm bg-muted hover:bg-muted/80"
            onClick={onSecondary}
            disabled={disabled || busy}
            aria-busy={busy}
          >
            {secondaryLabel}
          </button>
        ) : null}
        <button
          type="button"
          className="inline-flex items-center justify-center h-11 rounded-md bg-primary px-4 text-sm text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
          onClick={onPrimary}
          disabled={disabled || busy}
          aria-busy={busy}
        >
          {primaryLabel}
        </button>
      </div>
    </div>
  );
}

export default StickyActions;

