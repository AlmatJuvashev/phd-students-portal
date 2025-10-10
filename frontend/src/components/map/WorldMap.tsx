// components/map/WorldMap.tsx
import { useEffect, useMemo, useRef, useState } from "react";
import { NodeVM, Playbook, toViewModel, edgesForWorld } from "@/lib/playbook";
import { NodeToken } from "./NodeToken";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { NodeDetailsSheet } from "../node-details/NodeDetailsSheet";
import { EdgeConnector } from "./EdgeConnector";
import { ArrowDown, ChevronDown } from "lucide-react";
import { GatewayModal } from "../node-details/GatewayModal";
import { api } from "@/api/client";
import { patchJourneyState } from "@/features/journey/session";
import { ConfettiBurst } from "@/features/journey/components/ConfettiBurst";
import { useConditions } from "@/features/journey/useConditions";
import { AnimatePresence, motion } from "framer-motion";
import { useTranslation } from "react-i18next";
import { useSubmission } from "@/features/journey/hooks";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";

type Pos = { x: number; y: number };
type Layout = Record<string, Pos>;

function layoutWorld(nodes: NodeVM[]): Layout {
  // Simple stacked layout: spread nodes vertically; you can replace with your Stitch coordinates if you have them.
  const L: Layout = {};
  nodes.forEach((n, i) => (L[n.id] = { x: 0, y: i * 100 }));
  return L;
}

