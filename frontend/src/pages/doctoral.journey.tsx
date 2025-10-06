// pages/doctoral.journey.tsx
import { WorldMap } from "@/components/map/WorldMap";
import playbook from "@/playbooks/phd-doctorant.kz-v1.json";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export function DoctoralJourney() {
  const { t: T } = useTranslation("common");
  const qc = useQueryClient();
  const { data: state = {}, refetch, isLoading } = useQuery({
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
      <div className="p-4">
        <Button
          variant="secondary"
          onClick={async () => {
            if (!confirm(T("journey.reset_confirm", { defaultValue: "Reset your journey progress?" }))) return;
            await api("/journey/reset", { method: "POST" });
            await refetch();
          }}
        >
          {T("journey.reset", { defaultValue: "Reset Journey" })}
        </Button>
      </div>
    </div>
  );
}
