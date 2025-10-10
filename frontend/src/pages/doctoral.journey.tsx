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
      <div className="max-w-4xl mx-auto px-4 pt-4">
        <a href="/" className="text-sm text-muted-foreground hover:underline">← {T("common.back", { defaultValue: "Back" })}</a>
      </div>
      <WorldMap
        playbook={playbook as any}
        stateByNodeId={state as any}
        onStateChanged={() => refetch()}
      />
      <ResetBar />
    </div>
  );
}
