import * as React from "react";
import { cn } from "../../lib/utils";
export interface TextareaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {}
export const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, ...props }, ref) => (
    <textarea
      ref={ref}
      className={cn(
        "flex min-h-[100px] w-full rounded-lg border-2 border-input bg-background px-4 py-3 text-sm shadow-sm transition-all duration-200",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:border-primary",
        "hover:border-primary/40",
        "disabled:cursor-not-allowed disabled:opacity-50 disabled:bg-muted",
        "placeholder:text-muted-foreground",
        "resize-y",
        className
      )}
      {...props}
    />
  )
);
Textarea.displayName = "Textarea";
