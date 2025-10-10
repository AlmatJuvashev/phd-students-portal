import { Button } from "@/components/ui/button";
import { useJourneyState } from "../hooks";
import { useTranslation } from "react-i18next";

export function ResetBar() {
  const { t: T } = useTranslation("common");
  const { reset } = useJourneyState();
  return (
    <div className="fixed inset-x-0 bottom-0 z-40">
      <div className="mx-auto max-w-6xl px-4 pb-4">
        <div className="rounded-xl bg-muted/70 backdrop-blur border p-3 flex items-center justify-between">
          <div className="text-xs text-muted-foreground">
            {T("journey.reset_hint", { defaultValue: "You can reset your journey to start over." })}
          </div>
          <Button
            variant="secondary"
            size="sm"
            onClick={() => {
              if (
                confirm(
                  T("journey.reset_confirm", {
                    defaultValue: "Reset your journey progress?",
                  })
                )
              ) {
                reset();
              }
            }}
          >
            {T("journey.reset", { defaultValue: "Reset Journey" })}
          </Button>
        </div>
      </div>
    </div>
  );
}

