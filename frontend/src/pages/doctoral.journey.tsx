// pages/doctoral.journey.tsx
import { WorldMap } from "@/components/map/WorldMap";
import playbook from "@/playbooks/phd-doctorant.kz-v1.json";

export function DoctoralJourney() {
  // Example: hydrate with some fake node states; replace with real data from your backend
  const stateByNodeId = {
    S1_profile: "done",
    S1_text_ready: "done",
    S1_antiplag: "submitted",
    S1_publications_list: "active",
    E2_wait_30_days: "waiting",
  } as const;

  return (
    <WorldMap playbook={playbook as any} stateByNodeId={stateByNodeId as any} />
  );
}
