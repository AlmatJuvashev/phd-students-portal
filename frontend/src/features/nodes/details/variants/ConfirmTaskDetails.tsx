import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { allAssets, getAssetUrl, type PublicAsset } from "@/lib/assets";
import { t, safeText } from "@/lib/playbook";
import type { NodeVM } from "@/lib/playbook";
import i18n from "i18next";
import type { NodeSubmissionDTO } from "@/api/journey";
import { NodeAttachmentsSection } from "../NodeAttachmentsSection";
import { useAuth } from "@/contexts/AuthContext";
import { useProfileSnapshot } from "@/features/profile/useProfileSnapshot";
import {
  buildTemplateData,
  generateStudentTemplateDoc,
  supportsStudentDocTemplate,
} from "@/features/docgen/student-template";
import type { StudentTemplateData } from "@/features/docgen/student-template";
import { Loader2, Clock } from "lucide-react";

type ConfirmTaskDetailsProps = {
  node: NodeVM | any;
  onComplete?: () => void;
  onReset?: () => void;
  slots?: NodeSubmissionDTO["slots"];
  canEdit?: boolean;
  onRefresh?: () => void;
};

const ConfirmTaskDetails: React.FC<ConfirmTaskDetailsProps> = ({
  node,
  onComplete: _onComplete,
  onReset: _onReset,
  slots,
  canEdit,
  onRefresh,
}) => {
  void _onComplete;
  void _onReset;
  // Localized question (string or i18n map)
  const question: string | undefined = safeText(
    node?.screen?.question as any,
    ""
  );

  // Primary button (index 0) contains instructions
  const primaryBtn = Array.isArray(node?.screen?.buttons)
    ? node.screen.buttons[0]
    : undefined;
  // Normalize instruction text: could be array of strings or a locale->string[] map
  const instructionsRaw = primaryBtn?.instructions?.text as
    | string[]
    | Record<string, string[]>
    | undefined;
  const currentLang = (i18n?.language as "ru" | "kz" | "en") || "ru";
  const instructions: string[] = Array.isArray(instructionsRaw)
    ? instructionsRaw
    : Array.isArray((instructionsRaw as any)?.[currentLang])
    ? (instructionsRaw as any)[currentLang]
    : [];
  const { user } = useAuth();
  const locale = (currentLang as "ru" | "kz" | "en") || "ru";
  const { data: profileData, isLoading: profileLoading } = useProfileSnapshot(
    true
  );
  const templateData = React.useMemo<StudentTemplateData | null>(() => {
    return buildTemplateData(user, profileData as any, locale);
  }, [user, profileData, locale]);
  const [downloadingId, setDownloadingId] = React.useState<string | null>(null);
  const profileRequiredMsg = React.useMemo(
    () =>
      i18n.t("templates.profile_required", {
        defaultValue:
          "Please fill your doctoral profile form to auto-fill templates.",
      }) as string,
    [currentLang]
  );
  const errorFallbackMsg = React.useMemo(
    () =>
      i18n.t("common.error", {
        defaultValue: "Error",
      }) as string,
    [currentLang]
  );
  const downloadSingle = (primaryBtn?.instructions?.download || undefined) as
    | { label?: any; asset_id?: string; asset_path?: string }
    | undefined;
  const downloadsRaw = primaryBtn?.instructions?.downloads as
    | Array<{ label?: any; asset_id?: string; asset_path?: string }>
    | undefined;
  const downloadItems = React.useMemo(() => {
    const items =
      Array.isArray(downloadsRaw) && downloadsRaw.length > 0
        ? downloadsRaw
        : downloadSingle
        ? [downloadSingle]
        : [];

    if (items.length <= 1) return items;

    const matchByLabel = items.find((item) => {
      const lbl = item.label;
      return lbl && typeof lbl === "object" && lbl[currentLang];
    });

    if (matchByLabel) return [matchByLabel];

    const matchByAsset = items.find((item) => {
      const id = (item.asset_id || "").toLowerCase();
      if (!id) return false;
      const lang = currentLang.toLowerCase();
      return (
        id.includes(`_${lang}`) ||
        id.includes(`-${lang}`) ||
        id.endsWith(`${lang}.docx`) ||
        id.endsWith(`${lang}.pdf`)
      );
    });

    if (matchByAsset) return [matchByAsset];

    const english = items.find((item) => {
      const lbl = item.label;
      const id = (item.asset_id || "").toLowerCase();
      return (
        (lbl && typeof lbl === "object" && lbl.en) ||
        id.includes("_en") ||
        id.includes("-en")
      );
    });

    return [matchByLabel || matchByAsset || english || items[0]].filter(
      Boolean
    ) as Array<{
      label?: any;
      asset_id?: string;
      asset_path?: string;
    }>;
  }, [downloadsRaw, downloadSingle, currentLang]);
  React.useEffect(() => {
    if (typeof window !== "undefined") {
      console.log("[templated-download] init", {
        source: "confirm-task",
        userRole: user?.role,
        locale,
        downloadItems: downloadItems.map((d) => d.asset_id || d.asset_path),
        templatable: downloadItems
          .map((d) => d.asset_id)
          .filter(Boolean)
          .map((id) => allAssets().find((a) => a.id === id))
          .filter((a): a is PublicAsset => !!a && supportsStudentDocTemplate(a))
          .map((a) => a.id),
      });
    }
  }, [downloadItems, locale, user?.role]);
  const handleTemplatedDownload = React.useCallback(
    async ({
      asset,
      href,
      label,
      key,
    }: {
      asset: PublicAsset;
      href: string;
      label: string;
      key: string;
    }) => {
      if (typeof window !== "undefined") {
        console.log("[templated-download] click", {
          source: "confirm-task",
          assetId: asset.id,
          locale,
          hasTemplateData: !!templateData,
          href,
        });
      }
      if (!href) return;
      if (!templateData) {
        console.log("[templated-download] missing template data");
        window.alert(profileRequiredMsg);
        window.open(href, "_blank", "noopener,noreferrer");
        return;
      }
      try {
        setDownloadingId(key);
        await generateStudentTemplateDoc({
          asset,
          data: templateData,
          locale,
          // Don't pass fileLabel - let it use asset.title for more descriptive filename
        });
      } catch (err) {
        console.error(err);
        window.alert(err instanceof Error ? err.message : errorFallbackMsg);
        window.open(href, "_blank", "noopener,noreferrer");
      } finally {
        setDownloadingId(null);
      }
    },
    [templateData, locale, profileRequiredMsg, errorFallbackMsg]
  );
  const accordionLabel = t(
    primaryBtn?.label as any,
    t(
      {
        ru: "Инструкция по прохождению",
        kz: "Өту бойынша нұсқаулық",
        en: "How to complete",
      },
      "Инструкция по прохождению"
    )
  );

  const fallbackSlots = React.useMemo(() => {
    const uploads =
      ((node?.requirements as any)?.uploads as Array<any> | undefined) || [];
    if (!Array.isArray(uploads) || uploads.length === 0) return [];
    return uploads.map((up) => ({
      key: up?.key,
      required: !!up?.required,
      multiplicity: up?.multiplicity || "single",
      mime: Array.isArray(up?.mime) ? up.mime : [],
      attachments: [],
    }));
  }, [node]);

  const attachmentSlots =
    Array.isArray(slots) && slots.length > 0 ? slots : fallbackSlots;

  return (
    <Card className="bg-gradient-to-br from-card to-card/50">
      <CardContent className="space-y-6 p-4 sm:p-6">
        {question && (
          <p className="text-lg sm:text-xl font-semibold text-foreground leading-relaxed">
            {question}
          </p>
        )}

        {/* Submitted State Banner */}
        {node?.state === 'submitted' && (
            <div className="bg-orange-50 dark:bg-orange-950/20 border border-orange-200 dark:border-orange-800 rounded-xl p-4 flex items-center gap-3">
                <div className="bg-orange-100 dark:bg-orange-900 p-2 rounded-full text-orange-600 dark:text-orange-400">
                    <Clock className="w-5 h-5" />
                </div>
                <div>
                    <p className="font-semibold text-orange-800 dark:text-orange-300">
                        {t({ru: "Отправлено на проверку", kz: "Тексеруге жіберілді", en: "Submitted for review"}, "Submitted")}
                    </p>
                    <p className="text-xs text-orange-700 dark:text-orange-400">
                        {t({ru: "Ваш документ получен и ожидает проверки.", kz: "Сіздің құжатыңыз қабылданды және тексеруді күтуде.", en: "Your document has been received and is awaiting review."}, "Waiting for review.")}
                    </p>
                </div>
            </div>
        )}

        {/* Inline guidance (no collapsible) */}
        <div className="space-y-3 p-5 sm:p-6 rounded-xl bg-muted/30 border-l-4 border-primary/40">
          <div className="text-sm font-semibold text-primary flex items-center gap-2">
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {accordionLabel}
          </div>
          {Array.isArray(instructions) && instructions.length > 0 && (
            <ul className="list-disc pl-5 space-y-2 text-sm text-muted-foreground">
              {instructions.map((line: string, idx: number) => (
                <li key={idx} className="leading-relaxed">
                  {line}
                </li>
              ))}
            </ul>
          )}
          {downloadItems.length > 0 && (
            <div className="mt-3 flex flex-col gap-2">
              {downloadItems.map((item, idx) => {
            const asset = item.asset_id
              ? allAssets().find((a) => a.id === item.asset_id)
              : undefined;
            const resolved = item.asset_id
              ? getAssetUrl(item.asset_id)
                  : undefined;
                const href =
                  resolved && resolved !== "#" ? resolved : item.asset_path;
                if (!href) return null;
                const label = safeText(
                  item.label as any,
                  t(
                    {
                      ru: "Скачать пример письма",
                      kz: "Хаттың үлгісін жүктеу",
                      en: "Download sample letter",
                    },
                    "Скачать пример письма"
                  )
                );
                const key = item.asset_id || item.asset_path || `${idx}`;
                const isTemplated =
                  !!asset &&
                  supportsStudentDocTemplate(asset) &&
                  (!!templateData || user?.role === "student");
                const isBusy = downloadingId === key;
                const icon = isBusy ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                );

                if (!isTemplated || !asset) {
                  return (
                    <Button
                      key={key}
                      asChild
                      variant="secondary"
                      size="sm"
                      className="gap-2"
                    >
                      <a
                        href={href}
                        download
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {icon}
                        {label}
                      </a>
                    </Button>
                  );
                }

                return (
                  <Button
                    key={key}
                    variant="secondary"
                    size="sm"
                    className="gap-2"
                    disabled={isBusy || profileLoading}
                    onClick={() =>
                      handleTemplatedDownload({
                        asset,
                        href,
                        label,
                        key,
                      })
                    }
                  >
                    {icon}
                    {label}
                  </Button>
                );
              })}
            </div>
          )}
        </div>

        {attachmentSlots.length > 0 && (
          <div className="pt-2">
            <div className="text-sm font-semibold text-foreground mb-3">
              {t(
                {
                  ru: "Поддерживающие документы",
                  kz: "Қолдаушы құжаттар",
                  en: "Supporting Documents",
                },
                "Поддерживающие документы"
              )}
            </div>
            <NodeAttachmentsSection
              nodeId={node?.node_id || node?.id}
              slots={attachmentSlots}
              canEdit={canEdit ?? node?.state !== "done"}
              onRefresh={onRefresh}
            />
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ConfirmTaskDetails;
