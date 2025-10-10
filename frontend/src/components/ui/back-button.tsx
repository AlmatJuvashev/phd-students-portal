import { ArrowLeft } from "lucide-react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface BackButtonProps {
  to?: string;
  label?: string;
  showLabelOnMobile?: boolean;
  className?: string;
  variant?: "default" | "ghost";
}

export function BackButton({
  to = "/",
  label,
  showLabelOnMobile = false,
  className,
  variant = "default",
}: BackButtonProps) {
  const { t: T } = useTranslation("common");
  const buttonLabel = label || T("common.back", { defaultValue: "Back" });

  const baseStyles =
    "inline-flex items-center justify-center gap-2 px-3 sm:px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 active:scale-95";

  const variantStyles = {
    default:
      "bg-muted/50 hover:bg-muted text-foreground border border-border hover:border-primary/40 shadow-sm hover:shadow-md",
    ghost: "hover:bg-muted/50 text-foreground",
  };

  return (
    <Link to={to} className={cn("inline-block", className)}>
      <button
        className={cn(baseStyles, variantStyles[variant])}
        aria-label={buttonLabel}
      >
        <ArrowLeft className="w-4 h-4 flex-shrink-0" />
        <span className={cn(showLabelOnMobile ? "" : "hidden sm:inline")}>
          {buttonLabel}
        </span>
      </button>
    </Link>
  );
}
