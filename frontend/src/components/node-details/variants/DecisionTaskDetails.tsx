import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { NodeVM, t } from "@/lib/playbook";
import { ExternalLink, FileText, CheckCircle2 } from "lucide-react";
import StickyActions from "@/components/ui/sticky-actions";
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
  renderGuide,
}: {
  node: NodeVM;
  onSubmit?: () => void;
  disabled?: boolean;
  renderGuide?: () => React.ReactNode;
}) {
  const { t: T } = useTranslation("common");
  const desc = t(node.description as any, node.id);
  const url = useMemo(() => firstUrl(desc), [desc]);

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
          {/* Inline guide for mobile-first UX */}
          {renderGuide ? renderGuide() : null}
        </div>
      </div>
      <StickyActions
        primaryLabel={T("decision.submit")}
        onPrimary={() => onSubmit?.()}
        disabled={disabled}
        busy={disabled}
      />
    </Card>
  );
}
