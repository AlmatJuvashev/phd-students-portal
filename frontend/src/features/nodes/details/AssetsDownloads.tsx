// components/node-details/AssetsDownloads.tsx
import { PublicAsset } from "@/lib/assets";
import { NodeVM, t } from "@/lib/playbook";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useTemplatesForNode } from "@/features/nodes/useTemplatesForNode";
import { useAuth } from "@/contexts/AuthContext";
import { useProfileSnapshot } from "@/features/profile/useProfileSnapshot";
import {
  generateStudentTemplateDoc,
  supportsStudentDocTemplate,
  type StudentTemplateData,
} from "@/features/docgen/student-template";
import { Loader2 } from "lucide-react";

export function AssetsDownloads({ node }: { node: NodeVM }) {
  const { i18n, t: T } = useTranslation("common");
  const locale = (i18n.language as "ru" | "kz" | "en") || "ru";
  const { user } = useAuth();
  const assets: PublicAsset[] = useTemplatesForNode(node);
  const { data: profileData, isLoading: profileLoading } =
    useProfileSnapshot(user?.role === "student");
  const templateData = useMemo(() => {
    if (user?.role !== "student") {
      return null;
    }
    return buildTemplateData(user, profileData as any, locale);
  }, [user, profileData, locale]);
  const [downloadingId, setDownloadingId] = useState<string | null>(null);
  const { order, groups } = useMemo(() => {
    const grouped: Record<string, PublicAsset[]> = {};
    for (const asset of assets) {
      const id = asset.id.toLowerCase();
      const match = id.match(/(app\d+)/i);
      let key = match ? match[1].toLowerCase() : id;
      if (!match && id.includes("omid")) key = "omid";
      if (!match && !id.includes("omid")) {
        key = id.replace(/_(ru|kz|en)(_.+)?$/, "");
      }
      grouped[key] = grouped[key] || [];
      grouped[key].push(asset);
    }
    return {
      order: Object.keys(grouped).sort(),
      groups: grouped,
    };
  }, [assets]);

  if (!order.length) return null;
  return (
    <div className="space-y-4 sticky top-4">
      <Separator className="my-4" />
      <div className="space-y-3">
        <h3
          id="templates-section"
          className="font-bold text-base sm:text-lg flex items-center gap-2 text-foreground"
        >
          <svg
            className="w-5 h-5 text-primary"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
            />
          </svg>
          {T("templates.heading")}
        </h3>
        <div className="space-y-2">
          {order.map((g, idx) => {
            const items = groups[g];
            const preferred =
              items.find((a) => a.id.toLowerCase().includes(`_${locale}`)) ||
              items.find((a) => a.title?.[locale as any]) ||
              items[0];
            if (!preferred) return null;
            const label =
              preferred.title?.[locale as any] ||
              t(preferred.title, preferred.id);
            const isTemplated =
              user?.role === "student" &&
              supportsStudentDocTemplate(preferred);
            return (
              <div
                key={g}
                className="animate-in fade-in slide-in-from-right-2 duration-300"
                style={{ animationDelay: `${idx * 50}ms` }}
              >
                {isTemplated ? (
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full justify-start gap-2 hover:bg-primary/5 hover:border-primary/40 hover:shadow-sm transition-all duration-200 group"
                    disabled={
                      !!downloadingId || profileLoading || !templateData
                    }
                    onClick={() =>
                      handleTemplatedDownload({
                        asset: preferred,
                        href: `/${preferred.storage.key}`,
                        label,
                        locale,
                        templateData: templateData || undefined,
                        setDownloadingId,
                        profileAvailable: !!profileData,
                        showAlert: (msg) =>
                          window.alert(msg ?? "Unable to generate template"),
                        t: T,
                      })
                    }
                  >
                    {downloadingId === preferred.id ? (
                      <Loader2 className="w-4 h-4 animate-spin" />
                    ) : (
                      <svg
                        className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors"
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
                    )}
                    <span className="text-xs sm:text-sm truncate">{label}</span>
                  </Button>
                ) : (
                  <a
                    href={`/${preferred.storage.key}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    download
                    className="block"
                  >
                    <Button
                      variant="outline"
                      size="sm"
                      className="w-full justify-start gap-2 hover:bg-primary/5 hover:border-primary/40 hover:shadow-sm transition-all duration-200 group"
                    >
                      <svg
                        className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors"
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
                      <span className="text-xs sm:text-sm truncate">
                        {label}
                      </span>
                    </Button>
                  </a>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

function buildTemplateData(
  user: { full_name?: string; first_name?: string; last_name?: string } | null,
  profile: Record<string, any> | null | undefined,
  locale: "ru" | "kz" | "en"
): StudentTemplateData {
  const data = (profile?.form?.data ?? profile) as Record<string, any> | undefined;
  const advisorsValue = data?.advisors_full_names;
  const advisors = Array.isArray(advisorsValue)
    ? advisorsValue
    : typeof advisorsValue === "string"
      ? advisorsValue.split(/\r?\n/)
      : [];
  const fullName =
    data?.full_name ||
    user?.full_name ||
    [user?.first_name, user?.last_name].filter(Boolean).join(" ");
  const program = data?.program || "";
  const specialty = data?.specialty || program;
  const localeMap: Record<string, string> = {
    ru: "ru-RU",
    kz: "kk-KZ",
    en: "en-US",
  };
  const formatter = new Intl.DateTimeFormat(localeMap[locale] || "ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
  const submissionDate = formatter.format(new Date());
  return {
    student_full_name: fullName || "",
    student_program: program || "",
    student_specialty: specialty || "",
    student_supervisors: advisors
      .map((a: string) => (a || "").trim())
      .filter(Boolean)
      .join("\n"),
    submission_date: submissionDate,
  };
}

async function handleTemplatedDownload({
  asset,
  href,
  label,
  locale,
  templateData,
  setDownloadingId,
  profileAvailable,
  showAlert,
  t,
}: {
  asset: PublicAsset;
  href: string;
  label: string;
  locale: "ru" | "kz" | "en";
  templateData?: StudentTemplateData | null;
  setDownloadingId: (id: string | null) => void;
  profileAvailable: boolean;
  showAlert: (msg?: string) => void;
  t: (key: string, options?: Record<string, any>) => string;
}) {
  try {
    if (!templateData || !profileAvailable) {
      showAlert(
        t("templates.profile_required", {
          defaultValue:
            "Please fill your doctoral profile form to auto-fill templates.",
        })
      );
      window.open(href, "_blank", "noopener,noreferrer");
      return;
    }
    setDownloadingId(asset.id);
    await generateStudentTemplateDoc({
      asset,
      data: templateData,
      locale,
      fileLabel: label,
    });
  } catch (err) {
    console.error(err);
    showAlert(
      err instanceof Error
        ? err.message
        : t("common.error", { defaultValue: "Error" })
    );
    window.open(href, "_blank", "noopener,noreferrer");
  } finally {
    setDownloadingId(null);
  }
}
