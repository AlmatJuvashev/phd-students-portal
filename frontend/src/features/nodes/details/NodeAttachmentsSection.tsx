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
import { Download, Loader2, Upload, FileText, MessageSquare, Eye } from "lucide-react";
import { DocumentPreviewDrawer } from "@/components/ui/DocumentPreviewDrawer";

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

function formatDateTime(value?: string) {
  if (!value) return "";
  const d = new Date(value);
  return `${d.toLocaleDateString()} ${d.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  })}`;
}

type TimelineEvent = {
  id: string;
  type: "student_upload" | "advisor_review";
  timestamp: string;
  filename: string;
  downloadUrl: string;
  sizeBytes?: number;
  status?: string;
  reviewNote?: string;
  reviewedBy?: string;
  versionNumber?: number;
};

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
  const [previewState, setPreviewState] = useState<{
    isOpen: boolean;
    fileUrl: string;
    filename: string;
    contentType?: string;
  }>({
    isOpen: false,
    fileUrl: "",
    filename: "",
    contentType: undefined,
  });

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

  const handleDownload = async (downloadUrl: string, filename: string) => {
    try {
      const response = await fetch(`${API_URL}${downloadUrl}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (response.redirected) {
        window.open(response.url, "_blank");
      } else {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        a.click();
        URL.revokeObjectURL(url);
      }
    } catch (err) {
      console.error("Download failed:", err);
    }
  };

  const handlePreview = async (
    downloadUrl: string,
    filename: string,
    contentType?: string
  ) => {
    try {
      const response = await fetch(`${API_URL}${downloadUrl}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });

      let previewUrl: string;
      if (response.redirected) {
        // S3 presigned URL
        previewUrl = response.url;
      } else {
        // Create blob URL for local preview
        const blob = await response.blob();
        previewUrl = URL.createObjectURL(blob);
      }

      setPreviewState({
        isOpen: true,
        fileUrl: previewUrl,
        filename,
        contentType,
      });
    } catch (err) {
      console.error("Preview failed:", err);
    }
  };

  const closePreview = () => {
    // Cleanup blob URL if it was created
    if (
      previewState.fileUrl &&
      previewState.fileUrl.startsWith("blob:")
    ) {
      URL.revokeObjectURL(previewState.fileUrl);
    }
    setPreviewState({
      isOpen: false,
      fileUrl: "",
      filename: "",
      contentType: undefined,
    });
  };

  const buildTimeline = (attachments: any[]): TimelineEvent[] => {
    const events: TimelineEvent[] = [];
    const sortedAttachments = [...attachments].reverse(); // Oldest first

    sortedAttachments.forEach((att, index) => {
      const versionNumber = index + 1;

      // Student upload event
      events.push({
        id: `student-${att.version_id}`,
        type: "student_upload",
        timestamp: att.attached_at,
        filename: att.filename,
        downloadUrl: att.download_url,
        sizeBytes: att.size_bytes,
        status: att.status,
        reviewNote: att.review_note,
        versionNumber,
      });

      // Advisor review event (if exists)
      if (att.reviewed_document?.version_id) {
        events.push({
          id: `advisor-${att.reviewed_document.version_id}`,
          type: "advisor_review",
          timestamp:
            att.reviewed_document.reviewed_at || att.reviewed_document.created_at,
          filename:
            att.reviewed_document.filename || `Reviewed_${att.filename}`,
          downloadUrl: att.reviewed_document.download_url,
          sizeBytes: att.reviewed_document.size_bytes,
          reviewedBy: att.reviewed_document.reviewed_by,
        });
      }
    });

    return events;
  };

  return (
    <section className="space-y-4">
      <DocumentPreviewDrawer
        isOpen={previewState.isOpen}
        onClose={closePreview}
        fileUrl={previewState.fileUrl}
        filename={previewState.filename}
        contentType={previewState.contentType}
      />
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
              {/* Timeline View */}
              <div className="space-y-0">
                {attachments.length === 0 ? (
                  <p className="text-sm text-muted-foreground">
                    {t("uploads.empty", {
                      defaultValue: "No files uploaded yet.",
                    })}
                  </p>
                ) : (
                  <div className="relative">
                    {/* Timeline vertical line */}
                    <div className="absolute left-6 top-8 bottom-8 w-0.5 bg-gradient-to-b from-blue-200 via-emerald-200 to-blue-200" />

                    {buildTimeline(attachments).map((event, index) => {
                      const isStudent = event.type === "student_upload";
                      const statusClass = event.status
                        ? statusStyles[event.status] || statusStyles.submitted
                        : null;
                      const statusLabel = event.status
                        ? t(`uploads.status.${event.status}`, {
                            defaultValue:
                              event.status === "approved"
                                ? "Approved"
                                : event.status === "rejected"
                                ? "Needs fixes"
                                : "Pending",
                          })
                        : null;

                      return (
                        <div
                          key={event.id}
                          className={`relative flex gap-4 pb-6 ${
                            index === 0 ? "pt-0" : ""
                          }`}
                        >
                          {/* Timeline dot */}
                          <div
                            className={`flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center z-10 ${
                              isStudent
                                ? "bg-blue-100 border-2 border-blue-400"
                                : "bg-emerald-100 border-2 border-emerald-400"
                            }`}
                          >
                            {isStudent ? (
                              <FileText className="h-5 w-5 text-blue-600" />
                            ) : (
                              <MessageSquare className="h-5 w-5 text-emerald-600" />
                            )}
                          </div>

                          {/* Event content */}
                          <div
                            className={`flex-1 ${
                              isStudent ? "ml-0" : "ml-0"
                            }`}
                          >
                            <div
                              className={`rounded-lg border p-4 ${
                                isStudent
                                  ? "bg-blue-50/50 border-blue-200"
                                  : "bg-emerald-50/50 border-emerald-200"
                              }`}
                            >
                              {/* Header */}
                              <div className="flex items-start justify-between gap-2 mb-2">
                                <div className="flex-1">
                                  <div className="flex items-center gap-2 flex-wrap">
                                    <span
                                      className={`text-sm font-semibold ${
                                        isStudent
                                          ? "text-blue-700"
                                          : "text-emerald-700"
                                      }`}
                                    >
                                      {isStudent
                                        ? event.versionNumber
                                          ? `${t("uploads.version", {
                                              defaultValue: "Version",
                                            })} ${event.versionNumber}`
                                          : t("uploads.your_submission", {
                                              defaultValue: "Your submission",
                                            })
                                        : t("uploads.advisor_reviewed", {
                                            defaultValue: "Advisor reviewed",
                                          })}
                                    </span>
                                    {statusClass && statusLabel && isStudent && (
                                      <Badge
                                        variant="outline"
                                        className={statusClass}
                                      >
                                        {statusLabel}
                                      </Badge>
                                    )}
                                  </div>
                                  <p className="text-xs text-muted-foreground mt-0.5">
                                    {formatDateTime(event.timestamp)}
                                    {event.reviewedBy &&
                                      ` · ${t("uploads.reviewed_by_prefix", {
                                        defaultValue: "by",
                                      })} ${event.reviewedBy}`}
                                  </p>
                                </div>
                              </div>

                              {/* File info */}
                              <div className="flex items-center justify-between gap-3">
                                <button
                                  onClick={() =>
                                    handleDownload(
                                      event.downloadUrl,
                                      event.filename
                                    )
                                  }
                                  className={`text-sm font-medium break-all text-left hover:underline transition-colors cursor-pointer inline-flex items-center gap-1.5 ${
                                    isStudent
                                      ? "text-blue-600 hover:text-blue-700"
                                      : "text-emerald-600 hover:text-emerald-700"
                                  }`}
                                  title={t("uploads.click_to_download", {
                                    defaultValue: "Click to download",
                                  })}
                                >
                                  <Download className="h-3.5 w-3.5 flex-shrink-0" />
                                  {event.filename}
                                </button>
                                <div className="flex items-center gap-1">
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() =>
                                      handlePreview(
                                        event.downloadUrl,
                                        event.filename
                                      )
                                    }
                                    aria-label={t("uploads.preview", {
                                      defaultValue: "Preview",
                                    })}
                                    title={t("uploads.preview", {
                                      defaultValue: "Preview",
                                    })}
                                  >
                                    <Eye
                                      className={`h-4 w-4 ${
                                        isStudent
                                          ? "text-blue-600"
                                          : "text-emerald-600"
                                      }`}
                                    />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() =>
                                      handleDownload(
                                        event.downloadUrl,
                                        event.filename
                                      )
                                    }
                                    aria-label={t("uploads.download", {
                                      defaultValue: "Download",
                                    })}
                                  >
                                    <Download
                                      className={`h-4 w-4 ${
                                        isStudent
                                          ? "text-blue-600"
                                          : "text-emerald-600"
                                      }`}
                                    />
                                  </Button>
                                </div>
                              </div>
                              <p className="text-xs text-muted-foreground mt-1">
                                {formatBytes(event.sizeBytes)}
                              </p>

                              {/* Review note */}
                              {event.reviewNote && (
                                <div className="mt-3 p-2.5 bg-amber-50 border border-amber-200 rounded-md">
                                  <p className="text-xs font-semibold text-amber-900 mb-1 flex items-center gap-1.5">
                                    <MessageSquare className="h-3.5 w-3.5" />
                                    {t("uploads.advisor_note", {
                                      defaultValue: "Advisor's feedback:",
                                    })}
                                  </p>
                                  <p className="text-xs text-amber-900">
                                    {event.reviewNote}
                                  </p>
                                </div>
                              )}
                            </div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
