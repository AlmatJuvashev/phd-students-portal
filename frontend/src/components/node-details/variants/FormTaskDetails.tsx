// components/node-details/variants/FormTaskDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { NodeVM, FieldDef, t } from "@/lib/playbook";
import { useEffect, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import * as Dialog from "@radix-ui/react-dialog";
import { useTranslation } from "react-i18next";
import { AssetsDownloads } from "../AssetsDownloads";
import { assetsForNode, allAssets } from "@/lib/assets";
import { evalVisible as evalVisibleExpr } from "@/features/forms/Visibility";
import { FieldRenderer } from "@/features/forms/FieldRenderer";
import { ActionsBar } from "@/features/forms/ActionsBar";
import StickyActions from "@/components/ui/sticky-actions";
import { useConditions } from "@/features/journey/useConditions";

type Props = {
  node: NodeVM;
  onSubmit?: (payload: any) => void;
  initial?: Record<string, any>;
  canEdit?: boolean;
  disabled?: boolean;
};

export function FormTaskDetails({
  node,
  initial = {},
  onSubmit,
  canEdit = true,
  disabled = false,
}: Props) {
  const [values, setValues] = useState<Record<string, any>>(initial);
  const [showOmidConfirm, setShowOmidConfirm] = useState(false);
  const [showNkConfirm, setShowNkConfirm] = useState(false);
  useEffect(() => {
    setValues(initial ?? {});
  }, [initial]);
  const { t: T, i18n } = useTranslation("common");

  const fields: FieldDef[] = node.requirements?.fields ?? [];

  const cardsLayout = (node.requirements as any)?.ui_hints?.cards_layout;
  const buttonsStyle = (node.requirements as any)?.ui_hints?.buttons_style;

  const evalVisible = (expr?: string) => evalVisibleExpr(values, expr);

  // Update setField to include hearing_happened cascading resets
  function setField(k: string, v: any) {
    setValues((prev) => {
      const next = { ...prev, [k]: v };
      if (k === "hearing_happened") {
        delete next.remarks_exist;
        delete next.plan_prepared;
        delete next.remarks_resolved;
      }
      if (k === "remarks_exist") {
        delete next.plan_prepared;
        delete next.remarks_resolved;
      }
      if (k === "plan_prepared") {
        delete next.remarks_resolved;
      }
      return next;
    });
  }

  // E3_hearing_nk: interactive cards with guard; proceed to RP1 or D1 based on rp_required
  if (node.id === "E3_hearing_nk") {
    const h = values["hearing_happened"];
    const r = values["remarks_exist"];
    const p = values["plan_prepared"];
    const z = values["remarks_resolved"];
    const { rp_required } = useConditions();
    const targetNext = rp_required
      ? "RP1_overview_actualization"
      : "D1_normokontrol_ncste";
    let currentStep:
      | "q0"
      | "hearingReminder"
      | "q1"
      | "q2"
      | "reminderA"
      | "q3"
      | "reminderB"
      | "done";
    if (h === undefined) currentStep = "q0";
    else if (h === false) currentStep = "hearingReminder";
    else if (r === undefined) currentStep = "q1";
    else if (r === false) currentStep = "done";
    else if (p === undefined) currentStep = "q2";
    else if (p === false) currentStep = "reminderA";
    else if (z === undefined) currentStep = "q3";
    else if (z === false) currentStep = "reminderB";
    else currentStep = "done";
    const canProceed =
      h === true && (r === false || (p === true && z === true));

    return (
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 h-full">
        {/* Left: Form (max 60%) */}
        <div className="lg:col-span-3 min-h-0 overflow-auto pr-1 space-y-3">
          <AnimatePresence initial={false}>
            <motion.div
              key={`nk-step-${currentStep}`}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: 10 }}
              transition={{ duration: 0.2 }}
            >
              {currentStep === "q0" && (
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
              {currentStep === "hearingReminder" && (
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
              {currentStep === "q1" && (
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
              {currentStep === "q2" && (
                <Card className="p-4">
                  <div className="mb-2 font-medium">
                    {T(
                      "forms.nk.q2",
                      "Подготовлен план устранения замечаний?"
                    )}
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
              {currentStep === "reminderA" && (
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
              {currentStep === "q3" && (
                <Card className="p-4">
                  <div className="mb-2 font-medium">
                    {T(
                      "forms.nk.q3",
                      "Замечания устранены согласно плану?"
                    )}
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
              {currentStep === "reminderB" && (
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
          {currentStep === "done" && (
            <div className="space-y-2">
              <div className="text-sm text-muted-foreground">
                {T("forms.nk.done_info", "Ответы зафиксированы. Вы можете продолжить.")}
              </div>
              <StickyActions
                primaryLabel={T("forms.proceed_next", "Перейти к следующему шагу")}
                onPrimary={() => onSubmit?.({ ...values, __nextOverride: targetNext })}
              />
            </div>
          )}
          {/* Sticky back actions for reminder steps */}
          {currentStep === "hearingReminder" && (
            <div className="pt-2">
              <StickyActions
                primaryLabel={T("forms.save_draft")}
                onPrimary={() => onSubmit?.({ ...values, __draft: true })}
                secondaryLabel={T("forms.nk.back", "Назад")}
                onSecondary={() => setField("hearing_happened", undefined)}
              />
            </div>
          )}
          {currentStep === "reminderA" && (
            <div className="pt-2">
              <StickyActions
                primaryLabel={T("forms.save_draft")}
                onPrimary={() => onSubmit?.({ ...values, __draft: true })}
                secondaryLabel={T("forms.nk.back", "Назад")}
                onSecondary={() => setField("plan_prepared", undefined)}
              />
            </div>
          )}
          {currentStep === "reminderB" && (
            <div className="pt-2">
              <StickyActions
                primaryLabel={T("forms.save_draft")}
                onPrimary={() => onSubmit?.({ ...values, __draft: true })}
                secondaryLabel={T("forms.nk.back", "Назад")}
                onSecondary={() => {
                  setField("plan_prepared", undefined);
                  setField("remarks_resolved", undefined);
                }}
              />
            </div>
          )}
        </div>
        {/* Right: Templates (40%), sticky */}
        <div className="lg:col-span-2 border-l pl-4 overflow-auto">
          <AssetsDownloads node={node} />
        </div>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 h-full">
      {/* Left: Form (<=60%) */}
      <Card className="p-4 space-y-4 lg:col-span-3 overflow-auto flex-1">
        {node.requirements?.notes && (
          <p className="text-sm text-muted-foreground">
            {node.requirements.notes}
          </p>
        )}
        <div className="space-y-3">
          {cardsLayout?.style === "stacked" && buttonsStyle === "yes_no" ? (
            // Stacked card-by-card yes/no flow
            <div className="space-y-3">
              <AnimatePresence initial={false}>
                {fields.map((f, index) => {
                  const visible = evalVisible((f as any).visible_when);
                  if (!visible) return null;

                  // Первая карточка
                  if (f.key === "remarks_exist") {
                    return (
                      <motion.div
                        key={f.key}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: 10 }}
                        transition={{ duration: 0.2 }}
                      >
                        <Card className="p-4">
                          <div className="mb-2 font-medium">
                            {t(f.label, f.key)}
                          </div>
                          <div className="flex gap-2">
                            <Button
                              onClick={() => {
                                setField(f.key, true);
                                onSubmit?.({
                                  ...values,
                                  [f.key]: true,
                                  __draft: true,
                                });
                              }}
                              disabled={!canEdit || disabled}
                            >
                              {T("forms.yes")}
                            </Button>
                            <Button
                              variant="secondary"
                              onClick={() => {
                                setField(f.key, false);
                                onSubmit?.({
                                  ...values,
                                  [f.key]: false,
                                  __draft: true,
                                });
                              }}
                              disabled={!canEdit || disabled}
                            >
                              {T("forms.no")}
                            </Button>
                          </div>
                        </Card>
                      </motion.div>
                    );
                  }

                  // Вторая карточка
                  if (
                    f.key === "plan_prepared" &&
                    values["remarks_exist"] === true
                  ) {
                    return (
                      <motion.div
                        key={f.key}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: 10 }}
                        transition={{ duration: 0.2 }}
                      >
                        <Card className="p-4">
                          <div className="mb-2 font-medium">
                            {t(f.label, f.key)}
                          </div>
                          <div className="flex gap-2">
                            <Button
                              onClick={() => {
                                setField(f.key, true);
                                onSubmit?.({
                                  ...values,
                                  [f.key]: true,
                                  __draft: true,
                                });
                              }}
                              disabled={!canEdit || disabled}
                            >
                              {T("forms.yes")}
                            </Button>
                            <Button
                              variant="secondary"
                              onClick={() => {
                                setField(f.key, false);
                                onSubmit?.({
                                  ...values,
                                  [f.key]: false,
                                  __draft: true,
                                });
                              }}
                              disabled={!canEdit || disabled}
                            >
                              {T("forms.no")}
                            </Button>
                          </div>
                        </Card>
                      </motion.div>
                    );
                  }

                  // Третья карточка
                  if (
                    f.key === "remarks_resolved" &&
                    values["plan_prepared"] === true
                  ) {
                    return (
                      <motion.div
                        key={f.key}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: 10 }}
                        transition={{ duration: 0.2 }}
                      >
                        <Card className="p-4">
                          <div className="mb-2 font-medium">
                            {t(f.label, f.key)}
                          </div>
                          <div className="flex gap-2">
                            <Button
                              onClick={() => {
                                setField(f.key, true);
                                onSubmit?.({
                                  ...values,
                                  [f.key]: true,
                                  __draft: true,
                                });
                              }}
                              disabled={!canEdit || disabled}
                            >
                              {T("forms.yes")}
                            </Button>
                            <Button
                              variant="secondary"
                              onClick={() => {
                                setField(f.key, false);
                                onSubmit?.({
                                  ...values,
                                  [f.key]: false,
                                  __draft: true,
                                });
                              }}
                              disabled={!canEdit || disabled}
                            >
                              {T("forms.no")}
                            </Button>
                          </div>
                        </Card>
                      </motion.div>
                    );
                  }

                  // Напоминание
                  if (
                    f.type === "note" &&
                    !values["remarks_resolved"] &&
                    values["plan_prepared"] === true
                  ) {
                    return (
                      <motion.div
                        key={f.key}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: 10 }}
                        transition={{ duration: 0.2 }}
                      >
                        <Card className="p-4 bg-gray-50">
                          <div className="text-sm text-muted-foreground">
                            {t(f.label, "")}
                          </div>
                          <Button
                            variant="secondary"
                            onClick={() => {
                              setField("plan_prepared", false);
                            }}
                          >
                            {T("forms.back")}
                          </Button>
                        </Card>
                      </motion.div>
                    );
                  }

                  return null;
                })}
              </AnimatePresence>

              {/* Кнопка перехода */}
              {((node.requirements as any)?.actions ?? []).map((a: any) => {
                const visible = evalVisible(a.visible_when);
                const disabled =
                  a.key === "go_to_ds_all_resolved" &&
                  !values["remarks_resolved"];
                if (!visible) return null;
                const label = t(a.label, "");
                return (
                  <div key={a.key} className="pt-2">
                    <Button
                      onClick={() => {
                        onSubmit?.({ ...values, __nextOverride: a.to });
                      }}
                      disabled={disabled}
                    >
                      {label}
                    </Button>
                  </div>
                );
              })}
              {/* Sticky back/save for stacked flow */}
              {(() => {
                const canBack =
                  values["remarks_resolved"] !== undefined ||
                  values["plan_prepared"] !== undefined ||
                  values["remarks_exist"] !== undefined;
                const handleBack = () => {
                  if (values["remarks_resolved"] !== undefined) {
                    setField("plan_prepared", undefined);
                    setField("remarks_resolved", undefined);
                  } else if (values["plan_prepared"] !== undefined) {
                    setField("remarks_exist", undefined);
                    setField("plan_prepared", undefined);
                  } else if (values["remarks_exist"] !== undefined) {
                    setField("remarks_exist", undefined);
                  }
                };
                return (
                  <div className="pt-2">
                    <StickyActions
                      primaryLabel={T("forms.save_draft")}
                      onPrimary={() => onSubmit?.({ ...values, __draft: true })}
                      secondaryLabel={canBack ? T("forms.back") : undefined}
                      onSecondary={canBack ? handleBack : undefined}
                      disabled={disabled}
                    />
                  </div>
                );
              })()}
            </div>
          ) : (
            <div className="space-y-3">
              {fields.map((f) => (
                <FieldRenderer
                  key={f.key}
                  field={f as any}
                  value={values[f.key]}
                  onChange={(v) => setField(f.key, v)}
                  setField={(k, v) => setField(k, v)}
                  otherValue={values[`${f.key}_other`]}
                  canEdit={canEdit}
                  disabled={disabled}
                />
              ))}
            </div>
          )}
        </div>

        {/* Templates / Downloads (if any) */}
        {/* Placeholder: Templates rendered on the right column */}

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

        {canEdit && (
          <div className="space-y-2">
            {node.id === "S1_publications_list" ? (
              <>
                <div className="text-sm font-medium">
                  {T("forms.app7_prompt")}{" "}
                  <span className="text-destructive">*</span>
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={() => onSubmit?.(values)}
                    disabled={disabled}
                  >
                    {T("forms.yes")}
                  </Button>
                  <Button
                    variant="secondary"
                    disabled={disabled}
                    onClick={() => {
                      // open preferred Appendix 7 template in a new tab
                      const assets = assetsForNode(node);
                      const lang =
                        (i18n.language as "ru" | "kz" | "en") || "ru";
                      const preferred =
                        assets.find(
                          (a) =>
                            a.id.toLowerCase().includes("app7") &&
                            a.id.toLowerCase().includes(`_${lang}`)
                        ) ||
                        assets.find((a) =>
                          a.id.toLowerCase().includes("app7")
                        ) ||
                        assets[0];
                      if (preferred?.storage?.key) {
                        window.open(
                          `/${preferred.storage.key}`,
                          "_blank",
                          "noopener,noreferrer"
                        );
                      }
                      // scroll to templates section to draw attention
                      document
                        .getElementById("templates-section")
                        ?.scrollIntoView({
                          behavior: "smooth",
                          block: "start",
                        });
                      // mark as draft so the user can return after preparing the document
                      onSubmit?.({ ...values, __draft: true });
                    }}
                  >
                    {T("forms.no")}
                  </Button>
                </div>
              </>
            ) : node.id === "E1_apply_omid" ? (
              <>
                <div className="text-sm font-medium">
                  {T("forms.omid_prompt")}{" "}
                  <span className="text-destructive">*</span>
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={() => setShowOmidConfirm(true)}
                    disabled={disabled}
                  >
                    {T("forms.yes")}
                  </Button>
                  <Button
                    variant="secondary"
                    disabled={disabled}
                    onClick={() => {
                      // Open OMiD application template matching current locale
                      const assets = assetsForNode(node);
                      const lang =
                        (i18n.language as "ru" | "kz" | "en") || "ru";
                      const preferred =
                        assets.find(
                          (a) =>
                            a.id.toLowerCase().includes("omid") &&
                            a.id.toLowerCase().includes(`_${lang}`)
                        ) ||
                        assets.find((a) =>
                          a.id.toLowerCase().includes("omid")
                        ) ||
                        assets[0];
                      if (preferred?.storage?.key) {
                        window.open(
                          `/${preferred.storage.key}`,
                          "_blank",
                          "noopener,noreferrer"
                        );
                      }
                      // Scroll to templates block and save draft
                      document
                        .getElementById("templates-section")
                        ?.scrollIntoView({
                          behavior: "smooth",
                          block: "start",
                        });
                      onSubmit?.({ ...values, __draft: true });
                    }}
                  >
                    {T("forms.no")}
                  </Button>
                </div>
              </>
            ) : node.id === "NK_package" ? (
              (() => {
                const ui = (node.requirements as any)?.ui_hints;
                const btnAll = ui?.buttons?.find(
                  (b: any) => b.key === "all_ready"
                );
                const btnNeed = ui?.buttons?.find(
                  (b: any) => b.key === "need_lcb_letter"
                );
                const allLabel = t(btnAll?.label, T("forms.save_submit"));
                const needLabel = t(btnNeed?.label, T("forms.save_draft"));
                const ready =
                  !!values["chk_thesis_unbound"] &&
                  !!values["chk_advisor_reviews"] &&
                  !!values["chk_pubs_app7"] &&
                  !!values["chk_sc_extract"] &&
                  !!values["chk_lcb_defense"];
                return (
                  <>
                    <div className="font-semibold">
                      {T("forms.nk_required_package_title")}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {t((node as any).description, "")}
                    </div>
                    <div className="flex gap-2">
                      <Button
                        onClick={() => setShowNkConfirm(true)}
                        disabled={disabled || !ready}
                      >
                        {allLabel}
                      </Button>
                      <Button
                        variant="secondary"
                        disabled={disabled}
                        onClick={() => {
                          // Open localized LCB request template from explicit templates
                          const explicit = (node.requirements as any)
                            ?.templates as string[] | undefined;
                          const pool = allAssets();
                          const lang =
                            (i18n.language as "ru" | "kz" | "en") || "ru";
                          const candidates = (
                            explicit?.length
                              ? explicit
                                  .map((id) => pool.find((a) => a.id === id))
                                  .filter(Boolean)
                              : pool.filter((a) =>
                                  a.id.toLowerCase().includes("lcb_request")
                                )
                          ) as ReturnType<typeof allAssets>;
                          const preferred =
                            candidates.find((a) =>
                              a.id.toLowerCase().includes(`_${lang}`)
                            ) || candidates[0];
                          if (preferred?.storage?.key) {
                            window.open(
                              `/${preferred.storage.key}`,
                              "_blank",
                              "noopener,noreferrer"
                            );
                          }
                          document
                            .getElementById("templates-section")
                            ?.scrollIntoView({
                              behavior: "smooth",
                              block: "start",
                            });
                          onSubmit?.({ ...values, __draft: true });
                        }}
                      >
                        {needLabel}
                      </Button>
                    </div>
                  </>
                );
              })()
            ) : (
              <ActionsBar
                onSubmit={() => onSubmit?.(values)}
                onDraft={() => onSubmit?.({ ...values, __draft: true })}
                disabled={disabled}
              />
            )}
          </div>
        )}
        <Dialog.Root open={showOmidConfirm} onOpenChange={setShowOmidConfirm}>
          <Dialog.Portal>
            <Dialog.Overlay className="fixed inset-0 bg-black/50 z-[70]" />
            <Dialog.Content className="fixed z-[70] left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg p-6 w-full max-w-md shadow-lg outline-none">
              <div className="mb-4 text-sm text-muted-foreground whitespace-pre-line">
                {T("forms.omid_info_after_yes")}
              </div>
              <div className="flex gap-2 justify-end">
                <Button
                  variant="outline"
                  onClick={() => setShowOmidConfirm(false)}
                >
                  {T("common.cancel")}
                </Button>
                <Button
                  onClick={() => {
                    setShowOmidConfirm(false);
                    onSubmit?.(values);
                  }}
                  disabled={disabled}
                >
                  {T("forms.proceed_next")}
                </Button>
              </div>
            </Dialog.Content>
          </Dialog.Portal>
        </Dialog.Root>
        {/* NK package confirmation */}
        <Dialog.Root open={showNkConfirm} onOpenChange={setShowNkConfirm}>
          <Dialog.Portal>
            <Dialog.Overlay className="fixed inset-0 bg-black/50 z-[70]" />
            <Dialog.Content className="fixed z-[70] left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg p-6 w-full max-w-md shadow-lg outline-none">
              <div className="mb-4 text-sm text-muted-foreground whitespace-pre-line">
                {T("forms.nk_confirm_info")}
              </div>
              <div className="flex gap-2 justify-end">
                <Button
                  variant="outline"
                  onClick={() => setShowNkConfirm(false)}
                >
                  {T("common.cancel")}
                </Button>
                <Button
                  onClick={() => {
                    setShowNkConfirm(false);
                    const nextOnComplete =
                      node.outcomes?.find((o) => o.value === "complete")
                        ?.next?.[0] ||
                      (Array.isArray(node.next) ? node.next[0] : undefined) ||
                      undefined;
                    onSubmit?.({ ...values, __nextOverride: nextOnComplete });
                  }}
                  disabled={disabled}
                >
                  {T("forms.proceed_next")}
                </Button>
              </div>
            </Dialog.Content>
          </Dialog.Portal>
        </Dialog.Root>
      </Card>
      {/* Right: Templates (40%), sticky container */}
      <div className="lg:col-span-2 border-l pl-4 overflow-auto">
        <AssetsDownloads node={node} />
      </div>
    </div>
  );
}
