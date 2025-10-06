// components/node-details/variants/UploadTaskDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { NodeVM, UploadDef } from "@/lib/playbook";
import { useRef, useState } from "react";
import { AssetsDownloads } from "../AssetsDownloads";

type FileState = { def: UploadDef; file?: File | null };

export function UploadTaskDetails({
  node,
  onSubmit,
  canEdit = true,
}: {
  node: NodeVM;
  onSubmit?: (payload: { files: Record<string, File | null> }) => void;
  canEdit?: boolean;
}) {
  const defs = node.requirements?.uploads ?? [];
  const [files, setFiles] = useState<Record<string, File | null>>({});
  const inputs = useRef<Record<string, HTMLInputElement | null>>({});

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
          const label = u.label?.ru ?? u.label?.en ?? u.key;
          return (
            <div
              key={u.key}
              className="flex items-center justify-between gap-4 rounded-md border p-3"
            >
              <div className="min-w-0">
                <div className="truncate text-sm font-medium">
                  {label}{" "}
                  {u.required ? (
                    <span className="text-destructive">*</span>
                  ) : null}
                </div>
                {u.mime?.length ? (
                  <div className="text-xs text-muted-foreground">
                    Допустимые типы: {u.mime.join(", ")}
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
                  Выбрать файл
                </Button>
                <div className="w-40 truncate text-xs text-muted-foreground">
                  {files[u.key]?.name ?? "Нет файла"}
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
            <div className="mb-2 font-medium">Автопроверки</div>
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
          Загрузить и Отправить
        </Button>
      )}
    </Card>
  );
}
