// components/node-details/variants/InfoDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { NodeVM, t } from "@/lib/playbook";
import { useTranslation } from "react-i18next";

export default function InfoDetails({
  node,
  onContinue,
}: {
  node: NodeVM;
  onContinue?: () => void;
}) {
  const { t: T } = useTranslation("common");
  const fields = node.requirements?.fields ?? [];
  const ui = (node.requirements as any)?.ui_hints || {};
  const layout = ui.cards_layout || {};

  // Prefer next pointer from node.next when present
  const hasNext = Array.isArray(node.next) && node.next.length > 0;

  const notes = fields.filter((f: any) => f.type === "note");

  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 h-full">
      {/* Left informational panel (occupies full width since there's no right panel) */}
      <div className="lg:col-span-5 min-h-0 overflow-auto space-y-3">
        {/* Description */}
        {Boolean((node as any).description) && (
          <Card className="p-4">
            <div className="text-sm text-muted-foreground">
              {t((node as any).description, "")}
            </div>
          </Card>
        )}
        {/* Notes as stacked cards when requested */}
        {layout?.style === "stacked" ? (
          <div className="space-y-3">
            {notes.map((f) => (
              <Card key={f.key} className="p-4 bg-gray-50">
                <div className="text-sm text-muted-foreground">
                  {t(f.label, f.key)}
                </div>
              </Card>
            ))}
          </div>
        ) : (
          <Card className="p-4">
            <div className="space-y-2">
              {notes.map((f) => (
                <div key={f.key} className="text-sm text-muted-foreground">
                  {t(f.label, f.key)}
                </div>
              ))}
            </div>
          </Card>
        )}
        {/* Continue button if explicit next exists */}
        {hasNext && (
          <div className="pt-2">
            <Button onClick={() => onContinue?.()}>
              {T("forms.proceed_next")}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
