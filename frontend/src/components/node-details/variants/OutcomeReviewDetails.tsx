// components/node-details/variants/OutcomeReviewDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Separator } from "@/components/ui/separator";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { NodeVM, t } from "@/lib/playbook";
import { useState } from "react";
import { UploadTaskDetails } from "./UploadTaskDetails";

type Props = {
  node: NodeVM;
  canDecide?: boolean;
  canUpload?: boolean;
  onFinalize?: (payload: {
    outcome: string;
    note?: string;
    files?: Record<string, File | null>;
  }) => void;
};

export function OutcomeReviewDetails({
  node,
  canDecide = true,
  canUpload = true,
  onFinalize,
}: Props) {
  const outcomes = node.outcomes?.length
    ? node.outcomes
    : ([{ value: "accepted" }, { value: "fixes_required" }] as const);
  const [value, setValue] = useState(outcomes?.[0]?.value ?? "accepted");
  const [note, setNote] = useState("");
  const [files, setFiles] = useState<Record<string, File | null>>({});

  const requiresUploads = !!node.requirements?.uploads?.length;

  return (
    <Card className="p-4 space-y-4">
      <div>
        <div className="mb-2 font-medium">Решение</div>
        <RadioGroup
          value={value}
          onValueChange={setValue}
          className="grid gap-2 sm:grid-cols-2"
        >
          {outcomes!.map((o) => (
            <div
              key={o.value}
              className="flex items-center space-x-2 rounded-md border p-3"
            >
              <RadioGroupItem value={o.value} id={o.value} />
              <Label htmlFor={o.value} className="cursor-pointer">
                {t(
                    'label' in o ? o.label : undefined,
                    o.value
                )}
              </Label>
            </div>
          ))}
        </RadioGroup>
      </div>

      <div className="grid gap-1">
        <Label htmlFor="note">Комментарий (опционально)</Label>
        <Textarea
          id="note"
          value={note}
          onChange={(e) => setNote(e.target.value)}
          placeholder="Краткое обоснование решения…"
        />
      </div>

      {requiresUploads && (
        <>
          <Separator />
          <div className="font-medium">Материалы заседания</div>
          <UploadTaskDetails
            node={node}
            canEdit={canUpload}
            onSubmit={({ files: f }) => setFiles(f)}
          />
        </>
      )}

      {canDecide && (
        <Button onClick={() => onFinalize?.({ outcome: value, note, files })}>
          Зафиксировать решение
        </Button>
      )}
    </Card>
  );
}