export function WorldMap({
  playbook,
  locale = "ru",
  stateByNodeId = {},
  onStateChanged,
}: {
  playbook: Playbook;
  locale?: string;
  stateByNodeId?: Record<string, NodeVM["state"]>;
  onStateChanged?: () => void;
}) {
  const { t: T } = useTranslation("common");
  const vm = useMemo(
    () => toViewModel(playbook, stateByNodeId),
    [playbook, stateByNodeId],
  );
  const [openNode, setOpenNode] = useState<NodeVM | null>(null);
  const [gatewayNode, setGatewayNode] = useState<NodeVM | null>(null);
  const [confetti, setConfetti] = useState(false);
  const { rp_required } = useConditions();
  const worldRefs = useRef<Record<string, HTMLDivElement | null>>({});
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});
  const prevDoneRef = useRef<Record<string, boolean>>({});
  const [showCongrats, setShowCongrats] = useState(false);
  const { submission: profile } = useSubmission("S1_profile");
  const findNode = (id: string): NodeVM | null => {
    for (const world of vm.worlds) {
      const found = world.nodes.find((n) => n.id === id);
      if (found) return found;
    }
    return null;
  };

  const totalNodes = vm.worlds.reduce((acc, w) => acc + w.nodes.length, 0);
  const doneNodes = vm.worlds.reduce(
    (acc, w) => acc + w.nodes.filter((n) => n.state === "done").length,
    0,
  );
  const progress =
    totalNodes > 0 ? Math.round((doneNodes / totalNodes) * 100) : 0;

  // Show congratulations modal once when completed
  useEffect(() => {
    try {
      const key = "journey_congrats_shown";
      if (progress === 100 && !sessionStorage.getItem(key)) {
        setShowCongrats(true);
        sessionStorage.setItem(key, "1");
      }
    } catch {}
  }, [progress]);

  // Initialize/refresh expanded state per world
  useEffect(() => {
    const next: Record<string, boolean> = { ...expanded };
    vm.worlds.forEach((w) => {
      const worldDone = w.nodes.every((n) => n.state === "done");
      if (!(w.id in next)) {
        next[w.id] = !worldDone; // collapse when already done
      }
    });
    if (JSON.stringify(next) !== JSON.stringify(expanded)) {
      setExpanded(next);
    }
  }, [vm.worlds]);

  // Detect newly completed world -> collapse and scroll next into view
  useEffect(() => {
    const prev = prevDoneRef.current;
    vm.worlds.forEach((w, idx) => {
      const isDone = w.nodes.every((n) => n.state === "done");
      const wasDone = prev[w.id] || false;
      if (isDone && !wasDone) {
        // collapse
        setExpanded((e) => ({ ...e, [w.id]: false }));
        // scroll next world into view
        const nextWorld = vm.worlds[idx + 1];
        if (nextWorld && worldRefs.current[nextWorld.id]) {
          worldRefs.current[nextWorld.id]?.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      }
      prev[w.id] = isDone;
    });
  }, [vm.worlds]);

  return (
    <div className="p-4 space-y-4">
      <header className="sticky top-0 z-20 bg-background/80 backdrop-blur-sm shadow-sm -mx-4 -mt-4 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="w-10"></div>
          <h1 className="text-lg font-bold text-center">{T("map.title", { defaultValue: "My Dissertation Map" })}</h1>
          <button className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-primary/10">
            {/* <Settings /> */}
          </button>
        </div>
        <div className="px-4 pb-1 pt-2">
          <div className="flex justify-between items-center bg-gray-100 dark:bg-gray-800 p-1 rounded-lg">
            <span className="text-xs font-bold text-gray-500 dark:text-gray-400 ml-2">{T("map.progress", { defaultValue: "Progress" })}:</span>
            <div className="flex-grow bg-gray-200 dark:bg-gray-700 rounded-full h-2.5 mx-3">
              <div
                className="bg-primary h-2.5 rounded-full"
                style={{ width: `${progress}%` }}
              ></div>
            </div>
            <span className="text-sm font-bold text-primary mr-2">
              {progress}%
            </span>
          </div>
        </div>
      </header>

      {vm.worlds
        .filter((w) => (w.id === "W3" ? rp_required : true))
        .map((w, wi, arr) => {
          const worldDoneNodes = w.nodes.filter((n) => n.state === "done").length;
          const worldProgressText = `${worldDoneNodes}/${w.nodes.length} ${T("map.done_suffix", { defaultValue: "Done" })}`;
          const isWorldDone = worldDoneNodes === w.nodes.length;
          const isWorldLocked =
          wi > 0 &&
          (() => {
            const prev = arr[wi - 1];
            const allDone = prev.nodes.every((n) => n.state === "done");
            const gatewayDone = prev.nodes.some(
              (n) => n.type === "gateway" && n.state === "done",
            );
            return !(allDone || gatewayDone);
          })();

        const isExpanded = !!expanded[w.id];

        return (
          <div key={w.id} ref={(el) => (worldRefs.current[w.id] = el)}>
            <Card
              className={`rounded-xl shadow-md overflow-hidden ${
                isWorldLocked ? "opacity-50" : ""
              }`}
            >
              {/* Header: clickable to expand/collapse when done */}
              <button
                className="w-full text-left p-4 border-b border-gray-200 dark:border-gray-700"
                onClick={() => {
                  if (isWorldDone) setExpanded((e) => ({ ...e, [w.id]: !e[w.id] }));
                }}
                aria-expanded={isExpanded}
              >
                <div className="flex items-center justify-between">
                  <h2
                    className={`text-lg font-bold ${
                      isWorldLocked
                        ? "text-gray-500 dark:text-gray-400"
                        : "text-primary"
                    }`}
                  >
                    {w.title}
                  </h2>
                  <div className="flex items-center gap-2">
                    <div
                      className={`px-3 py-1 rounded-full text-xs font-semibold ${
                        isWorldDone
                          ? "bg-green-500/20 text-green-700 dark:text-green-300"
                          : "bg-yellow-500/20 text-yellow-700 dark:text-yellow-300"
                      }`}
                    >
                      {worldProgressText}
                    </div>
                    {/* Indication chevron when completed (collapsible) */}
                    {isWorldDone && (
                      <ChevronDown
                        className={`transition-transform ${isExpanded ? "rotate-180" : "rotate-0"}`}
                        aria-hidden
                        title={isExpanded ? "Collapse" : "Expand"}
                      />
                    )}
                  </div>
                </div>
              </button>
              <AnimatePresence initial={false}>
                {(!isWorldDone || isExpanded) && (
                  <motion.div
                    key={`${w.id}-body`}
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: "auto", opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    transition={{ duration: 0.2 }}
                    className="overflow-hidden"
                  >
                    <div className="p-6 relative">
                      <div className="absolute left-8 top-0 bottom-0 w-0.5 bg-gray-200 dark:bg-slate-700"></div>
                      <div className="space-y-8">
                        {w.nodes.map((n) => (
                          <NodeToken
                            key={n.id}
                            node={n}
                            onClick={(node) => {
                              if (node.type === "gateway") setGatewayNode(node);
                              else setOpenNode(node);
                            }}
                          />
                        ))}
                      </div>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </Card>
            {wi < arr.length - 1 && (
              <div className="relative h-16 flex items-center justify-center">
                <div
                  className="absolute left-1/2 -translate-x-1/2 top-0 bottom-0 w-0.5 bg-gray-200 dark:bg-slate-700"
                  style={{
                    backgroundImage:
                      "repeating-linear-gradient(to bottom, #cbd5e1, #cbd5e1 4px, transparent 4px, transparent 8px)",
                  }}
                ></div>
                <ArrowDown className="z-10 text-gray-400 dark:text-gray-500 bg-background p-1" />
              </div>
            )}
          </div>
        );
      })}

      <ConfettiBurst trigger={confetti} />

      {/* Congratulations modal */}
      <Modal open={showCongrats} onClose={() => setShowCongrats(false)}>
        <div className="space-y-3">
          <h3 className="text-lg font-semibold">{T("map.congrats_title", { defaultValue: "Congratulations!" })}</h3>
          <p className="text-sm text-muted-foreground">
            {T("map.congrats_message", {
              defaultValue: "Congratulations, {{name}}! You have successfully completed your dissertation journey.",
              name:
                (profile?.form?.data?.full_name as string) ||
                T("map.student_fallback_name", { defaultValue: "Student" }),
            })}
          </p>
          <div className="flex justify-end gap-2 pt-2">
            <Button onClick={() => setShowCongrats(false)}>{T("common.ok", { defaultValue: "OK" })}</Button>
          </div>
        </div>
      </Modal>

      <NodeDetailsSheet
        node={openNode}
        onOpenChange={(o) => !o && setOpenNode(null)}
        onStateRefresh={onStateChanged}
        onAdvance={(nextId) => {
          if (nextId) {
            const nextNode = findNode(nextId);
            if (nextNode) {
              setOpenNode(nextNode);
              return;
            }
          }
          setOpenNode(null);
        }}
      />

      <GatewayModal
        node={gatewayNode}
        open={!!gatewayNode}
        onClose={() => setGatewayNode(null)}
        onUnlock={async (node) => {
          try {
            await api("/journey/state", {
              method: "PUT",
              body: JSON.stringify({ node_id: node.id, state: "done" }),
            });
            // persist session progress immediately
            patchJourneyState({ [node.id]: "done" });
            setGatewayNode(null);
            setConfetti(true);
            setTimeout(() => setConfetti(false), 1200);
            onStateChanged?.();
          } catch (e) {
            console.error(e);
          }
        }}
      />
    </div>
  );
}
