import React from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export type Filters = {
  q?: string;
  program?: string;
  department?: string;
  cohort?: string;
  advisor_id?: string;
  rp_required?: boolean;
  overdue?: boolean;
  due_from?: string;
  due_to?: string;
};

export function FiltersBar({ value, onChange, onRefresh }: { value: Filters; onChange: (f: Filters) => void; onRefresh: () => void }) {
  const [local, setLocal] = React.useState<Filters>(value);
  React.useEffect(() => setLocal(value), [value]);

  return (
    <div className="sticky top-14 z-10 bg-background/80 backdrop-blur border-b p-3">
      <div className="flex flex-wrap items-center gap-2">
        <Input placeholder="Search (name, email, phone)" value={local.q || ""} onChange={(e) => setLocal({ ...local, q: e.target.value })} className="w-64" />
        <Input placeholder="Program" value={local.program || ""} onChange={(e) => setLocal({ ...local, program: e.target.value })} />
        <Input placeholder="Department" value={local.department || ""} onChange={(e) => setLocal({ ...local, department: e.target.value })} />
        <Input placeholder="Cohort" value={local.cohort || ""} onChange={(e) => setLocal({ ...local, cohort: e.target.value })} />
        <label className="flex items-center gap-2 text-sm px-2">
          <input type="checkbox" checked={!!local.rp_required} onChange={(e) => setLocal({ ...local, rp_required: e.target.checked })} />
          RP required only
        </label>
        <label className="flex items-center gap-2 text-sm px-2">
          <input type="checkbox" checked={!!local.overdue} onChange={(e) => setLocal({ ...local, overdue: e.target.checked })} />
          Overdue only
        </label>
        <Input type="date" value={local.due_from || ""} onChange={e => setLocal({ ...local, due_from: e.target.value })} placeholder="From" />
        <Input type="date" value={local.due_to || ""} onChange={e => setLocal({ ...local, due_to: e.target.value })} placeholder="To" />
        <Button variant="outline" onClick={() => onChange(local)}>Apply</Button>
        <Button variant="ghost" onClick={() => { setLocal({}); onChange({}); }}>Clear</Button>
        <div className="ml-auto flex items-center gap-2">
          <Button variant="outline" onClick={onRefresh}>Refresh</Button>
          <Button onClick={() => exportCSV(local)}>Export CSV</Button>
        </div>
      </div>
    </div>
  );
}

function exportCSV(filters: Filters) {
  // Signal to caller via DOM event; page listens and provides current rows snapshot
  const ev = new CustomEvent('students-monitor:export', { detail: { filters } });
  window.dispatchEvent(ev);
}
