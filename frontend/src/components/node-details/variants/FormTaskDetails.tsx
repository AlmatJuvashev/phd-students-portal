// components/node-details/variants/FormTaskDetails.tsx
import { Card } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { NodeVM, FieldDef, t } from "@/lib/playbook";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { AssetsDownloads } from "../AssetsDownloads";

type Props = {
  node: NodeVM;
  onSubmit?: (payload: any) => void;
  initial?: Record<string, any>;
  canEdit?: boolean;
};

export function FormTaskDetails({
  node,
  initial = {},
  onSubmit,
  canEdit = true,
}: Props) {
  const [values, setValues] = useState<Record<string, any>>(initial);
  useEffect(() => {
    setValues(initial ?? {});
  }, [initial]);
  const { t: T } = useTranslation("common");

  const fields: FieldDef[] = node.requirements?.fields ?? [];

  function setField(k: string, v: any) {
    setValues((prev) => ({ ...prev, [k]: v }));
  }

  return (
    <Card className="p-4 space-y-4">
      {node.requirements?.notes && (
        <p className="text-sm text-muted-foreground">
          {node.requirements.notes}
        </p>
      )}
      <div className="space-y-3">
        {fields.map((f) => (
          <div key={f.key} className="grid gap-1">
            <Label htmlFor={f.key}>
              {t(f.label, f.key)}{" "}
              {f.required ? <span className="text-destructive">*</span> : null}
            </Label>
            {f.type === "textarea" || f.type === "array" ? (
              <Textarea
                id={f.key}
                disabled={!canEdit}
                placeholder={
                  f.type === "array"
                    ? T("forms.array_hint")
                    : t(f.placeholder, "")
                }
                value={values[f.key] ?? ""}
                onChange={(e) => setField(f.key, e.target.value)}
              />
            ) : (
              <Input
                id={f.key}
                disabled={!canEdit}
                type={f.type === "number" ? "number" : "text"}
                placeholder={t(f.placeholder, "")}
                value={values[f.key] ?? ""}
                onChange={(e) => setField(f.key, e.target.value)}
              />
            )}
          </div>
        ))}
      </div>

      {/* Templates / Downloads (if any) */}
      <AssetsDownloads node={node} />

      {!!node.requirements?.validations?.length && (
        <>
          <Separator />
          <div>
            <div className="mb-2 font-medium">
              {T("forms.validations_title")}
            </div>
            <ul className="list-inside list-disc text-sm">
              {node.requirements.validations!.map((v, i) => (
                <li key={i}>
                  {v.rule}
                  {v.source ? ` @ ${v.source}` : ""}
                </li>
              ))}
            </ul>
          </div>
        </>
      )}

      {canEdit && (
        <div className="flex gap-2">
          <Button onClick={() => onSubmit?.(values)}>
            {T("forms.save_submit")}
          </Button>
          <Button
            variant="secondary"
            onClick={() => onSubmit?.({ ...values, __draft: true })}
          >
            {T("forms.save_draft")}
          </Button>
        </div>
      )}
    </Card>
  );
}
