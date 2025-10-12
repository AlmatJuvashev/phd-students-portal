import { Check } from "lucide-react";
import { cn } from "@/lib/utils";

interface ChecklistItemProps {
  checked: boolean;
  onChange?: (checked: boolean) => void;
  label: string;
  required?: boolean;
  disabled?: boolean;
  readOnly?: boolean;
}

export function ChecklistItem({
  checked,
  onChange,
  label,
  required = false,
  disabled = false,
  readOnly = false,
}: ChecklistItemProps) {
  // ReadOnly mode - show only completed items with green checkmarks
  if (readOnly) {
    // Only render checked items in readOnly mode
    if (checked) {
      return (
        <div className="flex items-start gap-3 p-4 rounded-xl bg-green-50 dark:bg-green-900/10 border-2 border-green-200 dark:border-green-800/30 transition-all duration-200">
          <div className="flex-shrink-0 w-6 h-6 rounded-full bg-green-500 flex items-center justify-center mt-0.5">
            <Check className="w-4 h-4 text-white" strokeWidth={3} />
          </div>
          <span className="text-sm text-green-900 dark:text-green-100 leading-relaxed flex-1">
            {label}
          </span>
        </div>
      );
    } else {
      // Don't render unchecked items in readOnly mode
      return null;
    }
  }

  return (
    <label
      className={cn(
        "group flex items-start gap-3 p-4 rounded-xl border-2 transition-all duration-200 cursor-pointer",
        {
          "bg-gradient-to-r from-primary/5 to-primary/10 border-primary/30 hover:border-primary/50":
            checked && !disabled,
          "bg-card border-border hover:border-primary/30 hover:bg-muted/30":
            !checked && !disabled,
          "opacity-60 cursor-not-allowed": disabled,
        }
      )}
    >
      <div className="flex-1 leading-relaxed">
        <span className="text-sm text-foreground">
          {label}
          {required && (
            <span className="ml-1 text-destructive font-semibold">*</span>
          )}
        </span>
      </div>

      <div className="flex-shrink-0">
        <div
          className={cn(
            "w-6 h-6 rounded-md border-2 flex items-center justify-center transition-all duration-200",
            {
              "bg-primary border-primary group-hover:scale-110": checked,
              "bg-background border-input group-hover:border-primary/50":
                !checked,
            }
          )}
        >
          {checked && (
            <Check
              className="w-4 h-4 text-primary-foreground"
              strokeWidth={3}
            />
          )}
        </div>
        <input
          type="checkbox"
          checked={checked}
          onChange={(e) => {
            e.preventDefault();
            onChange?.(e.target.checked);
          }}
          disabled={disabled || readOnly}
          className="sr-only"
        />
      </div>
    </label>
  );
}
