import { Card } from "@/components/ui/card";

export function GuideCard({
  title,
  items,
  body,
  href,
  ctaLabel,
  onCta,
  tone = "info",
}: {
  title?: string;
  items?: string[];
  body?: React.ReactNode;
  href?: string;
  ctaLabel?: string;
  onCta?: () => void;
  tone?: "info" | "success" | "warning";
}) {
  const toneClasses = {
    info: "bg-blue-50/60 dark:bg-blue-950/20 border-blue-200 dark:border-blue-800 text-blue-800 dark:text-blue-200",
    success:
      "bg-emerald-50/60 dark:bg-emerald-950/20 border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200",
    warning:
      "bg-amber-50/60 dark:bg-amber-950/20 border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-200",
  } as const;
  return (
    <Card className={`p-3 rounded-md border ${toneClasses[tone]}`}>
      {title ? (
        <div className="mb-1 text-sm font-semibold">{title}</div>
      ) : null}
      {items?.length ? (
        <ul className="list-disc pl-5 text-sm space-y-1">
          {items.map((it, i) => (
            <li key={i}>{it}</li>
          ))}
        </ul>
      ) : null}
      {body ? <div className="text-sm whitespace-pre-line">{body}</div> : null}
      {href || onCta ? (
        <div className="pt-2">
          {href ? (
            <a
              className="inline-flex items-center underline text-primary"
              href={href}
              target="_blank"
              rel="noreferrer"
            >
              {ctaLabel || "Open"}
            </a>
          ) : (
            <button className="underline text-primary" onClick={onCta}>
              {ctaLabel || "Open"}
            </button>
          )}
        </div>
      ) : null}
    </Card>
  );
}

export default GuideCard;

