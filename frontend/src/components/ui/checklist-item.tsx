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
  // ReadOnly mode - show all items but with different styling
  if (readOnly) {
    if (checked) {
      // Show checked items with green styling
      return (
        <div
          className="flex items-start gap-3 p-4 rounded-xl bg-green-50 dark:bg-green-900/10 border-2 border-green-200 dark:border-green-800/30 transition-all duration-200 min-w-0"
          style={{ contain: "layout" }}
        >
          <div className="flex-shrink-0 w-6 h-6 rounded-full bg-green-500 flex items-center justify-center mt-0.5">
            <Check className="w-4 h-4 text-white" strokeWidth={3} />
          </div>
          <span className="text-sm text-green-900 dark:text-green-100 leading-relaxed flex-1 min-w-0">
            {label}
          </span>
        </div>
      );
    } else {
      // Show unchecked items as grayed out but still visible
      return (
        <div
          className="flex items-start gap-3 p-4 rounded-xl bg-gray-50 dark:bg-gray-900/10 border-2 border-gray-200 dark:border-gray-800/30 opacity-60 transition-all duration-200 min-w-0"
          style={{ contain: "layout" }}
        >
          <div className="flex-shrink-0 w-6 h-6 rounded-full border-2 border-gray-300 dark:border-gray-600 mt-0.5"></div>
          <span className="text-sm text-gray-700 dark:text-gray-400 leading-relaxed flex-1 min-w-0">
            {label}
          </span>
        </div>
      );
    }
  }

  return (
    <label
      className={cn(
        "group flex items-start gap-3 p-4 rounded-xl border-2 transition-all duration-200 cursor-pointer min-w-0",
        {
          "bg-gradient-to-r from-primary/5 to-primary/10 border-primary/30 hover:border-primary/50":
            checked && !disabled,
          "bg-card border-border hover:border-primary/30 hover:bg-muted/30":
            !checked && !disabled,
          "opacity-60 cursor-not-allowed": disabled,
        }
      )}
      style={{ contain: "layout" }}
    >
      <div className="flex-1 leading-relaxed min-w-0">
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
            onChange?.(e.target.checked);
          }}
          disabled={disabled || readOnly}
          className="sr-only"
        />
      </div>
    </label>
  );
}
