// components/node-details/variants/WaitLockDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { NodeVM } from "@/lib/playbook";

export function WaitLockDetails({
  node,
  onSubscribe,
}: {
  node: NodeVM;
  onSubscribe?: () => void;
}) {
  const days = node.timer?.duration_days ?? 0;
  return (
    <Card className="p-4 space-y-3">
      <div className="text-sm text-muted-foreground">
        Этот этап заблокирован таймером.
      </div>
      <div className="text-lg font-semibold">Ожидание: {days} дней</div>
      <div className="text-sm">Старт: {node.timer?.start_on}</div>
      <Button variant="secondary" onClick={onSubscribe}>
        Напоминать об открытии
      </Button>
    </Card>
  );
}
