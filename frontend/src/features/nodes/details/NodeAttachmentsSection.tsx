import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { NodeSubmissionDTO, attachNodeUpload, presignNodeUpload } from "@/api/journey";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Download, Loader2, Upload } from "lucide-react";

const statusStyles: Record<string, { label: string; className: string }> = {
  submitted: {
    label: "Pending",
    className: "bg-amber-50 text-amber-700 border border-amber-200",
  },
  approved: {
    label: "Approved",
    className: "bg-emerald-50 text-emerald-700 border border-emerald-200",
  },
  rejected: {
    label: "Needs fixes",
    className: "bg-rose-50 text-rose-700 border border-rose-200",
  },
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

export function NodeAttachmentsSection({ nodeId, slots = [], canEdit, onRefresh }: Props) {
  const { t } = useTranslation("common");
  const [uploadingSlot, setUploadingSlot] = useState<string | null>(null);
  const [message, setMessage] = useState<{ text: string; tone: "error" | "success" } | null>(null);
  const fileInputs = useRef<Record<string, HTMLInputElement | null>>({});

  if (!slots || slots.length === 0) {
    return null;
  }

  const handleFileSelected = async (slotKey: string, files: FileList | null) => {
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
      const etag = uploadResp.headers.get("ETag")?.replace(/"/g, "") ?? undefined;
      await attachNodeUpload(nodeId, {
        slot_key: slotKey,
        filename: file.name,
        object_key: presign.object_key,
        content_type: contentType,
        size_bytes: file.size,
        etag,
      });
      setMessage({ text: t("uploads.success", { defaultValue: "File uploaded" }), tone: "success" });
      onRefresh?.();
    } catch (error: any) {
      setMessage({
        text:
          error?.message ||
          t("uploads.error", { defaultValue: "Failed to upload file. Try again." }),
        tone: "error",
      });
    } finally {
      setUploadingSlot(null);
      const input = fileInputs.current[slotKey];
      if (input) input.value = "";
    }
  };

  const acceptFor = (mime: string[]) => (mime.length ? mime.join(",") : undefined);

  return (
    <section className="space-y-4">
      <div>
        <h3 className="text-base font-semibold">{t("uploads.title", { defaultValue: "Supporting documents" })}</h3>
        <p className="text-sm text-muted-foreground">
          {t("uploads.subtitle", { defaultValue: "Attach the required files for this node." })}
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
          const attachments = (slot.attachments || []).filter((att) => att.is_active);
          return (
            <Card key={slot.key} className="p-4 space-y-3">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="space-y-1">
                  <p className="text-sm font-medium">
                    {slot.key}
                    {slot.required && <span className="text-red-500 ml-1">*</span>}
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
                      onChange={(event) => handleFileSelected(slot.key, event.target.files)}
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
                        ? t("uploads.replace", { defaultValue: "Replace file" })
                        : t("uploads.add", { defaultValue: "Upload file" })}
                    </Button>
                  </div>
                )}
              </div>
              <div className="space-y-2">
                {attachments.length === 0 ? (
                  <p className="text-sm text-muted-foreground">
                    {t("uploads.empty", { defaultValue: "No files uploaded yet." })}
                  </p>
                ) : (
                  attachments.map((att) => {
                    const status = att.status ? statusStyles[att.status] : null;
                    return (
                      <div
                        key={att.version_id}
                        className="flex flex-wrap items-center justify-between gap-3 rounded-md border border-dashed px-3 py-2"
                      >
                        <div>
                          <p className="text-sm font-medium text-foreground break-all">{att.filename}</p>
                          <p className="text-xs text-muted-foreground">
                            {formatBytes(att.size_bytes)}
                            {att.attached_at ? ` · ${formatDate(att.attached_at)}` : ""}
                          </p>
                          {att.review_note && (
                            <p className="text-xs text-amber-700 mt-1">
                              {t("uploads.note", { defaultValue: "Reviewer note:" })} {att.review_note}
                            </p>
                          )}
                        </div>
                        <div className="flex items-center gap-2">
                          {status && (
                            <Badge variant="outline" className={status.className}>
                              {status.label}
                            </Badge>
                          )}
                          <Button variant="ghost" size="icon" asChild>
                            <a
                              href={att.download_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              aria-label={t("uploads.download", { defaultValue: "Download" })}
                            >
                              <Download className="h-4 w-4" />
                            </a>
                          </Button>
                        </div>
                      </div>
                    );
                  })
                )}
              </div>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
