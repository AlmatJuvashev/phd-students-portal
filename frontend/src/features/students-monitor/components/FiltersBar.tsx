import React from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";
import { Search, Download, Mail, Plus, Filter } from "lucide-react";

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

export function FiltersBar({
  value,
  onChange,
  onRefresh,
}: {
  value: Filters;
  onChange: (f: Filters) => void;
  onRefresh: () => void;
}) {
  const [local, setLocal] = React.useState<Filters>(value);
  React.useEffect(() => setLocal(value), [value]);

  const handleApply = () => {
    onChange(local);
  };

  const handleClear = () => {
    const cleared = {};
    setLocal(cleared);
    onChange(cleared);
  };

  return (
    <div className="sticky top-0 z-40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 border-b">
      <div className="px-8 py-4">
        <div className="flex items-center gap-3 flex-wrap">
          {/* Search Input */}
          <div className="relative flex-1 min-w-[320px]">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search by name, ID, phone, or email..."
              className="pl-10 h-10"
              value={local.q || ""}
              onChange={(e) => setLocal({ ...local, q: e.target.value })}
              onKeyDown={(e) => e.key === "Enter" && handleApply()}
            />
          </div>

          {/* Program Filter */}
          <Select
            value={local.program || "all"}
            onValueChange={(val) =>
              setLocal({ ...local, program: val === "all" ? undefined : val })
            }
          >
            <SelectTrigger className="w-[180px] h-10">
              <SelectValue placeholder="Program" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Programs</SelectItem>
              <SelectItem value="PhD Computer Science">
                PhD Computer Science
              </SelectItem>
              <SelectItem value="PhD Physics">PhD Physics</SelectItem>
              <SelectItem value="PhD Chemistry">PhD Chemistry</SelectItem>
              <SelectItem value="PhD Biomedical Engineering">
                PhD Biomedical Engineering
              </SelectItem>
              <SelectItem value="PhD Applied Mathematics">
                PhD Applied Mathematics
              </SelectItem>
            </SelectContent>
          </Select>

          {/* Department Filter */}
          <Select
            value={local.department || "all"}
            onValueChange={(val) =>
              setLocal({
                ...local,
                department: val === "all" ? undefined : val,
              })
            }
          >
            <SelectTrigger className="w-[180px] h-10">
              <SelectValue placeholder="Department" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Departments</SelectItem>
              <SelectItem value="Computer Science">Computer Science</SelectItem>
              <SelectItem value="Physics">Physics</SelectItem>
              <SelectItem value="Chemistry">Chemistry</SelectItem>
              <SelectItem value="Biomedical Engineering">
                Biomedical Engineering
              </SelectItem>
              <SelectItem value="Applied Mathematics">
                Applied Mathematics
              </SelectItem>
            </SelectContent>
          </Select>

          {/* Cohort Filter */}
          <Select
            value={local.cohort || "all"}
            onValueChange={(val) =>
              setLocal({ ...local, cohort: val === "all" ? undefined : val })
            }
          >
            <SelectTrigger className="w-[150px] h-10">
              <SelectValue placeholder="Cohort" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Cohorts</SelectItem>
              <SelectItem value="2025">2025</SelectItem>
              <SelectItem value="2024">2024</SelectItem>
              <SelectItem value="2023">2023</SelectItem>
              <SelectItem value="2022">2022</SelectItem>
            </SelectContent>
          </Select>

          {/* More Filters Button */}
          <Button variant="outline" size="sm" className="gap-2 h-10">
            <Filter className="h-4 w-4" />
            More Filters
          </Button>

          {/* RP Required Toggle */}
          <div className="flex items-center gap-2 ml-2">
            <Checkbox
              id="rp-only"
              checked={!!local.rp_required}
              onCheckedChange={(checked) =>
                setLocal({ ...local, rp_required: !!checked })
              }
            />
            <label
              htmlFor="rp-only"
              className="text-sm text-muted-foreground cursor-pointer"
            >
              RP required only
            </label>
          </div>

          {/* Overdue Toggle */}
          <div className="flex items-center gap-2">
            <Checkbox
              id="overdue-only"
              checked={!!local.overdue}
              onCheckedChange={(checked) =>
                setLocal({ ...local, overdue: !!checked })
              }
            />
            <label
              htmlFor="overdue-only"
              className="text-sm text-muted-foreground cursor-pointer"
            >
              Overdue only
            </label>
          </div>

          <div className="flex-1" />

          {/* Action Buttons */}
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={handleClear}>
              Clear
            </Button>
            <Button variant="default" size="sm" onClick={handleApply}>
              Apply Filters
            </Button>
            <div className="h-6 w-px bg-border mx-1" />
            <Button
              variant="outline"
              size="sm"
              className="gap-2"
              onClick={() => exportCSV(local)}
            >
              <Download className="h-4 w-4" />
              Export CSV
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="gap-2"
              onClick={() => bulkReminder(local)}
            >
              <Mail className="h-4 w-4" />
              Bulk Message
            </Button>
            <Button
              size="sm"
              className="gap-2"
              onClick={() => bulkReminder(local)}
            >
              <Plus className="h-4 w-4" />
              New Reminder
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

function exportCSV(filters: Filters) {
  // Signal to caller via DOM event; page listens and provides current rows snapshot
  const ev = new CustomEvent("students-monitor:export", {
    detail: { filters },
  });
  window.dispatchEvent(ev);
}

function bulkReminder(filters: Filters) {
  const ev = new CustomEvent("students-monitor:bulk-reminder", {
    detail: { filters },
  });
  window.dispatchEvent(ev);
}
