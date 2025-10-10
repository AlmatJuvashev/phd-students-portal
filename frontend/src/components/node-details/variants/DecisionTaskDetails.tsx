import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { NodeVM, t } from "@/lib/playbook";
import { ExternalLink, FileText, CheckCircle2 } from "lucide-react";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

function firstUrl(text?: string) {
  if (!text) return null;
  const m = text.match(/https?:\/\/\S+/i);
  return m ? m[0] : null;
}

export function DecisionTaskDetails({
  node,
  onSubmit,
  disabled = false,
}: {
  node: NodeVM;
  onSubmit?: () => void;
  disabled?: boolean;
}) {
  const { t: T } = useTranslation("common");
  const desc = t(node.description as any, node.id);
  const url = useMemo(() => firstUrl(desc), [desc]);

  const checklist = [
    "title page, abstract, table of contents",
    "introduction, objectives, novelty",
    "chapters with methods and results",
    "conclusion (findings, contribution)",
    "references list (GOST/APA per rules)",
    "appendices (if applicable)",
  ];

  return (
    <Card className="p-4 space-y-4">
      <div className="flex items-start gap-3">
        <div className="mt-0.5 rounded-full bg-primary/10 p-2 text-primary">
          <FileText className="h-5 w-5" />
        </div>
        <div className="space-y-2">
          <div className="text-sm text-muted-foreground whitespace-pre-wrap">
            {desc}
          </div>
          {url && (
            <div className="flex items-center gap-2">
              <a
                className="inline-flex items-center gap-1 underline text-primary"
                href={url}
                target="_blank"
                rel="noreferrer"
              >
                <ExternalLink className="h-4 w-4" />
                <span>Open standard</span>
              </a>
            </div>
          )}
        </div>
      </div>

      <div className="rounded-md border bg-blue-50/60 dark:bg-blue-950/20 border-blue-200 dark:border-blue-800 p-3">
        <div className="mb-1 text-xs font-semibold text-blue-700 dark:text-blue-300">
          What to check
        </div>
        <ul className="list-disc pl-5 text-xs text-blue-800 dark:text-blue-200 space-y-0.5">
          {checklist.map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      </div>

      <div className="flex flex-wrap gap-2">
        {url && (
          <a href={url} target="_blank" rel="noreferrer">
            <Button variant="outline" type="button">
              <ExternalLink className="mr-1 h-4 w-4" /> Open standard
            </Button>
          </a>
        )}
        <Button disabled={disabled} aria-busy={disabled} onClick={onSubmit}>
          <CheckCircle2 className="mr-1 h-4 w-4" /> {T("decision.submit")}
        </Button>
      </div>
    </Card>
  );
}
