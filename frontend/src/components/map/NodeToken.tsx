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
} from "lucide-react";
import { NodeVM, t } from "@/lib/playbook";
import clsx from "clsx";

const typeIcon: Record<NodeVM["type"], LucideIcon> = {
  form: FormInput,
  upload: Upload,
  decision: GitMerge,
  meeting: Users,
  waiting: Hourglass,
  external: ExternalLink,
  boss: Trophy,
  gateway: MapPin,
  info: MapPin,
  confirmTask: Check,
  uploadTask: Upload,
};

// More accurate icons by node id when available
const idIcon: Record<string, LucideIcon> = {
  // Section 1 — Student preparation
  S1_text_ready: FileText,
  S1_antiplag: ShieldCheck,
  S1_publications_list: ListChecks,
  // External application to OMiD
  E1_apply_omid: ClipboardCheck,
  // Hearing at NK (department/committee)
  E3_hearing_nk: Megaphone,
  // NCSTE (НЦГНТЭ) steps
  D1_normokontrol_ncste: GraduationCap,
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
    iconBg: "bg-gray-300 dark:bg-gray-600",
    iconColor: "text-gray-500 dark:text-gray-400",
    badge: "bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300",
    ring: "",
    opacity: "opacity-60",
  },
  active: {
    iconBg: "bg-primary",
    iconColor: "text-white",
    badge: "bg-primary/20 text-primary",
    ring: "ring-4 ring-primary/30 animate-pulse",
    opacity: "",
  },
  submitted: {
    iconBg: "bg-amber-500",
    iconColor: "text-white",
    badge: "bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300",
    ring: "",
    opacity: "",
  },
  waiting: {
    iconBg: "bg-blue-500",
    iconColor: "text-white",
    badge: "bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300",
    ring: "",
    opacity: "",
  },
  needs_fixes: {
    iconBg: "bg-red-500",
    iconColor: "text-white",
    badge: "bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-300",
    ring: "",
    opacity: "",
  },
  done: {
    iconBg: "bg-green-500",
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

  return (
    <div
      role="button"
      onClick={() => isClickable && onClick?.(node)}
      className={clsx(
        "flex items-center gap-4 relative group transition-all duration-200",
        styles.opacity,
        {
          "cursor-pointer hover:scale-[1.02]": isClickable,
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
          <div className="absolute -inset-1.5 rounded-full bg-yellow-400 animate-glow"></div>
        )}
        <div
          className={clsx(
            "relative w-16 h-16 sm:w-18 sm:h-18 rounded-full flex items-center justify-center shadow-lg hover:shadow-xl transition-shadow duration-300 border-4 border-card dark:border-card-dark backdrop-blur-sm",
            styles.iconBg,
            styles.ring,
            { "boss-node": isBossNode }
          )}
        >
          <Icon
            className={clsx(
              "h-8 w-8 sm:h-9 sm:w-9 transition-transform group-hover:scale-110 duration-200",
              styles.iconColor
            )}
          />
          {node.state === "done" && (
            <div className="absolute -bottom-1 -right-1 bg-green-500 text-white rounded-full p-1 w-6 h-6 shadow-md animate-in zoom-in duration-300">
              <Check className="w-4 h-4" strokeWidth={3} />
            </div>
          )}
        </div>
      </div>

      <div className={clsx("flex-1 min-w-0", { "ml-1": isBossNode })}>
        <h3
          className={clsx(
            "font-bold text-sm sm:text-base leading-tight transition-colors duration-200 truncate",
            {
              "text-primary": node.state === "active",
              "group-hover:text-primary":
                isClickable && node.state !== "active",
            }
          )}
        >
          {t(node.title, node.id)}
        </h3>
        <div
          className={clsx(
            "mt-1 text-xs font-semibold px-2.5 py-1 rounded-full inline-flex items-center gap-1.5 shadow-sm transition-all duration-200",
            styles.badge,
            "group-hover:shadow-md"
          )}
        >
          {node.state === "locked" && <Lock className="w-3 h-3" />}
          <span className="capitalize">{node.state.replace("_", " ")}</span>
        </div>
      </div>
    </div>
  );
}
