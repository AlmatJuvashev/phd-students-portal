import React from "react";

export function StatusChip({ state }: { state: string }) {
  const map: Record<string, string> = {
    locked: "bg-gray-100 text-gray-700",
    active: "bg-blue-100 text-blue-800",
    submitted: "bg-amber-100 text-amber-800",
    waiting: "bg-amber-100 text-amber-800",
    under_review: "bg-purple-100 text-purple-800",
    needs_fixes: "bg-red-100 text-red-800",
    done: "bg-green-100 text-green-800",
  };
  const label: Record<string, string> = {
    locked: "Locked",
    active: "Active",
    submitted: "Submitted",
    waiting: "Waiting",
    under_review: "Under review",
    needs_fixes: "Needs fixes",
    done: "Done",
  };
  const cls = map[state] || "bg-gray-100 text-gray-800";
  return (
    <span className={`px-2 py-0.5 rounded text-xs ${cls}`}>
      {label[state] || state}
    </span>
  );
}
