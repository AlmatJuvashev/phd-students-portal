import { useTranslation } from "react-i18next";
import { useFormContext } from "./FormProvider";
import { t } from "@/lib/playbook";

export interface FormStatusProps {
  className?: string;
}

/**
 * Displays form submission status and timestamp
 */
export function FormStatus({ className = "" }: FormStatusProps) {
  const { t: T } = useTranslation("common");
  const { readOnly, initial, node } = useFormContext();

  if (!readOnly) return null;

  const submittedAt: string | undefined = (initial as any)?.__submittedAt;

  return (
    <div
      className={`mt-3 text-sm text-muted-foreground whitespace-pre-line ${className}`}
    >
      {t(
        {
          ru: `Форма отправлена${
            submittedAt
              ? ` (дата: ${new Date(submittedAt).toLocaleDateString("ru-RU")})`
              : ""
          }.`,
          kz: `Форма жіберілді${
            submittedAt
              ? ` (күні: ${new Date(submittedAt).toLocaleDateString("kk-KZ")})`
              : ""
          }.`,
          en: `Form submitted${
            submittedAt
              ? ` (date: ${new Date(submittedAt).toLocaleDateString("en-US")})`
              : ""
          }.`,
        },
        ""
      )}
    </div>
  );
}
