// pages/doctoral.journey.tsx
import { WorldMap } from "@/components/map/WorldMap";
import playbook from "@/playbooks/playbook.json";
import { useJourneyState } from "@/features/journey/hooks";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ResetBar } from "@/features/journey/components/ResetBar";

export function DoctoralJourney() {
  const { t: T } = useTranslation("common");
  const { state = {}, refetch } = useJourneyState();

  return (
    <div>
      <WorldMap
        playbook={playbook as any}
        stateByNodeId={state as any}
        onStateChanged={() => refetch()}
      />
      <ResetBar />
    </div>
  );
}
