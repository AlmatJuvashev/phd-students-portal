// pages/doctoral.journey.tsx
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { JourneyMap } from "@/components/map/JourneyMap";
import { ResetBar } from "@/features/journey/components/ResetBar";
import { useJourneyState } from "@/features/journey/hooks";
import type { Playbook } from "@/lib/playbook";
import { useRequireAuth } from '@/hooks/useRequireAuth'

export function DoctoralJourney() {
  const { t: T, i18n } = useTranslation("common");
  const { isLoading } = useRequireAuth()
  const { state = {}, refetch } = useJourneyState();
  const [playbook, setPlaybook] = useState<Playbook | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    import("@/playbooks/playbook.json")
      .then((mod) => {
        if (!mounted) return;
        const next = (mod as { default?: Playbook }).default || (mod as Playbook);
        setPlaybook(next);
      })
      .finally(() => {
        if (mounted) setLoading(false);
      });
    return () => {
      mounted = false;
    };
  }, []);

  if (isLoading || loading || !playbook) {
    return (
      <div className="flex items-center justify-center py-16">
        <p className="text-sm text-muted-foreground animate-pulse">
          {T("map.loading", { defaultValue: "Loading dissertation map…" })}
        </p>
      </div>
    );
  }

  return (
    <div>
      <JourneyMap
        playbook={playbook}
        locale={i18n.language}
        stateByNodeId={state as any}
        onStateChanged={refetch}
      />
      <ResetBar />
    </div>
  );
}
