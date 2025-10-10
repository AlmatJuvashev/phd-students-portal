// pages/doctoral.journey.tsx
import { WorldMap } from "@/components/map/WorldMap";
import playbook from "@/playbooks/playbook.json";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ResetBar } from "@/features/journey/components/ResetBar";

export function DoctoralJourney() {
  const { t: T } = useTranslation("common");
  const qc = useQueryClient();
  const {
    data: state = {},
    refetch,
    isLoading,
  } = useQuery({
    queryKey: ["journey", "state"],
    queryFn: () => api("/journey/state"),
  });

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
