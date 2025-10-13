/**
 * Refactored ChecklistDetails using unified FormProvider
 * Phase 1.4 - Example of form unification
 *
 * This is a cleaner, simplified version that demonstrates:
 * - Reusable state management via FormProvider
 * - Automatic field rendering via FormFields
 * - Unified action buttons via FormActions
 * - Status display via FormStatus
 */

import type { NodeVM } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import {
  FormProvider,
  FormFields,
  FormActions,
  FormStatus,
} from "@/features/forms";
import { TemplatesPanel } from "@/features/forms/TemplatesPanel";

export default function ChecklistDetailsUnified({
  node,
  initial = {},
  disabled,
  onSubmit,
  canEdit,
}: {
  node: NodeVM;
  initial?: Record<string, any>;
  disabled?: boolean;
  onSubmit?: (payload: any) => void;
  canEdit?: boolean;
}) {
  return (
    <FormProvider
      node={node}
      initial={initial}
      onSubmit={onSubmit}
      canEdit={canEdit}
      disabled={disabled}
    >
      <form className="h-full">
        <div className="lg:grid lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)] lg:gap-6 space-y-6 lg:space-y-0 min-w-0">
          <div className="space-y-4 min-w-0">
            {/* Description */}
            {Boolean((node as any)?.description) && (
              <div className="text-sm text-muted-foreground mb-4 p-4 rounded-lg bg-muted/30 border-l-4 border-primary/50">
                {t((node as any).description, "")}
              </div>
            )}

            {/* Auto-render all boolean fields as checklist items */}
            <FormFields types={["boolean"]} />

            {/* Action buttons with confirm modal */}
            <FormActions showConfirm={true} />

            {/* Status message for submitted forms */}
            <FormStatus />
          </div>

          {/* Templates panel */}
          <TemplatesPanel node={node} className="lg:border-l lg:pl-6" />
        </div>
      </form>
    </FormProvider>
  );
}
