// components/node-details/variants/InfoDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { NodeVM, t } from "@/lib/playbook";
import { useTranslation } from "react-i18next";
import StickyActions from "@/components/ui/sticky-actions";

export default function InfoDetails({
  node,
  onContinue,
  renderGuide,
}: {
  node: NodeVM;
  onContinue?: () => void;
  renderGuide?: () => React.ReactNode;
}) {
  const { t: T } = useTranslation("common");
  const fields = node.requirements?.fields ?? [];
  const ui = (node.requirements as any)?.ui_hints || {};
  const layout = ui.cards_layout || {};

  // Prefer next pointer from node.next when present
  const hasNext = Array.isArray(node.next) && node.next.length > 0;

  const notes = fields.filter((f: any) => f.type === "note");

  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 sm:gap-6 h-full">
      {/* Left informational panel (occupies full width since there's no right panel) */}
      <div className="lg:col-span-5 min-h-0 overflow-auto space-y-4">
        {/* Description */}
        {Boolean((node as any).description) && (
          <Card className="p-5 sm:p-6 shadow-sm hover:shadow-md transition-shadow duration-300 border-l-4 border-primary/30 bg-gradient-to-br from-card to-card/50">
            <div className="text-sm sm:text-base leading-relaxed text-muted-foreground">
              {t((node as any).description, "")}
            </div>
          </Card>
        )}
        {/* Optional standardized guide */}
        {renderGuide ? renderGuide() : null}
        {/* Notes as stacked cards when requested */}
        {layout?.style === "stacked" ? (
          layout?.card_variant === "info" ? (
            <Card className="p-5 sm:p-6 bg-gradient-to-br from-primary/5 via-primary/3 to-transparent border border-primary/20 shadow-sm hover:shadow-md transition-all duration-300 animate-in fade-in slide-in-from-bottom-2">
              <div className="space-y-4">
                {notes.map((f, idx) => {
                  const raw = t(f.label, f.key);
                  // Fix: String indexing/spread
                  const first = [...raw][0] || "";
                  const rest = [...raw].slice(1).join("");
                  const leadingIcon =
                    first && /[\p{Emoji}\p{Extended_Pictographic}]/u.test(first)
                      ? first
                      : "";
                  const body = leadingIcon ? rest : raw;
                  return (
                    <div key={f.key} className="space-y-3">
                      <div className="flex items-start gap-3">
                        <div className="flex-shrink-0 text-xl mt-0.5">
                          {leadingIcon || idx + 1}
                        </div>
                        <div className="text-sm sm:text-base leading-relaxed text-muted-foreground flex-1 whitespace-pre-wrap">
                          {leadingIcon ? body.trimStart() : raw}
                        </div>
                      </div>
                      {idx < notes.length - 1 && (
                        <div className="h-px bg-gradient-to-r from-transparent via-primary/25 to-transparent opacity-80" />
                      )}
                    </div>
                  );
                })}
              </div>
            </Card>
          ) : (
            <div className="space-y-3 sm:space-y-4">
              {notes.map((f, idx) => (
                <Card
                  key={f.key}
                  className="p-4 sm:p-5 bg-gradient-to-br from-muted/40 to-muted/20 border border-border/60 shadow-sm hover:shadow-md hover:border-primary/30 transition-all duration-300 animate-in fade-in slide-in-from-bottom-2"
                  style={{ animationDelay: `${idx * 50}ms` }}
                >
                  <div className="flex items-start gap-3">
                    <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold text-sm">
                      {idx + 1}
                    </div>
                    <div className="text-sm sm:text-base leading-relaxed text-muted-foreground flex-1">
                      {t(f.label, f.key)}
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          )
        ) : (
          <Card className="p-5 sm:p-6 shadow-sm hover:shadow-md transition-shadow duration-300 border border-border/60">
            <div className="space-y-3">
              {notes.map((f, idx) => (
                <div
                  key={f.key}
                  className="flex items-start gap-3 pb-3 border-b border-border/40 last:border-0 last:pb-0"
                >
                  <div className="flex-shrink-0 w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold text-xs">
                    {idx + 1}
                  </div>
                  <div className="text-sm sm:text-base leading-relaxed text-muted-foreground flex-1">
                    {t(f.label, f.key)}
                  </div>
                </div>
              ))}
            </div>
          </Card>
        )}
        {/* Continue button if explicit next exists */}
        {hasNext && (
          <StickyActions
            primaryLabel={T("forms.proceed_next")}
            onPrimary={() => onContinue?.()}
          />
        )}
      </div>
    </div>
  );
}
