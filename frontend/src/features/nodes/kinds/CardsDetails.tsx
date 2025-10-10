import type { NodeVM } from "@/lib/playbook";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useState } from "react";
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
      {current && (
        <div className="text-sm text-muted-foreground whitespace-pre-line">
          {t(current.label, current.key)}
        </div>
      )}
      <div className="flex gap-2">
        <Button
          variant="secondary"
          onClick={() => setIndex((i) => Math.max(0, i - 1))}
          disabled={disabled || index === 0}
        >
          Prev
        </Button>
        {!last ? (
          <Button onClick={() => setIndex((i) => Math.min(items.length - 1, i + 1))} disabled={disabled}>
            Next
          </Button>
        ) : (
          <Button onClick={() => onSubmit?.({ acknowledged: true })} disabled={disabled}>
            Finish
          </Button>
        )}
      </div>
    </Card>
  );
}

