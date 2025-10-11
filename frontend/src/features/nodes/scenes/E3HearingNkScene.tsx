import { useEffect, useMemo, useState } from "react";
import { NodeVM, FieldDef } from "@/lib/playbook";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import StickyActions from "@/components/ui/sticky-actions";
import { useTranslation } from "react-i18next";
import { AnimatePresence, motion } from "framer-motion";
import { useConditions } from "@/features/journey/useConditions";
import { AssetsDownloads } from "@/features/nodes/details/AssetsDownloads";

type Props = {
  node: NodeVM;
  onSubmit?: (payload: any) => void;
  initial?: Record<string, any>;
  canEdit?: boolean;
  disabled?: boolean;
};

export function E3HearingNkScene({
  node,
  initial = {},
  onSubmit,
  canEdit = true,
  disabled = false,
}: Props) {
  const [values, setValues] = useState<Record<string, any>>(initial);
  useEffect(() => {
    setValues(initial ?? {});
  }, [initial]);

  const { t: T } = useTranslation("common");
  const { rp_required } = useConditions();

  const setField = (key: string, value: any) => {
    setValues((prev) => {
      const next = { ...prev, [key]: value };
      if (key === "hearing_happened") {
        delete next.remarks_exist;
        delete next.plan_prepared;
        delete next.remarks_resolved;
      }
      if (key === "remarks_exist") {
        delete next.plan_prepared;
        delete next.remarks_resolved;
      }
      if (key === "plan_prepared") {
        delete next.remarks_resolved;
      }
      return next;
    });
  };

  const flow = useMemo(() => {
    const h = values["hearing_happened"];
    const r = values["remarks_exist"];
    const p = values["plan_prepared"];
    const z = values["remarks_resolved"];

    if (h === undefined) return "q0" as const;
    if (h === false) return "hearingReminder" as const;
    if (r === undefined) return "q1" as const;
    if (r === false) return "done" as const;
    if (p === undefined) return "q2" as const;
    if (p === false) return "reminderA" as const;
    if (z === undefined) return "q3" as const;
    if (z === false) return "reminderB" as const;
    return "done" as const;
  }, [values]);

  const targetNext = rp_required
    ? "RP1_overview_actualization"
    : "D1_normokontrol_ncste";
  const canProceed =
    values["hearing_happened"] === true &&
    (values["remarks_exist"] === false ||
      (values["plan_prepared"] === true && values["remarks_resolved"] === true));

  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 h-full">
      <div className="lg:col-span-3 min-h-0 overflow-auto pr-1 space-y-3">
        <AnimatePresence initial={false}>
          <motion.div
            key={`nk-step-${flow}`}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 10 }}
            transition={{ duration: 0.2 }}
          >
            {flow === "q0" && (
              <Card className="p-4">
                <div className="mb-2 font-medium">
                  {T("forms.nk.q0", "Заслушивание состоялось?")}
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={() => setField("hearing_happened", true)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.yes")}
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => setField("hearing_happened", false)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.no")}
                  </Button>
                </div>
              </Card>
            )}

            {flow === "hearingReminder" && (
              <Card className="p-4 bg-gray-50">
                <div className="text-sm text-muted-foreground">
                  {T(
                    "forms.nk.reminder_hearing",
                    "Назначьте и проведите заслушивание НК. Вернитесь к этому шагу после завершения."
                  )}
                </div>
                <Button
                  variant="secondary"
                  className="mt-2"
                  onClick={() => setField("hearing_happened", undefined)}
                  disabled={!canEdit || disabled}
                >
                  {T("forms.nk.back", "Назад")}
                </Button>
              </Card>
            )}

            {flow === "q1" && (
              <Card className="p-4">
                <div className="mb-2 font-medium">
                  {T(
                    "forms.nk.q1",
                    "Имеются зафиксированные замечания рецензентов/членов НК?"
                  )}
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={() => setField("remarks_exist", true)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.yes")}
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => setField("remarks_exist", false)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.no")}
                  </Button>
                </div>
                <div className="pt-2">
                  <Button
                    variant="secondary"
                    onClick={() => setField("hearing_happened", undefined)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.nk.back", "Назад")}
                  </Button>
                </div>
              </Card>
            )}

            {flow === "q2" && (
              <Card className="p-4">
                <div className="mb-2 font-medium">
                  {T("forms.nk.q2", "Подготовлен план устранения замечаний?")}
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={() => setField("plan_prepared", true)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.yes")}
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => setField("plan_prepared", false)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.no")}
                  </Button>
                </div>
                <div className="pt-2">
                  <Button
                    variant="secondary"
                    onClick={() => setField("remarks_exist", undefined)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.nk.back", "Назад")}
                  </Button>
                </div>
              </Card>
            )}

            {flow === "reminderA" && (
              <Card className="p-4 bg-gray-50">
                <div className="text-sm text-muted-foreground">
                  {T(
                    "forms.nk.reminder_plan",
                    "Подготовьте план устранения замечаний и вернитесь к этому шагу."
                  )}
                </div>
                <Button
                  variant="secondary"
                  className="mt-2"
                  onClick={() => setField("plan_prepared", undefined)}
                  disabled={!canEdit || disabled}
                >
                  {T("forms.nk.back", "Назад")}
                </Button>
              </Card>
            )}

            {flow === "q3" && (
              <Card className="p-4">
                <div className="mb-2 font-medium">
                  {T("forms.nk.q3", "Замечания устранены согласно плану?")}
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={() => setField("remarks_resolved", true)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.yes")}
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => setField("remarks_resolved", false)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.no")}
                  </Button>
                </div>
                <div className="pt-2">
                  <Button
                    variant="secondary"
                    onClick={() => setField("plan_prepared", undefined)}
                    disabled={!canEdit || disabled}
                  >
                    {T("forms.nk.back", "Назад")}
                  </Button>
                </div>
              </Card>
            )}

            {flow === "reminderB" && (
              <Card className="p-4 bg-gray-50">
                <div className="text-sm text-muted-foreground">
                  {T(
                    "forms.nk.reminder_resolve",
                    "Устраните все замечания согласно плану, затем подтвердите 'Да'."
                  )}
                </div>
                <Button
                  variant="secondary"
                  className="mt-2"
                  onClick={() => {
                    setField("plan_prepared", undefined);
                    setField("remarks_resolved", undefined);
                  }}
                  disabled={!canEdit || disabled}
                >
                  {T("forms.nk.back", "Назад")}
                </Button>
              </Card>
            )}
          </motion.div>
        </AnimatePresence>

        {flow === "done" && (
          <div className="space-y-2">
            <div className="text-sm text-muted-foreground">
              {T(
                "forms.nk.done_info",
                "Ответы зафиксированы. Вы можете продолжить."
              )}
            </div>
            <StickyActions
              primaryLabel={T("forms.proceed_next", "Перейти к следующему шагу")}
              onPrimary={() =>
                onSubmit?.({ ...values, __nextOverride: targetNext })
              }
              secondaryLabel={T("forms.save_draft")}
              onSecondary={() => onSubmit?.({ ...values, __draft: true })}
              disabled={!canProceed || disabled}
            />
          </div>
        )}

        {flow === "hearingReminder" && (
          <StickyActions
            primaryLabel={T("forms.save_draft")}
            onPrimary={() => onSubmit?.({ ...values, __draft: true })}
            secondaryLabel={T("forms.nk.back", "Назад")}
            onSecondary={() => setField("hearing_happened", undefined)}
            disabled={disabled}
          />
        )}

        {flow === "reminderA" && (
          <StickyActions
            primaryLabel={T("forms.save_draft")}
            onPrimary={() => onSubmit?.({ ...values, __draft: true })}
            secondaryLabel={T("forms.nk.back", "Назад")}
            onSecondary={() => setField("plan_prepared", undefined)}
            disabled={disabled}
          />
        )}

        {flow === "reminderB" && (
          <StickyActions
            primaryLabel={T("forms.save_draft")}
            onPrimary={() => onSubmit?.({ ...values, __draft: true })}
            secondaryLabel={T("forms.nk.back", "Назад")}
            onSecondary={() => {
              setField("plan_prepared", undefined);
              setField("remarks_resolved", undefined);
            }}
            disabled={disabled}
          />
        )}
      </div>

      <div className="lg:col-span-2 border-l pl-4">
        <AssetsDownloads node={node} />
      </div>
    </div>
  );
}

export default E3HearingNkScene;
