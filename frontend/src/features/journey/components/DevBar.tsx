import React from "react";
import { Button } from "@/components/ui/button";

export function DevBar() {
  const [unlockAll, setUnlockAll] = React.useState<boolean>(() => {
    try {
      return localStorage.getItem("dev_unlock_all_nodes") === "true";
    } catch {
      return false;
    }
  });

  const toggleUnlock = () => {
    const next = !unlockAll;
    setUnlockAll(next);
    try {
      localStorage.setItem("dev_unlock_all_nodes", String(next));
    } catch {}
    // Reload to allow maps/hooks to re-evaluate
    setTimeout(() => window.location.reload(), 50);
  };

  const clearCongrats = () => {
    try {
      sessionStorage.removeItem("journey_congrats_shown");
    } catch {}
  };

  return (
    <div className="fixed right-3 bottom-20 z-40">
      <div className="rounded-lg border bg-background/90 backdrop-blur p-3 text-xs shadow-md">
        <div className="font-semibold mb-2">Dev Tools</div>
        <label className="inline-flex items-center gap-2 mb-2">
          <input
            type="checkbox"
            checked={unlockAll}
            onChange={toggleUnlock}
          />
          <span>Unlock all nodes</span>
        </label>
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" onClick={clearCongrats}>
            Clear congrats
          </Button>
        </div>
      </div>
    </div>
  );
}

export default DevBar;

