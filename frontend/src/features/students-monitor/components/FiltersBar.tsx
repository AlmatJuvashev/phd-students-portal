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
import { useTranslation } from "react-i18next";

const PROGRAM_OPTIONS = [
  { value: "PhD Computer Science", key: "phd_cs", fallback: "PhD Computer Science" },
  { value: "PhD Physics", key: "phd_physics", fallback: "PhD Physics" },
  { value: "PhD Chemistry", key: "phd_chemistry", fallback: "PhD Chemistry" },
  {
    value: "PhD Biomedical Engineering",
    key: "phd_biomed",
    fallback: "PhD Biomedical Engineering",
  },
  {
    value: "PhD Applied Mathematics",
    key: "phd_math",
    fallback: "PhD Applied Mathematics",
  },
];

const DEPARTMENT_OPTIONS = [
  { value: "Computer Science", key: "cs", fallback: "Computer Science" },
  { value: "Physics", key: "physics", fallback: "Physics" },
  { value: "Chemistry", key: "chemistry", fallback: "Chemistry" },
  {
    value: "Biomedical Engineering",
    key: "biomed",
    fallback: "Biomedical Engineering",
  },
  {
    value: "Applied Mathematics",
    key: "math",
    fallback: "Applied Mathematics",
  },
];

const COHORT_OPTIONS = ["2025", "2024", "2023", "2022"];

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
  const { t } = useTranslation("common");
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
              placeholder={t("admin.monitor.filters.search_placeholder", {
                defaultValue: "Search by name, ID, phone, or email...",
              })}
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
              <SelectValue
                placeholder={t("admin.monitor.filters.program_placeholder", {
                  defaultValue: "Program",
                })}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">
                {t("admin.monitor.filters.program_all", {
                  defaultValue: "All Programs",
                })}
              </SelectItem>
              {PROGRAM_OPTIONS.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {t(`admin.monitor.filters.programs.${option.key}`, {
                    defaultValue: option.fallback,
                  })}
                </SelectItem>
              ))}
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
              <SelectValue
                placeholder={t("admin.monitor.filters.department_placeholder", {
                  defaultValue: "Department",
                })}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">
                {t("admin.monitor.filters.department_all", {
                  defaultValue: "All Departments",
                })}
              </SelectItem>
              {DEPARTMENT_OPTIONS.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {t(`admin.monitor.filters.departments.${option.key}`, {
                    defaultValue: option.fallback,
                  })}
                </SelectItem>
              ))}
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
              <SelectValue
                placeholder={t("admin.monitor.filters.cohort_placeholder", {
                  defaultValue: "Cohort",
                })}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">
                {t("admin.monitor.filters.cohort_all", {
                  defaultValue: "All Cohorts",
                })}
              </SelectItem>
              {COHORT_OPTIONS.map((cohort) => (
                <SelectItem key={cohort} value={cohort}>
                  {cohort}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {/* More Filters Button */}
          <Button variant="outline" size="sm" className="gap-2 h-10">
            <Filter className="h-4 w-4" />
            {t("admin.monitor.filters.more", {
              defaultValue: "More Filters",
            })}
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
              {t("admin.monitor.filters.rp_only", {
                defaultValue: "RP required only",
              })}
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
              {t("admin.monitor.filters.overdue_only", {
                defaultValue: "Overdue only",
              })}
            </label>
          </div>

          <div className="flex-1" />

          {/* Action Buttons */}
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" onClick={handleClear}>
              {t("admin.monitor.filters.clear", { defaultValue: "Clear" })}
            </Button>
            <Button variant="default" size="sm" onClick={handleApply}>
              {t("admin.monitor.filters.apply", {
                defaultValue: "Apply Filters",
              })}
            </Button>
            <div className="h-6 w-px bg-border mx-1" />
            <Button
              variant="outline"
              size="sm"
              className="gap-2"
              onClick={() => exportCSV(local)}
            >
              <Download className="h-4 w-4" />
              {t("admin.monitor.filters.export", {
                defaultValue: "Export CSV",
              })}
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="gap-2"
              onClick={() => bulkReminder(local)}
            >
              <Mail className="h-4 w-4" />
              {t("admin.monitor.filters.bulk", {
                defaultValue: "Bulk Message",
              })}
            </Button>
            <Button
              size="sm"
              className="gap-2"
              onClick={() => bulkReminder(local)}
            >
              <Plus className="h-4 w-4" />
              {t("admin.monitor.filters.new_reminder", {
                defaultValue: "New Reminder",
              })}
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
