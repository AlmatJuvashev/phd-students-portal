// components/map/WorldMap.tsx
import { useEffect, useMemo, useRef, useState } from "react";
import { NodeVM, Playbook, toViewModel, edgesForWorld } from "@/lib/playbook";
import { NodeToken } from "./NodeToken";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { BackButton } from "@/components/ui/back-button";
import { NodeDetailsSheet } from "@/features/nodes/details/NodeDetailsSheet";
import { EdgeConnector } from "./EdgeConnector";
import { ArrowDown, ChevronDown } from "lucide-react";
import { GatewayModal } from "@/features/nodes/details/GatewayModal";
import ModuleGuardModal from "@/features/nodes/details/ModuleGuardModal";
import { api } from "@/api/client";
import { patchJourneyState } from "@/features/journey/session";
import { ConfettiBurst } from "@/features/journey/components/ConfettiBurst";
import { useConditions } from "@/features/journey/useConditions";
import { AnimatePresence, motion } from "framer-motion";
import { useTranslation } from "react-i18next";
import { useSubmission } from "@/features/journey/hooks";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import DevBar from "@/features/journey/components/DevBar";
import clsx from "clsx";
import { detectTerminalNodeIds } from "@/features/journey/moduleGraph";
import { useSwipeable } from "react-swipeable";
import { useAuth } from '@/contexts/AuthContext'

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
  const { user } = useAuth()
  const { t: T } = useTranslation("common");
  const [unlockAll, setUnlockAll] = useState(() => {
    try {
      return (
        (import.meta as any).env?.VITE_UNLOCK_ALL_NODES === "true" ||
        localStorage.getItem("dev_unlock_all_nodes") === "true"
      );
    } catch {
      return (import.meta as any).env?.VITE_UNLOCK_ALL_NODES === "true";
    }
  });

  // Listen to storage changes for unlock all
  useEffect(() => {
    const handleUnlockChange = () => {
      try {
        const isUnlocked =
          (import.meta as any).env?.VITE_UNLOCK_ALL_NODES === "true" ||
          localStorage.getItem("dev_unlock_all_nodes") === "true";
        setUnlockAll(isUnlocked);
      } catch {}
    };
    window.addEventListener("storage", handleUnlockChange);
    window.addEventListener("dev_unlock_changed", handleUnlockChange);
    return () => {
      window.removeEventListener("storage", handleUnlockChange);
      window.removeEventListener("dev_unlock_changed", handleUnlockChange);
    };
  }, []);

  // Stabilize stateByNodeId to prevent unnecessary re-renders
  // Use a ref to track previous state and only update when actually different
  const prevStateRef = useRef<Record<string, NodeVM["state"]> | undefined>();
  const stableStateByNodeId = useMemo(() => {
    const current = stateByNodeId || {};
    const prev = prevStateRef.current || {};

    // Check if states actually changed
    const keys = new Set([...Object.keys(current), ...Object.keys(prev)]);
    let hasChanged = false;
    for (const key of keys) {
      if (current[key] !== prev[key]) {
        hasChanged = true;
        break;
      }
    }

    if (hasChanged) {
      prevStateRef.current = current;
      return current;
    }
    return prev;
  }, [stateByNodeId]);

  const stateOverride = useMemo(() => {
    if (!unlockAll) return stableStateByNodeId;
    const m: Record<string, NodeVM["state"]> = {};
    for (const w of playbook.worlds) {
      for (const n of w.nodes) m[n.id] = "active";
    }
    return m;
  }, [unlockAll, playbook.worlds, stableStateByNodeId]);

  const { rp_required } = useConditions();

  const activeConditions = useMemo(() => {
    const conditions = new Set<string>();
    if (rp_required) {
      conditions.add("rp_required");
    }
    return conditions;
  }, [rp_required]);

  const vm = useMemo(
    () => toViewModel(playbook, stateOverride, activeConditions),
    [playbook, stateOverride, activeConditions]
  );
  const [openNode, setOpenNode] = useState<NodeVM | null>(null);
  const [gatewayNode, setGatewayNode] = useState<NodeVM | null>(null);
  const [confetti, setConfetti] = useState(false);
  const [pendingModule, setPendingModule] = useState<{
    fromWorldId: string;
    toWorldId: string;
    toFirstNodeId: string;
  } | null>(null);
  const worldRefs = useRef<Record<string, HTMLDivElement | null>>({});
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});
  const prevDoneRef = useRef<Record<string, boolean>>({});
  const [showCongrats, setShowCongrats] = useState(false);
  const { submission: profile } = useSubmission("S1_profile");
  const [focusedWorldIndex, setFocusedWorldIndex] = useState<number | null>(
    null
  );

  const terminalNodeSet = useMemo(() => {
    const perWorld = detectTerminalNodeIds(playbook);
    const set = new Set<string>();
    Object.values(perWorld).forEach((ids) => ids.forEach((id) => set.add(id)));
    return set;
  }, [playbook]);
  const findNode = (id: string): NodeVM | null => {
    for (const world of vm.worlds) {
      const found = world.nodes.find((n) => n.id === id);
      if (found) return found;
    }
    return null;
  };

  // Filter visible worlds based on unlock criteria
  const visibleWorlds = useMemo(
    () =>
      vm.worlds.filter((w) =>
        w.id === "W3" ? unlockAll || rp_required : true
      ),
    [vm.worlds, unlockAll, rp_required]
  );

  const progress = useMemo(() => {
    const totals = visibleWorlds.reduce(
      (acc, w) => {
        acc.total += w.nodes.length;
        acc.done += w.nodes.filter((n) => n.state === "done").length;
        return acc;
      },
      { done: 0, total: 0 }
    );
    return totals.total > 0
      ? Math.round((totals.done / totals.total) * 100)
      : 0;
  }, [visibleWorlds]);

  const worldSignature = useMemo(
    () =>
      vm.worlds
        .map(
          (w) => `${w.id}:${w.nodes.map((n) => `${n.id}:${n.state}`).join("|")}`
        )
        .join(";"),
    [vm.worlds]
  );

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

  // Initialize expanded state per world only once
  useEffect(() => {
    setExpanded((prev) => {
      const next: Record<string, boolean> = {};
      let changed = false;
      vm.worlds.forEach((w) => {
        const isDone = w.nodes.every((n) => n.state === "done");
        const current = w.id in prev ? prev[w.id] : !isDone;
        if (current !== prev[w.id]) changed = true;
        next[w.id] = current;
      });
      if (Object.keys(prev).length !== Object.keys(next).length) changed = true;
      return changed ? next : prev;
    });
  }, [worldSignature, vm.worlds]);

  // Detect newly completed world -> collapse and scroll next into view
  useEffect(() => {
    const prev = prevDoneRef.current;
    vm.worlds.forEach((w, idx) => {
      const isDone = w.nodes.every((n) => n.state === "done");
      const wasDone = prev[w.id] || false;

      if (isDone && !wasDone) {
        // Check if we already showed the guard modal for this module
        let alreadyShown = false;
        try {
          alreadyShown = !!sessionStorage.getItem(`module_guard_${w.id}_shown`);
        } catch {}

        if (!alreadyShown) {
          // collapse
          setExpanded((e) => ({ ...e, [w.id]: false }));
          // scroll next world into view
          const nextWorld = vm.worlds[idx + 1];
          if (nextWorld && worldRefs.current[nextWorld.id]) {
            worldRefs.current[nextWorld.id]?.scrollIntoView({
              behavior: "smooth",
              block: "start",
            });
            // queue guard modal for next module, without auto-opening its first node
            const firstNode = nextWorld.nodes[0];
            if (firstNode) {
              setPendingModule({
                fromWorldId: w.id,
                toWorldId: nextWorld.id,
                toFirstNodeId: firstNode.id,
              });
              try {
                sessionStorage.setItem(`module_guard_${w.id}_shown`, "1");
              } catch {}
            }
          }
        }
      }
      prev[w.id] = isDone;
    });
  }, [vm.worlds, worldSignature]);

  // Swipe handlers for mobile navigation between modules
  const handleSwipe = (direction: "left" | "right") => {
    if (focusedWorldIndex === null) return;

    const nextIndex =
      direction === "left"
        ? Math.min(focusedWorldIndex + 1, visibleWorlds.length - 1)
        : Math.max(focusedWorldIndex - 1, 0);

    if (nextIndex !== focusedWorldIndex) {
      setFocusedWorldIndex(nextIndex);
      const targetWorld = visibleWorlds[nextIndex];
      if (targetWorld && worldRefs.current[targetWorld.id]) {
        worldRefs.current[targetWorld.id]?.scrollIntoView({
          behavior: "smooth",
          block: "center",
        });
      }
    }
  };

  const swipeHandlers = useSwipeable({
    onSwipedLeft: () => handleSwipe("left"),
    onSwipedRight: () => handleSwipe("right"),
    trackMouse: false, // Only track touch, not mouse
    preventScrollOnSwipe: false,
    delta: 50, // Minimum distance for swipe
  });

  // Track which world is in viewport for swipe context
  useEffect(() => {
    if (typeof window === "undefined") return;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && entry.intersectionRatio > 0.5) {
            const worldId = entry.target.getAttribute("data-world-id");
            const index = visibleWorlds.findIndex((w) => w.id === worldId);
            if (index >= 0) {
              setFocusedWorldIndex(index);
            }
          }
        });
      },
      { threshold: [0.5], rootMargin: "-20% 0px -20% 0px" }
    );

    Object.values(worldRefs.current).forEach((ref) => {
      if (ref) observer.observe(ref);
    });

    return () => observer.disconnect();
  }, [visibleWorlds]);

  return (
    <div className="p-4 space-y-6 max-w-4xl mx-auto" {...swipeHandlers}>
      <header className="sticky top-0 z-20 bg-gradient-to-b from-background via-background to-background/80 backdrop-blur-md shadow-lg -mx-4 -mt-4 px-4 py-4 rounded-b-xl">
        <div className="flex items-center justify-between mb-3 gap-2">
          <BackButton to="/" />
          <h1 className="text-lg sm:text-xl font-bold text-center bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent flex-1">
            {T("map.title", { defaultValue: "My Dissertation Map" })}
          </h1>
          <div className="w-[78px] sm:w-10"></div>
        </div>
        <div className="px-2 space-y-3">
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

          {/* Module navigation dots (mobile) */}
          {visibleWorlds.length > 1 && (
            <div className="flex justify-center items-center gap-2 sm:hidden py-2">
              {visibleWorlds.map((w, idx) => {
                const isActive = focusedWorldIndex === idx;
                const isDone = w.nodes.every((n) => n.state === "done");
                return (
                  <button
                    key={w.id}
                    onClick={() => {
                      setFocusedWorldIndex(idx);
                      worldRefs.current[w.id]?.scrollIntoView({
                        behavior: "smooth",
                        block: "center",
                      });
                    }}
                    aria-label={`${T("map.navigate_to_module", {
                      defaultValue: "Navigate to module",
                    })} ${idx + 1}`}
                    className={clsx(
                      "transition-all duration-300 rounded-full touch-manipulation",
                      {
                        "w-8 h-2 bg-primary shadow-lg": isActive,
                        "w-2 h-2": !isActive,
                        "bg-green-500": !isActive && isDone,
                        "bg-muted-foreground/30": !isActive && !isDone,
                      }
                    )}
                  />
                );
              })}
            </div>
          )}
        </div>
      </header>

      {visibleWorlds.map((w, wi, arr) => {
        const worldDoneNodes = w.nodes.filter((n) => n.state === "done").length;
        const worldProgressText = `${worldDoneNodes}/${w.nodes.length} ${T(
          "map.done_suffix",
          { defaultValue: "Done" }
        )}`;
        const isWorldDone = worldDoneNodes === w.nodes.length;
        const isWorldLocked = (() => {
          // RP world is visible/unlocked when required or forced
          if (w.id === "W3" && (rp_required || unlockAll)) return false;
          if (wi === 0) return false;
          const prev = arr[wi - 1];
          const allDone = prev.nodes.every((n) => n.state === "done");
          const gatewayDone = prev.nodes.some(
            (n) => n.type === "gateway" && n.state === "done"
          );
          if (allDone || gatewayDone) return false;
          // Edge-based unlock: if any done node in prev points to a node in this world
          const currentIds = new Set(w.nodes.map((n) => n.id));
          const edgeUnlock = prev.nodes.some((pn) => {
            if (pn.state !== "done") return false;
            const next = Array.isArray(pn.next) ? pn.next : [];
            return next.some((nid: string) => currentIds.has(nid));
          });
          return !edgeUnlock;
        })();

        const isExpanded = !!expanded[w.id];

        const isFocused = focusedWorldIndex === wi;

        return (
          <div
            key={w.id}
            ref={(el) => (worldRefs.current[w.id] = el)}
            data-world-id={w.id}
            className={clsx(
              "animate-in fade-in slide-in-from-bottom-4 duration-500 transition-all",
              {
                "scale-[1.01] sm:scale-100": isFocused, // Subtle scale on mobile when focused
              }
            )}
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
              <div className="relative h-20 flex items-center justify-center my-2">
                {/* Animated dashed line */}
                <div className="absolute left-1/2 -translate-x-1/2 top-0 bottom-0 w-1 overflow-hidden rounded-full">
                  <div
                    className="w-full h-full bg-gradient-to-b from-primary/30 to-primary/20"
                    style={{
                      backgroundImage:
                        "repeating-linear-gradient(to bottom, transparent, transparent 4px, rgba(0,0,0,0.1) 4px, rgba(0,0,0,0.1) 8px)",
                    }}
                  />
                </div>

                {/* Arrow icon with better styling */}
                <div className="z-10 bg-background border-2 border-primary/20 rounded-full p-2 shadow-lg">
                  <ArrowDown
                    className="w-5 h-5 text-primary animate-bounce"
                    style={{ animationDuration: "2s" }}
                  />
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
        role={(user?.role as any) || 'student'}
        onStateRefresh={onStateChanged}
        closeOnComplete={openNode ? terminalNodeSet.has(openNode.id) : false}
        onAdvance={(nextId, currentId) => {
          if (currentId && terminalNodeSet.has(currentId)) {
            setOpenNode(null);
            return;
          }
          if (nextId) {
            const nextNode = findNode(nextId);
            if (nextNode) {
              // Check if next node is in a different world/module than current
              const currentNode = currentId ? findNode(currentId) : null;
              if (currentNode && currentNode.worldId !== nextNode.worldId) {
                // Cross-module transition: close sheet, let ModuleGuardModal handle it
                setOpenNode(null);
                return;
              }
              // Same module: keep sheet open and animate content switch
              setOpenNode(nextNode);
              return;
            }
          }
          setOpenNode(null);
        }}
      />

      <ModuleGuardModal
        open={!!pendingModule}
        title={T("module.unlock_title", { defaultValue: "Unlock Next Module" })}
        onClose={() => setPendingModule(null)}
        onConfirm={() => {
          if (!pendingModule) return;
          setPendingModule(null);
          // Expand and center next world, do not open its first node
          setExpanded((e) => ({ ...e, [pendingModule.toWorldId]: true }));
          const el = worldRefs.current[pendingModule.toWorldId];
          el?.scrollIntoView({ behavior: "smooth", block: "start" });
          setConfetti(true);
          setTimeout(() => setConfetti(false), 1200);
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

      {(import.meta as any).env?.DEV && <DevBar />}
    </div>
  );
}
