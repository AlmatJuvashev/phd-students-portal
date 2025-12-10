// components/node-details/NodeDetailSwitch.tsx
import { NodeVM } from "@/lib/playbook";
import type { NodeSubmissionDTO } from "@/api/journey";
import ConfirmTaskDetails from "./variants/ConfirmTaskDetails";
import InfoDetails from "./variants/InfoDetails";
import { deriveNodeKind } from "@/features/nodes/deriveNodeKind";
import FormEntryDetails from "@/features/nodes/kinds/FormEntryDetails";
import ChecklistDetails from "@/features/nodes/kinds/ChecklistDetails";
import CardsDetails from "@/features/nodes/kinds/CardsDetails";
import React, { Suspense, lazy } from "react";
import { FormTaskDetails } from "./variants/FormTaskDetails";
import useGuideForNode from "@/features/guides/useGuideForNode";

type SceneProps = {
  node: NodeVM;
  initial?: Record<string, any>;
  disabled?: boolean;
  canEdit?: boolean;
  onSubmit?: (payload: any) => void;
};

type FormRendererArgs = {
  node: NodeVM;
  initialForm: Record<string, any>;
  saving: boolean;
  canEdit?: boolean;
  onEvent?: (evt: { type: string; payload?: any }) => void;
};

const sceneLoaders: Record<
  string,
  () => Promise<{ default: React.ComponentType<SceneProps> }>
> = {
  D2_apply_to_ds: () => import("@/features/nodes/scenes/D2ApplyToDsScene"),
  V1_reinstatement_package: () =>
    import("@/features/nodes/scenes/V1ReinstatementScene"),
  RP2_sc_hearing_prep: () =>
    import("@/features/nodes/scenes/RP2HearingPrepScene"),

  S1_profile: () => import("@/features/nodes/scenes/S1ProfileDetails"),
  S1_publications_list: () =>
    import("@/features/nodes/scenes/S1PublicationsDetails"),
  E1_apply_omid: () => import("@/features/nodes/scenes/E1ApplyOmidDetails"),
  NK_package: () => import("@/features/nodes/scenes/NkPackageDetails"),
  E3_hearing_nk: () => import("@/features/nodes/scenes/E3HearingNkScene"),
};

const sceneComponents: Record<
  string,
  React.LazyExoticComponent<React.ComponentType<SceneProps>>
> = Object.fromEntries(
  Object.entries(sceneLoaders).map(([id, loader]) => [
    id,
    lazy(() =>
      loader().then((mod) => ({
        default: mod.default || (mod as any),
      }))
    ),
  ])
);

const formRenderers: Record<string, (args: FormRendererArgs) => JSX.Element> = {
  formEntry: ({ node, initialForm, saving, canEdit, onEvent }) => (
    <FormEntryDetails
      node={node}
      initial={initialForm}
      disabled={saving || canEdit === false}
      onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
    />
  ),
  checklist: ({ node, initialForm, saving, canEdit, onEvent }) => (
    <ChecklistDetails
      node={node}
      initial={initialForm}
      disabled={saving}
      canEdit={canEdit}
      onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
    />
  ),
  cards: ({ node, initialForm, saving, onEvent }) => (
    <CardsDetails
      node={node}
      disabled={saving}
      onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
    />
  ),
};

type Props = {
  node: NodeVM;
  onEvent?: (evt: { type: string; payload?: any }) => void; // bubble up submit/finalize/etc.
  submission?: NodeSubmissionDTO | null;
  saving?: boolean;
  canEdit?: boolean;
  onAttachmentsRefresh?: () => void;
};

export function NodeDetailSwitch({
  node,
  onEvent,
  submission,
  saving = false,
  canEdit,
  onAttachmentsRefresh,
}: Props) {
  const uiKind = deriveNodeKind(node);
  const renderGuide = useGuideForNode(node) || undefined;
  const initialForm = submission?.form?.data ?? {};

  const SceneComponent = sceneComponents[node.id];
  if (SceneComponent) {
    return (
      <Suspense fallback={<div className="p-4 text-sm">Loading…</div>}>
        <SceneComponent
          node={node}
          initial={initialForm}
          disabled={saving}
          canEdit={canEdit ?? !saving}
          onSubmit={(payload: any) =>
            onEvent?.({ type: "submit-form", payload })
          }
        />
      </Suspense>
    );
  }

  if (node.type === "form") {
    const renderer = uiKind ? formRenderers[uiKind] : undefined;
    if (renderer) {
      return renderer({
        node,
        initialForm,
        saving,
        canEdit,
        onEvent,
      });
    }
  }

  if (node.type === "confirmTask") {
    return (
      <ConfirmTaskDetails
        node={node}
        slots={submission?.slots}
        canEdit={canEdit}
        onRefresh={onAttachmentsRefresh}
        onComplete={() =>
          onEvent?.({
            type: "submit-decision",
            payload: {},
          })
        }
        onReset={() => onEvent?.({ type: "reset-node" })}
      />
    );
  }

  if (node.type === "info") {
    return (
      <InfoDetails
        node={node}
        renderGuide={renderGuide}
        onContinue={() => onEvent?.({ type: "continue" })}
      />
    );
  }

  // Default to a simple form view when type is unrecognised.
  return (
    <FormTaskDetails
      node={node}
      canEdit={canEdit ?? !saving}
      initial={initialForm}
      disabled={saving}
      onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
    />
  );
}
