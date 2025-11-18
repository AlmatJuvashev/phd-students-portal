import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  NodeSubmissionDTO,
  attachNodeUpload,
  presignNodeUpload,
} from "@/api/journey";
import { API_URL } from "@/api/client";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Download, Loader2, Upload } from "lucide-react";

const statusStyles: Record<string, string> = {
  submitted: "bg-amber-50 text-amber-700 border border-amber-200",
  approved: "bg-emerald-50 text-emerald-700 border border-emerald-200",
  rejected: "bg-rose-50 text-rose-700 border border-rose-200",
};

function formatBytes(bytes?: number) {
  if (!bytes && bytes !== 0) return "";
  const thresh = 1024;
  if (Math.abs(bytes) < thresh) {
    return `${bytes} B`;
  }
  const units = ["KB", "MB", "GB", "TB"];
  let u = -1;
  let value = bytes;
  do {
    value /= thresh;
    ++u;
  } while (Math.abs(value) >= thresh && u < units.length - 1);
  return `${value.toFixed(1)} ${units[u]}`;
}

function formatDate(value?: string) {
  if (!value) return "";
  const d = new Date(value);
  return d.toLocaleDateString();
}

type Props = {
  nodeId: string;
  slots?: NodeSubmissionDTO["slots"];
  canEdit?: boolean;
  onRefresh?: () => void;
};

