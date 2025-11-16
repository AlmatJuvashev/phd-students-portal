// components/map/NodeToken.tsx
import { Badge } from "@/components/ui/badge";
import {
  LucideIcon,
  FormInput,
  Upload,
  GitMerge,
  Users,
  Hourglass,
  ExternalLink,
  Trophy,
  MapPin,
  Check,
  Lock,
  FileText,
  FileCheck2,
  ShieldCheck,
  ListChecks,
  ClipboardCheck,
  Megaphone,
  GraduationCap,
  Award,
  Package,
  FileSignature,
  RefreshCw,
  ClipboardList,
  BookOpen,
  ScrollText,
  CheckCircle2,
  Clock,
  AlertCircle,
  ArrowRight,
} from "lucide-react";
import { NodeVM, t } from "@/lib/playbook";
import clsx from "clsx";

const typeIcon: Record<NodeVM["type"], LucideIcon> = {
  form: ClipboardList,
  upload: Upload,
  decision: GitMerge,
  meeting: Users,
  waiting: Clock,
  external: ExternalLink,
  boss: Trophy,
  gateway: ArrowRight,
  info: BookOpen,
  confirmTask: CheckCircle2,
};

// More accurate icons by node id when available
const idIcon: Record<string, LucideIcon> = {
  // Section 1 — Student preparation
  S1_text_ready: ScrollText,
  S1_antiplag: ShieldCheck,
  S1_publications_list: ListChecks,
  // External application to OMiD
  E1_apply_omid: ClipboardCheck,
  // Hearing at NK (department/committee)
  E3_hearing_nk: Users,
  // NCSTE (НЦГНТЭ) steps
  D1_normokontrol_ncste: FileCheck2,
  IV3_publication_certificate_ncste: Award,
  NK_package: Package,
  // DS application / reinstatement
  D2_apply_to_ds: FileSignature,
  V1_reinstatement_package: RefreshCw,
  // Special scenes
  RP2_sc_hearing_prep: Users,
  VI_attestation_file: FileCheck2,
};

const stateStyles = {
  locked: {
    iconBg:
      "bg-gradient-to-br from-gray-300 to-gray-400 dark:from-gray-600 dark:to-gray-700",
    iconColor: "text-gray-600 dark:text-gray-400",
    badge: "bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300",
    ring: "",
    opacity: "opacity-60",
  },
  active: {
    iconBg: "bg-gradient-to-br from-primary via-primary to-blue-600",
    iconColor: "text-white",
    badge: "bg-primary/20 text-primary",
    ring: "ring-4 ring-primary/25 shadow-[0_0_12px_rgba(56,139,253,0.25)]",
    opacity: "",
  },
  submitted: {
    iconBg: "bg-gradient-to-br from-amber-400 to-amber-600",
    iconColor: "text-white",
    badge: "bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300",
    ring: "",
    opacity: "",
  },
  waiting: {
    iconBg: "bg-gradient-to-br from-blue-400 to-blue-600",
    iconColor: "text-white",
    badge: "bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300",
    ring: "",
    opacity: "",
  },
  under_review: {
    iconBg: "bg-gradient-to-br from-purple-400 to-purple-600",
    iconColor: "text-white",
    badge:
      "bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300",
    ring: "",
    opacity: "",
  },
  needs_fixes: {
    iconBg: "bg-gradient-to-br from-red-400 to-red-600",
    iconColor: "text-white",
    badge: "bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-300",
    ring: "",
    opacity: "",
  },
  done: {
    iconBg: "bg-gradient-to-br from-green-400 to-green-600",
    iconColor: "text-white",
    badge: "bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300",
    ring: "",
    opacity: "",
  },
};

export function NodeToken({
  node,
  onClick,
}: {
  node: NodeVM;
  onClick?: (n: NodeVM) => void;
}) {
  const Icon = idIcon[node.id] || typeIcon[node.type];
  const styles = stateStyles[node.state];

  const isBossNode = node.type === "boss";
  const isClickable = node.state !== "locked";
  const isDone = node.state === "done";

  return (
    <button
      type="button"
      onClick={() => isClickable && onClick?.(node)}
      disabled={!isClickable}
      aria-label={`${t(node.title, node.id)} - ${node.state}`}
      aria-disabled={!isClickable}
      className={clsx(
        "flex items-center gap-4 relative group transition-all duration-200 w-full text-left",
        "min-h-[48px] sm:min-h-[52px] touch-manipulation", // Minimum touch target
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-lg",
        styles.opacity,
        {
          "cursor-pointer hover:scale-[1.02] active:scale-[0.98]": isClickable,
          "cursor-not-allowed": !isClickable,
        }
      )}
    >
      <div
        className={clsx("z-10 relative transition-transform duration-300", {
          "transform scale-110": isBossNode,
          "group-hover:scale-110": isClickable && !isBossNode,
        })}
      >
        {isBossNode && node.state === "active" && (
          <div className="absolute -inset-2 rounded-full bg-gradient-to-r from-yellow-400 via-amber-400 to-yellow-400 blur-sm animate-glow"></div>
        )}
        <div
          className={clsx(
            "relative w-16 h-16 sm:w-18 sm:h-18 rounded-full flex items-center justify-center transition-all duration-300",
            {
              "shadow-lg hover:shadow-2xl": isClickable,
              "shadow-md": !isClickable,
              "ring-2 ring-white/50 dark:ring-gray-700/50":
                !isDone && node.state !== "locked",
              "ring-2 ring-green-300/60 dark:ring-green-700/60": isDone,
            },
            styles.iconBg,
            styles.ring,
            { "boss-node": isBossNode }
          )}
        >
          <Icon
            className={clsx(
              "h-10 w-10 sm:h-11 sm:w-11 transition-transform group-hover:scale-110 duration-200",
              styles.iconColor
            )}
            strokeWidth={isDone ? 2.5 : 2}
            fill={isDone ? "currentColor" : "none"}
            fillOpacity={isDone ? 0.2 : 0}
          />
          {isDone && (
            <div className="absolute -bottom-0.5 -right-0.5 bg-gradient-to-br from-green-400 to-green-600 text-white rounded-full w-7 h-7 flex items-center justify-center shadow-lg ring-2 ring-white dark:ring-gray-800 animate-in zoom-in duration-300">
              <Check className="w-4 h-4" strokeWidth={3} />
            </div>
          )}
        </div>
      </div>

      <div className={clsx("flex-1 min-w-0", { "ml-1": isBossNode })}>
        <h3
          className={clsx(
            "font-bold text-base sm:text-lg leading-tight transition-colors duration-200",
            {
              "text-primary": node.state === "active",
              "text-green-700 dark:text-green-400": isDone,
              "group-hover:text-primary":
                isClickable && node.state !== "active" && !isDone,
              truncate: !isBossNode,
              "line-clamp-2": isBossNode,
            }
          )}
        >
          {t(node.title, node.id)}
        </h3>
        <div
          className={clsx(
            "mt-1.5 text-xs sm:text-sm font-semibold px-3 py-1.5 rounded-full inline-flex items-center gap-1.5 shadow-sm transition-all duration-200",
            styles.badge,
            "group-hover:shadow-md"
          )}
        >
          {node.state === "locked" && (
            <Lock className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
          )}
          {node.state === "active" && (
            <AlertCircle className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
          )}
          {node.state === "done" && (
            <CheckCircle2 className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
          )}
          {node.state === "waiting" && (
            <Clock className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
          )}
          <span className="capitalize">{node.state.replace("_", " ")}</span>
        </div>
      </div>
    </button>
  );
}
