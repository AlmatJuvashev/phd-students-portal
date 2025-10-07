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
import { assetsForNode } from "@/lib/assets";

type Props = {
  node: NodeVM;
  onSubmit?: (payload: any) => void;
  initial?: Record<string, any>;
  canEdit?: boolean;
  disabled?: boolean;
};

export function FormTaskDetails({
  node,
  initial = {},
  onSubmit,
  canEdit = true,
  disabled = false,
}: Props) {
  const [values, setValues] = useState<Record<string, any>>(initial);
  const [showOmidConfirm, setShowOmidConfirm] = useState(false);
  useEffect(() => {
    setValues(initial ?? {});
  }, [initial]);
  const { t: T, i18n } = useTranslation("common");

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
        <div className="space-y-2">
          {node.id === "S1_publications_list" ? (
            <>
              <div className="text-sm font-medium">
                {T("forms.app7_prompt")}{" "}
                <span className="text-destructive">*</span>
              </div>
              <div className="flex gap-2">
                <Button onClick={() => setShowOmidConfirm(true)} disabled={disabled}>
                  {T("forms.yes")}
                </Button>
                <Button
                  variant="secondary"
                  disabled={disabled}
                  onClick={() => {
                    // open preferred Appendix 7 template in a new tab
                    const assets = assetsForNode(node);
                    const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
                    const preferred =
                      assets.find(
                        (a) =>
                          a.id.toLowerCase().includes("app7") &&
                          a.id.toLowerCase().includes(`_${lang}`)
                      ) ||
                      assets.find((a) => a.id.toLowerCase().includes("app7")) ||
                      assets[0];
                    if (preferred?.storage?.key) {
                      window.open(
                        `/${preferred.storage.key}`,
                        "_blank",
                        "noopener,noreferrer"
                      );
                    }
                    // scroll to templates section to draw attention
                    document
                      .getElementById("templates-section")
                      ?.scrollIntoView({ behavior: "smooth", block: "start" });
                    // mark as draft so the user can return after preparing the document
                    onSubmit?.({ ...values, __draft: true });
                  }}
                >
                  {T("forms.no")}
                </Button>
              </div>
            </>
          ) : node.id === "E1_apply_omid" ? (
            <>
              <div className="text-sm font-medium">
                {T("forms.omid_prompt")} {" "}
                <span className="text-destructive">*</span>
              </div>
              <div className="flex gap-2">
                <Button onClick={() => onSubmit?.(values)} disabled={disabled}>
                  {T("forms.yes")}
                </Button>
                <Button
                  variant="secondary"
                  disabled={disabled}
                  onClick={() => {
                    // Open OMiD application template matching current locale
                    const assets = assetsForNode(node);
                    const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
                    const preferred =
                      assets.find(
                        (a) =>
                          a.id.toLowerCase().includes("omid") &&
                          a.id.toLowerCase().includes(`_${lang}`)
                      ) ||
                      assets.find((a) => a.id.toLowerCase().includes("omid")) ||
                      assets[0];
                    if (preferred?.storage?.key) {
                      window.open(`/${preferred.storage.key}`, "_blank", "noopener,noreferrer");
                    }
                    // Scroll to templates block and save draft
                    document.getElementById("templates-section")?.scrollIntoView({ behavior: "smooth", block: "start" });
                    onSubmit?.({ ...values, __draft: true });
                  }}
                >
                  {T("forms.no")}
                </Button>
              </div>
            </>
          ) : (
            <div className="flex gap-2">
              <Button onClick={() => onSubmit?.(values)} disabled={disabled}>
                {T("forms.save_submit")}
              </Button>
              <Button
                variant="secondary"
                disabled={disabled}
                onClick={() => onSubmit?.({ ...values, __draft: true })}
              >
                {T("forms.save_draft")}
              </Button>
            </div>
          )}
        </div>
      )}
      {showOmidConfirm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md shadow-lg">
            <div className="mb-4 text-sm text-muted-foreground whitespace-pre-line">
              {T("forms.omid_info_after_yes")}
            </div>
            <div className="flex gap-2 justify-end">
              <Button
                variant="outline"
                onClick={() => setShowOmidConfirm(false)}
              >
                {T("common.cancel")}
              </Button>
              <Button
                onClick={() => {
                  setShowOmidConfirm(false);
                  onSubmit?.(values);
                }}
                disabled={disabled}
              >
                {T("forms.proceed_next")}
              </Button>
            </div>
          </div>
        </div>
      )}
    </Card>
  );
}
