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
  const Icon = typeIcon[node.type];
  const styles = stateStyles[node.state];

  const isBossNode = node.type === "boss";

  return (
    <div
      role="button"
      onClick={() => onClick?.(node)}
      className={clsx("flex items-center gap-4 relative", styles.opacity)}
    >
      <div
        className={clsx("z-10 relative", { "transform scale-110": isBossNode })}
      >
        {isBossNode && node.state === "active" && (
          <div className="absolute -inset-1.5 rounded-full bg-yellow-400 animate-glow"></div>
        )}
        <div
          className={clsx(
            "relative w-16 h-16 rounded-full flex items-center justify-center shadow-lg border-4 border-card dark:border-card-dark",
            styles.iconBg,
            styles.ring,
            { "boss-node": isBossNode }
          )}
        >
          <Icon className={clsx("h-8 w-8", styles.iconColor)} />
          {node.state === "done" && (
            <Check className="absolute -bottom-1 -right-1 bg-green-500 text-white rounded-full p-0.5 w-5 h-5" />
          )}
        </div>
      </div>

      <div className={clsx({ "ml-1": isBossNode })}>
        <h3
          className={clsx("font-bold", {
            "text-primary": node.state === "active",
          })}
        >
          {t(node.title, node.id)}
        </h3>
        <div
          className={clsx(
            "text-xs font-semibold px-2 py-0.5 rounded-full inline-flex items-center gap-1",
            styles.badge
          )}
        >
          {node.state === "locked" && <Lock className="text-sm" />}
          {node.state.replace("_", " ")}
        </div>
      </div>
    </div>
  );
}
