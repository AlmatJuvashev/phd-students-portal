// components/node-details/variants/GatewayInfoDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { NodeVM } from "@/lib/playbook";

export function GatewayInfoDetails({
  node,
  onContinue,
}: {
  node: NodeVM;
  onContinue?: () => void;
}) {
  return (
    <Card className="p-4 space-y-3">
      <div className="text-sm">Системный переходный узел.</div>
      {node.condition && (
        <div className="text-sm text-muted-foreground">
          Условие: <code>{node.condition}</code>
        </div>
      )}
      <Button onClick={() => onContinue?.()}>Продолжить</Button>
    </Card>
  );
}