export function NodeAttachmentsSection({
  nodeId,
  slots = [],
  canEdit,
  onRefresh,
}: Props) {
  const { t, i18n } = useTranslation("common");
  const [uploadingSlot, setUploadingSlot] = useState<string | null>(null);
  const [message, setMessage] = useState<{
    text: string;
    tone: "error" | "success";
  } | null>(null);
  const fileInputs = useRef<Record<string, HTMLInputElement | null>>({});
  const [slotMeta, setSlotMeta] = useState<
    Record<string, { label?: string; required?: boolean }>
  >({});

  if (!slots || slots.length === 0) {
    return null;
  }

  // Load slot labels from playbook for nicer display (fallback to key)
  useEffect(() => {
    let mounted = true;
    import("@/playbooks/playbook.json")
      .then((mod: any) => {
        if (!mounted) return;
        const pb = (mod && (mod.default || mod)) as any;
        const worlds = (pb.worlds || pb.Worlds || []) as any[];
        const lang = i18n?.language || "en";
        const pick = (obj: any, key: string) =>
          obj?.[key] ||
          obj?.[key?.toUpperCase?.()] ||
          (key ? obj?.[key.charAt(0).toUpperCase() + key.slice(1)] : undefined);
        const findNode = () => {
          for (const w of worlds) {
            const nodesArr = (w.nodes || w.Nodes || []) as any[];
            for (const n of nodesArr) {
              const id = n.id || n.ID;
              if (id === nodeId) return n;
            }
          }
          return null;
        };
        const node = findNode();
        const uploads =
          node?.requirements?.uploads || node?.Requirements?.Uploads || [];
        const map: Record<string, { label?: string; required?: boolean }> =
          {} as any;
        for (const up of uploads) {
          const key = up.key || up.Key;
          const lbl = up.label || up.Label || {};
          const label = pick(lbl, lang.toLowerCase()) || pick(lbl, "en") || key;
          map[key] = { label, required: !!(up.required ?? up.Required) };
        }
        setSlotMeta(map);
      })
      .catch(() => setSlotMeta({}));
    return () => {
      mounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [nodeId, i18n?.language]);

  const handleFileSelected = async (
    slotKey: string,
    files: FileList | null
  ) => {
    const file = files?.[0];
    if (!file) return;
    setUploadingSlot(slotKey);
    setMessage(null);
    try {
      const contentType = file.type || "application/octet-stream";
      const presign = await presignNodeUpload(nodeId, {
        slot_key: slotKey,
        filename: file.name,
        content_type: contentType,
        size_bytes: file.size,
      });
      const headers = new Headers({ "Content-Type": contentType });
      if (presign?.required_headers) {
        Object.entries(presign.required_headers).forEach(([key, value]) => {
          headers.set(key, value as string);
        });
      }
      const uploadResp = await fetch(presign.upload_url, {
        method: "PUT",
        headers,
        body: file,
      });
      if (!uploadResp.ok) {
        throw new Error(`Upload failed (${uploadResp.status})`);
      }
      const etag =
        uploadResp.headers.get("ETag")?.replace(/"/g, "") ?? undefined;
      await attachNodeUpload(nodeId, {
        slot_key: slotKey,
        filename: file.name,
        object_key: presign.object_key,
        content_type: contentType,
        size_bytes: file.size,
        etag,
      });
      setMessage({
        text: t("uploads.success", { defaultValue: "File uploaded" }),
        tone: "success",
      });
      onRefresh?.();
    } catch (error: any) {
      setMessage({
        text:
          error?.message ||
          t("uploads.error", {
            defaultValue: "Failed to upload file. Try again.",
          }),
        tone: "error",
      });
    } finally {
      setUploadingSlot(null);
      const input = fileInputs.current[slotKey];
      if (input) input.value = "";
    }
  };

  const acceptFor = (mime: string[]) =>
    mime.length ? mime.join(",") : undefined;

  return (
    <section className="space-y-4">
      <div>
        <h3 className="text-base font-semibold">
          {t("uploads.title", { defaultValue: "Supporting documents" })}
        </h3>
        <p className="text-sm text-muted-foreground">
          {t("uploads.subtitle", {
            defaultValue: "Attach the required files for this node.",
          })}
        </p>
      </div>
      {message && (
        <div
          className={`rounded-md border p-3 text-sm ${
            message.tone === "error"
              ? "border-destructive/40 bg-destructive/5 text-destructive"
              : "border-emerald-200 bg-emerald-50 text-emerald-700"
          }`}
        >
          {message.text}
        </div>
      )}
      <div className="space-y-4">
        {slots.map((slot) => {
          const attachments = (slot.attachments || []).filter(
            (att) => att.is_active
          );
          return (
            <Card key={slot.key} className="p-4 space-y-3">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="space-y-1">
                  <p className="text-sm font-medium">
                    {slotMeta[slot.key]?.label || slot.key}
                    {slot.required && (
                      <span className="text-red-500 ml-1">*</span>
                    )}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {slot.mime && slot.mime.length
                      ? t("uploads.allowed", {
                          defaultValue: "Allowed: {{mime}}",
                          mime: slot.mime.join(", "),
                        })
                      : t("uploads.any_format", { defaultValue: "Any format" })}
                  </p>
                </div>
                {canEdit && (
                  <div className="flex items-center gap-2">
                    <input
                      type="file"
                      className="hidden"
                      accept={acceptFor(slot.mime)}
                      ref={(el) => (fileInputs.current[slot.key] = el)}
                      onChange={(event) =>
                        handleFileSelected(slot.key, event.target.files)
                      }
                    />
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      disabled={uploadingSlot === slot.key}
                      onClick={() => fileInputs.current[slot.key]?.click()}
                      className="flex items-center gap-2"
                    >
                      {uploadingSlot === slot.key ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Upload className="h-4 w-4" />
                      )}
                      {attachments.length > 0
                        ? t("uploads.add_version", {
                            defaultValue: "Add new version",
                          })
                        : t("uploads.add", { defaultValue: "Upload file" })}
                    </Button>
                  </div>
                )}
              </div>
              <div className="space-y-4">
                {attachments.length === 0 ? (
                  <p className="text-sm text-muted-foreground">
                    {t("uploads.empty", {
                      defaultValue: "No files uploaded yet.",
                    })}
                  </p>
                ) : (
                  <>
                    {attachments.map((att, index) => {
                    const statusClass = att.status
                      ? statusStyles[att.status] || statusStyles.submitted
                      : null;
                    const statusLabel = att.status
                      ? t(`uploads.status.${att.status}`, {
                          defaultValue:
                            att.status === "approved"
                              ? "Approved"
                              : att.status === "rejected"
                              ? "Needs fixes"
                              : "Pending",
                        })
                      : null;
                    const hasReviewedDoc = att.reviewed_document?.version_id;
                    const versionNumber = attachments.length - index;
                    return (
                      <div
                        key={att.version_id}
                        className={`rounded-md border ${
                          hasReviewedDoc
                            ? "border-blue-200 bg-blue-50/30"
                            : "border-dashed"
                        } p-3 space-y-3`}
                      >
                        {/* Version header */}
                        <div className="flex items-center gap-2 pb-2 border-b border-gray-200">
                          <span className="text-xs font-semibold text-gray-500">
                            {t("uploads.version", {
                              defaultValue: "Version",
                            })}{" "}
                            {versionNumber}
                          </span>
                          <span className="text-xs text-gray-400">·</span>
                          <span className="text-xs text-gray-500">
                            {formatDate(att.attached_at)}
                          </span>
                        </div>
                        
                        {/* Student's original submission */}
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-1">
                              <p className="text-xs font-semibold text-blue-600">
                                {t("uploads.your_submission", {
                                  defaultValue: "📤 Your submission",
                                })}
                              </p>
                              {statusClass && statusLabel && (
                                <Badge
                                  variant="outline"
                                  className={statusClass}
                                >
                                  {statusLabel}
                                </Badge>
                              )}
                            </div>
                            <button
                              onClick={async () => {
                                try {
                                  const response = await fetch(
                                    `${API_URL}${att.download_url}`,
                                    {
                                      headers: {
                                        Authorization: `Bearer ${localStorage.getItem(
                                          "token"
                                        )}`,
                                      },
                                    }
                                  );
                                  if (response.redirected) {
                                    window.open(response.url, "_blank");
                                  } else {
                                    const blob = await response.blob();
                                    const url = URL.createObjectURL(blob);
                                    const a = document.createElement("a");
                                    a.href = url;
                                    a.download = att.filename;
                                    a.click();
                                    URL.revokeObjectURL(url);
                                  }
                                } catch (err) {
                                  console.error("Download failed:", err);
                                }
                              }}
                              className="text-sm font-medium text-blue-600 break-all text-left hover:underline hover:text-blue-700 transition-colors cursor-pointer inline-flex items-center gap-1.5"
                              title={t("uploads.click_to_download", {
                                defaultValue: "Click to download",
                              })}
                            >
                              <Download className="h-3.5 w-3.5 flex-shrink-0" />
                              {att.filename}
                            </button>
                            <p className="text-xs text-muted-foreground">
                              {formatBytes(att.size_bytes)}
                              {att.attached_at
                                ? ` · ${formatDate(att.attached_at)}`
                                : ""}
                            </p>
                            {att.review_note && (
                              <div className="mt-2 p-2 bg-amber-50 border border-amber-200 rounded text-xs text-amber-900">
                                <p className="font-semibold mb-1">
                                  💬{" "}
                                  {t("uploads.advisor_note", {
                                    defaultValue: "Advisor's feedback:",
                                  })}
                                </p>
                                <p>{att.review_note}</p>
                              </div>
                            )}
                          </div>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={async () => {
                              try {
                                const response = await fetch(
                                  `${API_URL}${att.download_url}`,
                                  {
                                    headers: {
                                      Authorization: `Bearer ${localStorage.getItem(
                                        "token"
                                      )}`,
                                    },
                                  }
                                );
                                if (response.redirected) {
                                  window.open(response.url, "_blank");
                                } else {
                                  const blob = await response.blob();
                                  const url = URL.createObjectURL(blob);
                                  const a = document.createElement("a");
                                  a.href = url;
                                  a.download = att.filename;
                                  a.click();
                                  URL.revokeObjectURL(url);
                                }
                              } catch (err) {
                                console.error("Download failed:", err);
                              }
                            }}
                            aria-label={t("uploads.download", {
                              defaultValue: "Download",
                            })}
                          >
                            <Download className="h-4 w-4" />
                          </Button>
                        </div>

                        {/* Advisor's reviewed document (if exists) */}
                        {hasReviewedDoc && att.reviewed_document && (
                          <div className="flex flex-wrap items-start justify-between gap-3 pt-3 border-t border-blue-200">
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center gap-2 mb-1">
                                <p className="text-xs font-semibold text-emerald-600">
                                  {t("uploads.advisor_reviewed", {
                                    defaultValue: "📥 Advisor reviewed file",
                                  })}
                                </p>
                              </div>
                              <button
                                onClick={async () => {
                                  try {
                                    const response = await fetch(
                                      `${API_URL}${att.reviewed_document!.download_url}`,
                                      {
                                        headers: {
                                          Authorization: `Bearer ${localStorage.getItem(
                                            "token"
                                          )}`,
                                        },
                                      }
                                    );
                                    if (response.redirected) {
                                      window.open(response.url, "_blank");
                                    } else {
                                      const blob = await response.blob();
                                      const url = URL.createObjectURL(blob);
                                      const a = document.createElement("a");
                                      a.href = url;
                                      a.download =
                                        att.reviewed_document!.filename ||
                                        `Reviewed_${att.filename}`;
                                      a.click();
                                      URL.revokeObjectURL(url);
                                    }
                                  } catch (err) {
                                    console.error("Download failed:", err);
                                  }
                                }}
                                className="text-sm font-medium text-emerald-600 break-all text-left hover:underline hover:text-emerald-700 transition-colors cursor-pointer inline-flex items-center gap-1.5"
                                title={t("uploads.click_to_download", {
                                  defaultValue: "Click to download",
                                })}
                              >
                                <Download className="h-3.5 w-3.5 flex-shrink-0" />
                                {att.reviewed_document.filename ||
                                  `Reviewed_${att.filename}`}
                              </button>
                              <p className="text-xs text-muted-foreground">
                                {formatBytes(att.reviewed_document.size_bytes)}
                                {att.reviewed_document.reviewed_at
                                  ? ` · ${formatDate(
                                      att.reviewed_document.reviewed_at
                                    )}`
                                  : ""}
                                {att.reviewed_document.reviewed_by
                                  ? ` · ${t("uploads.reviewed_by_prefix", {
                                      defaultValue: "by",
                                    })} ${att.reviewed_document.reviewed_by}`
                                  : ""}
                              </p>
                            </div>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={async () => {
                                try {
                                  const response = await fetch(
                                    `${API_URL}${att.reviewed_document!.download_url}`,
                                    {
                                      headers: {
                                        Authorization: `Bearer ${localStorage.getItem(
                                          "token"
                                        )}`,
                                      },
                                    }
                                  );
                                  if (response.redirected) {
                                    window.open(response.url, "_blank");
                                  } else {
                                    const blob = await response.blob();
                                    const url = URL.createObjectURL(blob);
                                    const a = document.createElement("a");
                                    a.href = url;
                                    a.download =
                                      att.reviewed_document!.filename ||
                                      `Reviewed_${att.filename}`;
                                    a.click();
                                    URL.revokeObjectURL(url);
                                  }
                                } catch (err) {
                                  console.error("Download failed:", err);
                                }
                              }}
                              aria-label={t("uploads.download_reviewed", {
                                defaultValue: "Download reviewed file",
                              })}
                            >
                              <Download className="h-4 w-4 text-emerald-600" />
                            </Button>
                          </div>
                        )}
                      </div>
                    );
                    })}
                  </>
                )}
              </div>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
