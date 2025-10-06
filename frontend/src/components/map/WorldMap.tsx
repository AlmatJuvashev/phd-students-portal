// components/map/WorldMap.tsx
import { useMemo, useState } from "react";
import { NodeVM, Playbook, toViewModel, edgesForWorld } from "@/lib/playbook";
import { NodeToken } from "./NodeToken";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { NodeDetailsSheet } from "../node-details/NodeDetailsSheet";
import { EdgeConnector } from "./EdgeConnector";
import { ArrowDown } from "lucide-react";
import { GatewayModal } from "../node-details/GatewayModal";
import { api } from "@/api/client";

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
  const vm = useMemo(
    () => toViewModel(playbook, stateByNodeId),
    [playbook, stateByNodeId]
  );
  const [openNode, setOpenNode] = useState<NodeVM | null>(null);
  const [gatewayNode, setGatewayNode] = useState<NodeVM | null>(null);

  const totalNodes = vm.worlds.reduce((acc, w) => acc + w.nodes.length, 0);
  const doneNodes = vm.worlds.reduce(
    (acc, w) => acc + w.nodes.filter((n) => n.state === "done").length,
    0
  );
  const progress =
    totalNodes > 0 ? Math.round((doneNodes / totalNodes) * 100) : 0;

  return (
    <div className="p-4 space-y-4">
      <header className="sticky top-0 z-20 bg-background/80 backdrop-blur-sm shadow-sm -mx-4 -mt-4 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="w-10"></div>
          <h1 className="text-lg font-bold text-center">My Journey</h1>
          <button className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-primary/10">
            {/* <Settings /> */}
          </button>
        </div>
        <div className="px-4 pb-1 pt-2">
          <div className="flex justify-between items-center bg-gray-100 dark:bg-gray-800 p-1 rounded-lg">
            <span className="text-xs font-bold text-gray-500 dark:text-gray-400 ml-2">
              Progress:
            </span>
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

      {vm.worlds.map((w, wi) => {
        const worldDoneNodes = w.nodes.filter((n) => n.state === "done").length;
        const worldProgressText = `${worldDoneNodes}/${w.nodes.length} Done`;
        const isWorldDone = worldDoneNodes === w.nodes.length;
        const isWorldLocked = wi > 0 && (() => {
          const prev = vm.worlds[wi - 1];
          const allDone = prev.nodes.every((n) => n.state === "done");
          const gatewayDone = prev.nodes.some((n) => n.type === "gateway" && n.state === "done");
          return !(allDone || gatewayDone);
        })();

        return (
          <div key={w.id}>
            <Card
              className={`rounded-xl shadow-md overflow-hidden ${
                isWorldLocked ? "opacity-50" : ""
              }`}
            >
              <div className="p-4 border-b border-gray-200 dark:border-gray-700">
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
                  <div
                    className={`px-3 py-1 rounded-full text-xs font-semibold ${
                      isWorldDone
                        ? "bg-green-500/20 text-green-700 dark:text-green-300"
                        : "bg-yellow-500/20 text-yellow-700 dark:text-yellow-300"
                    }`}
                  >
                    {worldProgressText}
                  </div>
                </div>
              </div>
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
            </Card>
            {wi < vm.worlds.length - 1 && (
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

      <NodeDetailsSheet
        node={openNode}
        onOpenChange={(o) => !o && setOpenNode(null)}
        onStateRefresh={onStateChanged}
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
            setGatewayNode(null);
            onStateChanged?.();
          } catch (e) {
            console.error(e);
          }
        }}
      />
    </div>
  );
}
