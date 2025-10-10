import * as React from "react";
import { cn } from "../../lib/utils";
export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {}
export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, ...props }, ref) => (
    <input
      ref={ref}
      className={cn(
        "flex h-10 w-full rounded-lg border-2 border-input bg-background px-4 py-2 text-sm shadow-sm transition-all duration-200",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:border-primary",
        "hover:border-primary/40",
        "disabled:cursor-not-allowed disabled:opacity-50 disabled:bg-muted",
        "placeholder:text-muted-foreground",
        className
      )}
      {...props}
    />
  )
);
Input.displayName = "Input";
