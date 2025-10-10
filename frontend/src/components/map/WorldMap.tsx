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
import clsx from "clsx";

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
    [playbook, stateByNodeId]
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
    0
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
          worldRefs.current[nextWorld.id]?.scrollIntoView({
            behavior: "smooth",
            block: "start",
          });
        }
      }
      prev[w.id] = isDone;
    });
  }, [vm.worlds]);

  return (
    <div className="p-4 space-y-6">
      <header className="sticky top-0 z-20 bg-gradient-to-b from-background via-background to-background/80 backdrop-blur-md shadow-lg -mx-4 -mt-4 px-4 py-4 rounded-b-xl">
        <div className="flex items-center justify-between mb-3">
          <div className="w-10"></div>
          <h1 className="text-xl font-bold text-center bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            {T("map.title", { defaultValue: "My Dissertation Map" })}
          </h1>
          <button className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-primary/10 transition-colors duration-200">
            {/* <Settings /> */}
          </button>
        </div>
        <div className="px-2">
          <div className="flex justify-between items-center bg-gradient-to-r from-muted/50 to-muted/30 p-3 rounded-xl shadow-sm">
            <span className="text-xs font-semibold text-muted-foreground ml-1">
              {T("map.progress", { defaultValue: "Progress" })}:
            </span>
            <div className="flex-grow bg-muted/30 rounded-full h-3 mx-3 overflow-hidden shadow-inner">
              <div
                className="bg-gradient-to-r from-primary to-primary/80 h-3 rounded-full transition-all duration-500 ease-out shadow-sm"
                style={{ width: `${progress}%` }}
              ></div>
            </div>
            <span className="text-sm font-bold text-primary mr-1">
              {progress}%
            </span>
          </div>
        </div>
      </header>

      {vm.worlds
        .filter((w) => (w.id === "W3" ? rp_required : true))
        .map((w, wi, arr) => {
          const worldDoneNodes = w.nodes.filter(
            (n) => n.state === "done"
          ).length;
          const worldProgressText = `${worldDoneNodes}/${w.nodes.length} ${T(
            "map.done_suffix",
            { defaultValue: "Done" }
          )}`;
          const isWorldDone = worldDoneNodes === w.nodes.length;
          const isWorldLocked =
            wi > 0 &&
            (() => {
              const prev = arr[wi - 1];
              const allDone = prev.nodes.every((n) => n.state === "done");
              const gatewayDone = prev.nodes.some(
                (n) => n.type === "gateway" && n.state === "done"
              );
              return !(allDone || gatewayDone);
            })();

          const isExpanded = !!expanded[w.id];

          return (
            <div
              key={w.id}
              ref={(el) => (worldRefs.current[w.id] = el)}
              className="animate-in fade-in slide-in-from-bottom-4 duration-500"
              style={{ animationDelay: `${wi * 100}ms` }}
            >
              <Card
                className={`rounded-2xl shadow-lg overflow-hidden bg-gradient-to-br from-card via-card to-card/50 border-2 transition-all duration-300 hover:shadow-xl ${
                  isWorldLocked
                    ? "opacity-50 grayscale"
                    : "hover:border-primary/20"
                }`}
              >
                {/* Header: clickable to expand/collapse when done */}
                <button
                  className={`w-full text-left p-5 border-b border-border/50 transition-all duration-200 ${
                    isWorldDone
                      ? "hover:bg-muted/30 cursor-pointer"
                      : "cursor-default"
                  }`}
                  onClick={() => {
                    if (isWorldDone)
                      setExpanded((e) => ({ ...e, [w.id]: !e[w.id] }));
                  }}
                  aria-expanded={isExpanded}
                >
                  <div className="flex items-center justify-between">
                    <h2
                      className={`text-lg sm:text-xl font-bold transition-colors duration-200 ${
                        isWorldLocked ? "text-muted-foreground" : "text-primary"
                      }`}
                    >
                      {w.title}
                    </h2>
                    <div className="flex items-center gap-2">
                      <div
                        className={`px-3 py-1.5 rounded-full text-xs font-semibold shadow-sm transition-all duration-200 ${
                          isWorldDone
                            ? "bg-green-500/20 text-green-700 dark:text-green-300 border border-green-500/30"
                            : "bg-yellow-500/20 text-yellow-700 dark:text-yellow-300 border border-yellow-500/30"
                        }`}
                      >
                        {worldProgressText}
                      </div>
                      {/* Indication chevron when completed (collapsible) */}
                      {isWorldDone && (
                        <ChevronDown
                          className={`transition-transform duration-300 text-muted-foreground ${
                            isExpanded ? "rotate-180" : "rotate-0"
                          }`}
                          aria-hidden="true"
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
                      transition={{ duration: 0.3, ease: "easeInOut" }}
                      className="overflow-hidden"
                    >
                      <div className="p-6 sm:p-8 relative bg-gradient-to-b from-background/50 via-muted/5 to-muted/10">
                        <div className="space-y-5 sm:space-y-6">
                          {w.nodes.map((n, idx) => (
                            <NodeToken
                              key={n.id}
                              node={n}
                              onClick={(node) => {
                                if (node.type === "gateway")
                                  setGatewayNode(node);
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
                <div className="relative h-20 flex items-center justify-center my-2">
                  {/* Animated dashed line */}
                  <div className="absolute left-1/2 -translate-x-1/2 top-0 bottom-0 w-1 overflow-hidden rounded-full">
                    <div 
                      className="w-full h-full bg-gradient-to-b from-primary/30 to-primary/20"
                      style={{
                        backgroundImage: "repeating-linear-gradient(to bottom, transparent, transparent 4px, rgba(0,0,0,0.1) 4px, rgba(0,0,0,0.1) 8px)",
                      }}
                    />
                  </div>
                  
                  {/* Arrow icon with better styling */}
                  <div className="z-10 bg-background border-2 border-primary/20 rounded-full p-2 shadow-lg">
                    <ArrowDown className="w-5 h-5 text-primary animate-bounce" style={{ animationDuration: "2s" }} />
                  </div>
                </div>
              )}
            </div>
          );
        })}

      <ConfettiBurst trigger={confetti} />

      {/* Congratulations modal */}
      <Modal open={showCongrats} onClose={() => setShowCongrats(false)}>
        <div className="space-y-4 p-2">
          <div className="text-center">
            <div className="mx-auto w-16 h-16 bg-gradient-to-br from-green-400 to-green-600 rounded-full flex items-center justify-center mb-4 shadow-lg animate-in zoom-in duration-500">
              <svg
                className="w-10 h-10 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2.5}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <h3 className="text-2xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent mb-2">
              {T("map.congrats_title", { defaultValue: "Congratulations!" })}
            </h3>
          </div>
          <p className="text-sm text-muted-foreground text-center leading-relaxed">
            {T("map.congrats_message", {
              defaultValue:
                "Congratulations, {{name}}! You have successfully completed your dissertation journey.",
              name:
                ((profile as any)?.form?.data?.full_name as string) ||
                T("map.student_fallback_name", { defaultValue: "Student" }),
            })}
          </p>
          <div className="flex justify-center pt-2">
            <Button
              onClick={() => setShowCongrats(false)}
              className="min-w-[120px] shadow-lg hover:shadow-xl transition-shadow duration-200"
            >
              {T("common.ok", { defaultValue: "OK" })}
            </Button>
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
