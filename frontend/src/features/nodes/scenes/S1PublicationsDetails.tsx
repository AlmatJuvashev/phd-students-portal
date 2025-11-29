import { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { assetsForNode } from "@/lib/assets";
import { t as pbT } from "@/lib/playbook";
import {
  FormTaskDetails,
  FormTaskDetailsProps,
  FormTaskContext,
} from "@/features/nodes/details/variants/FormTaskDetails";
import { generateApp7FromTemplate } from "@/features/docgen/app7-templated";
import { useAuth } from "@/contexts/AuthContext";
import { useProfileSnapshot } from "@/features/profile/useProfileSnapshot";
import { buildTemplateData } from "@/features/docgen/student-template";

const scrollToTemplates = () =>
  document
    .getElementById("templates-section")
    ?.scrollIntoView({ behavior: "smooth", block: "start" });

const SECTION_KEYS = ["wos_scopus", "kokson", "conferences", "ip"] as const;

type SectionKey = (typeof SECTION_KEYS)[number];

type SectionCount = {
  key: SectionKey;
  label: string;
  count: number;
};

export function S1PublicationsDetails({
  node,
  initial,
  canEdit,
  ...rest
}: Omit<FormTaskDetailsProps, "renderActions">) {
  const { t: T, i18n } = useTranslation("common");
  const assets = useMemo(() => assetsForNode(node), [node]);
  const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
  const { user } = useAuth();
  const { data: profileData } = useProfileSnapshot(user?.role === "student");

  const fieldMap = useMemo(() => {
    const map = new Map<string, any>();
    node.requirements?.fields?.forEach((field: any) => {
      if (field?.key) {
        map.set(field.key, field);
      }
    });
    return map;
  }, [node]);

  const computeCounts = useCallback(
    (values: Record<string, any> = {}) =>
      SECTION_KEYS.map<SectionCount>((key) => {
        const fieldDef = fieldMap.get(key);
        const label = pbT(fieldDef?.label, key) || key;
        const items = Array.isArray(values[key]) ? values[key] : [];
        return { key, label, count: items.length };
      }),
    [fieldMap]
  );

  const openTemplate = useCallback(() => {
    const preferred =
      assets.find(
        (asset) =>
          asset.id.toLowerCase().includes("app7") &&
          asset.id.toLowerCase().includes(`_${lang}`)
      ) ||
      assets.find((asset) => asset.id.toLowerCase().includes("app7")) ||
      assets[0];
    if (preferred?.storage?.key) {
      window.open(`/${preferred.storage.key}`, "_blank", "noopener,noreferrer");
    }
  }, [assets, lang]);

  const [validationErrors, setValidationErrors] = useState<Record<string, Record<number, string>>>({});

  const validatePublications = useCallback((values: Record<string, any>) => {
    const errors: Record<string, Record<number, string>> = {};
    let hasErrors = false;

    // DOI Regex: starts with 10., followed by 4-9 digits, slash, then allowed chars
    const doiRegex = /^10\.\d{4,9}\/[-._;()/:A-Z0-9]+$/i;

    SECTION_KEYS.forEach((key) => {
      const items = Array.isArray(values[key]) ? values[key] : [];
      const fieldDef = fieldMap.get(key);
      const itemFields = fieldDef?.item_fields || [];
      const sectionErrors: Record<number, string> = {};

      items.forEach((item: any, index: number) => {
        const itemErrors: string[] = [];

        itemFields.forEach((f: any) => {
          // Check required fields
          if (f.required && !item[f.key]) {
             // We won't show individual field errors on the card, just a general "Invalid"
             // But we could if we wanted to be more specific.
             // For now, let's just flag the item.
             if (!itemErrors.includes("Required fields missing")) {
                itemErrors.push("Required fields missing");
             }
          }

          // Check DOI format
          if (f.key === "doi" && item[f.key]) {
             if (!doiRegex.test(item[f.key].trim())) {
                itemErrors.push("Invalid DOI format");
             }
          }
        });

        if (itemErrors.length > 0) {
          sectionErrors[index] = itemErrors.join(", ");
          hasErrors = true;
        }
      });

      if (Object.keys(sectionErrors).length > 0) {
        errors[key] = sectionErrors;
      }
    });

    setValidationErrors(errors);
    return !hasErrors;
  }, [fieldMap]);

  const renderActions = useCallback(
    (ctx: FormTaskContext) => {
      if (!ctx.canEdit) return null;
      const counts = computeCounts(ctx.values || {});
      const total = counts.reduce((acc, item) => acc + item.count, 0);

      return (
        <>
          <div className="flex gap-2 mb-3">
            <Button onClick={() => ctx.submit()} disabled={ctx.disabled}>
              {T("forms.submit_publications", "Submit publications")}
            </Button>
            <Button
              variant="secondary"
              onClick={() => {
                const isValid = validatePublications(ctx.values);
                if (!isValid) {
                   // Don't generate if invalid, but keep data
                   return;
                }

                const templateData =
                  user?.role === "student"
                    ? buildTemplateData(user, profileData as any, lang)
                    : undefined;
                generateApp7FromTemplate(
                  ctx.values as any,
                  lang,
                  templateData
                ).catch((err) => console.error(err));
              }}
              disabled={ctx.disabled}
            >
              {T("forms.generate_app7")}
            </Button>
          </div>
          <div className="mb-3 space-y-2 rounded-md border border-dashed p-3 text-sm">
            <div className="font-medium">
              {T("forms.collection_summary", "Записи по разделам")}
            </div>
            <div className="grid gap-1">
              {counts.map(({ key, label, count }) => (
                <div
                  key={key}
                  className="flex items-center justify-between gap-2"
                >
                  <span className="text-muted-foreground">{label}</span>
                  <span className="font-semibold">{count}</span>
                </div>
              ))}
            </div>
            <div className="text-xs text-muted-foreground">
              {T("forms.collection_total", "Всего записей")}: {total}
            </div>
          </div>
          <div className="text-sm font-medium">
            {T("forms.app7_prompt")} <span className="text-destructive">*</span>
          </div>
          <div className="flex gap-2">
            <Button onClick={() => ctx.submit()} disabled={ctx.disabled}>
              {T("forms.yes")}
            </Button>
            <Button
              variant="secondary"
              disabled={ctx.disabled}
              onClick={() => {
                openTemplate();
                scrollToTemplates();
                ctx.saveDraft();
              }}
            >
              {T("forms.no")}
            </Button>
          </div>
        </>
      );
    },
    [T, computeCounts, openTemplate, validatePublications, user, profileData, lang]
  );

  const readOnlyCounts = useMemo(
    () => computeCounts(initial ?? {}),
    [computeCounts, initial]
  );
  const readOnlyTotal = useMemo(
    () => readOnlyCounts.reduce((acc, item) => acc + item.count, 0),
    [readOnlyCounts]
  );

  return (
    <div className="space-y-4">
      <FormTaskDetails
        node={node}
        initial={initial}
        canEdit={canEdit}
        {...rest}
        renderActions={renderActions}
        validationErrors={validationErrors}
        persistKey={`form_draft_${node.id}`}
        onValuesChange={() => setValidationErrors({})}
      />

      {canEdit === false && (
        <div className="space-y-2 rounded-md border bg-muted/30 p-4 text-sm">
          <div className="font-medium">
            {T("forms.collection_summary", "Записи по разделам")}
          </div>
          <div className="grid gap-1">
            {readOnlyCounts.map(({ key, label, count }) => (
              <div
                key={key}
                className="flex items-center justify-between gap-2"
              >
                <span className="text-muted-foreground">{label}</span>
                <span className="font-semibold">{count}</span>
              </div>
            ))}
          </div>
          <div className="text-xs text-muted-foreground">
            {T("forms.collection_total", "Всего записей")}: {readOnlyTotal}
          </div>
        </div>
      )}
    </div>
  );
}

export default S1PublicationsDetails;
