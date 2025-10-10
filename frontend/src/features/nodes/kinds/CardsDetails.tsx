import type { NodeVM } from "@/lib/playbook";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import StickyActions from "@/components/ui/sticky-actions";
import { t } from "@/lib/playbook";

// Lightweight sequential cards skeleton; can be extended later
export default function CardsDetails({
  node,
  onSubmit,
  disabled,
}: {
  node: NodeVM;
  onSubmit?: (payload: any) => void;
  disabled?: boolean;
}) {
  const items = (node.requirements as any)?.fields?.filter((f: any) => f?.type === "note") || [];
  const [index, setIndex] = useState(0);
  const last = index >= items.length - 1;
  const current = items[index];

  return (
    <Card className="p-4 space-y-4">
      {/* Progress dots */}
      <div className="flex items-center gap-1">
        {items.map((_, i) => (
          <div
            key={i}
            className={`h-2 w-2 rounded-full ${i === index ? "bg-primary" : "bg-muted"}`}
          />
        ))}
        <div className="ml-auto text-xs text-muted-foreground">
          {index + 1}/{items.length}
        </div>
      </div>
      {/* Animated content */}
      <div className="min-h-[120px]">
        <AnimatePresence initial={false}>
          <motion.div
            key={current?.key || index}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -10 }}
            transition={{ duration: 0.2 }}
          >
            {current && (
              <div className="text-sm text-muted-foreground whitespace-pre-line">
                {t(current.label, current.key)}
              </div>
            )}
          </motion.div>
        </AnimatePresence>
      </div>
      {/* Sticky actions for mobile */}
      <StickyActions
        secondaryLabel="Prev"
        onSecondary={() => setIndex((i) => Math.max(0, i - 1))}
        primaryLabel={last ? "Finish" : "Next"}
        onPrimary={() =>
          last
            ? onSubmit?.({ acknowledged: true })
            : setIndex((i) => Math.min(items.length - 1, i + 1))
        }
        disabled={disabled}
      />
    </Card>
  );
}
