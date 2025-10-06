// components/node-details/variants/UploadTaskDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { NodeVM, UploadDef, t as tLabel } from "@/lib/playbook";
import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { AssetsDownloads } from "../AssetsDownloads";

type FileState = { def: UploadDef; file?: File | null };

export function UploadTaskDetails({
  node,
  onSubmit,
  canEdit = true,
  existing,
}: {
  node: NodeVM;
  onSubmit?: (payload: { files: Record<string, File | null> }) => void;
  canEdit?: boolean;
  existing?: Map<
    string,
    Array<{
      version_id: string;
      filename: string;
      download_url: string;
      size_bytes: number;
      is_active: boolean;
      attached_at?: string;
    }>
  >;
}) {
  const defs = node.requirements?.uploads ?? [];
  const [files, setFiles] = useState<Record<string, File | null>>({});
  const inputs = useRef<Record<string, HTMLInputElement | null>>({});
  const { t: T } = useTranslation("common");

  function pick(key: string, file: File | null) {
    setFiles((prev) => ({ ...prev, [key]: file ?? null }));
  }

  return (
    <Card className="p-4 space-y-4">
      {node.requirements?.notes && (
        <p className="text-sm text-muted-foreground">
          {node.requirements.notes}
        </p>
      )}

      <div className="space-y-3">
        {defs.map((u) => {
          const accept = u.accept ?? u.mime?.join(",") ?? undefined;
          const labelText = tLabel(u.label as any, u.key);
          const attachments = existing?.get(u.key) ?? [];
          return (
            <div
              key={u.key}
              className="flex items-center justify-between gap-4 rounded-md border p-3"
            >
              <div className="min-w-0">
                <div className="truncate text-sm font-medium">
                  {labelText}{" "}
                  {u.required ? (
                    <span className="text-destructive">*</span>
                  ) : null}
                </div>
                {attachments.length > 0 && (
                  <div className="mt-1 space-y-1 text-xs">
                    {attachments.map((att) => (
                      <div
                        key={att.version_id}
                        className="flex items-center gap-1"
                      >
                        <a
                          className="underline"
                          href={att.download_url}
                          target="_blank"
                          rel="noreferrer"
                        >
                          {att.filename}
                        </a>
                        {!att.is_active && (
                          <span className="text-muted-foreground">(old)</span>
                        )}
                      </div>
                    ))}
                  </div>
                )}
                {u.mime?.length ? (
                  <div className="text-xs text-muted-foreground">
                    {T("upload.allowed_types")}: {u.mime.join(", ")}
                  </div>
                ) : null}
              </div>
              <div className="flex items-center gap-2">
                <input
                  ref={(el) => (inputs.current[u.key] = el)}
                  type="file"
                  className="hidden"
                  accept={accept}
                  onChange={(e) => pick(u.key, e.target.files?.[0] ?? null)}
                  disabled={!canEdit}
                />
                <Button
                  variant="secondary"
                  onClick={() => inputs.current[u.key]?.click()}
                  disabled={!canEdit}
                >
                  {T("upload.select_file")}
                </Button>
                <div className="w-40 truncate text-xs text-muted-foreground">
                  {files[u.key]?.name ?? T("upload.no_file")}
                </div>
              </div>
            </div>
          );
        })}
      </div>

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

      {/* Templates / Downloads (if any) */}
      <AssetsDownloads node={node} />

      {canEdit && (
        <Button onClick={() => onSubmit?.({ files })}>
          {T("upload.submit")}
        </Button>
      )}
    </Card>
  );
}
