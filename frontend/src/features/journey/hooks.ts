import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api/client";
import { getNodeSubmission, saveNodeSubmission } from "@/api/journey";
import {
  loadJourneyState,
  saveJourneyState,
  patchJourneyState,
} from "./session";

export function useJourneyState() {
  const qc = useQueryClient();
  const query = useQuery({
    queryKey: ["journey", "state"],
    queryFn: async () => {
      const res = await api("/journey/state");
      // Persist in session for offline fallback
      if (res && typeof res === "object") saveJourneyState(res);
      return res;
    },
    initialData: loadJourneyState() || undefined,
    retry: 0,
    staleTime: 30 * 1000, // 30 seconds - shorter cache to catch state changes faster
    refetchOnWindowFocus: true, // Refetch when user returns to tab
  });
  const reset = useMutation({
    mutationFn: async () => api("/journey/reset", { method: "POST" }),
    onSuccess: () => {
      // Clear sessionStorage journey state
      saveJourneyState({});
      // Invalidate all journey-related queries
      qc.invalidateQueries({ queryKey: ["journey"] });
      // Force reload to ensure clean state
      window.location.reload();
    },
  });

  const setNodeState = useMutation({
    mutationFn: async ({
      node_id,
      state,
    }: {
      node_id: string;
      state: string;
    }) =>
      api("/journey/state", {
        method: "PUT",
        body: JSON.stringify({ node_id, state }),
      }),
    onSuccess: (_data, vars) => {
      // Update session copy as well
      patchJourneyState({ [vars.node_id]: vars.state });
      qc.invalidateQueries({ queryKey: ["journey", "state"] });
    },
  });

  return {
    state: query.data as Record<string, string> | undefined,
    isLoading: query.isLoading,
    refetch: query.refetch,
    reset: () => reset.mutate(),
    setNodeState: (node_id: string, state: string) =>
      setNodeState.mutate({ node_id, state }),
  };
}

export function useSubmission(nodeId?: string | null) {
  const enabled = !!nodeId;
  const qc = useQueryClient();
  const query = useQuery({
    queryKey: ["journey", "node", nodeId, "submission"],
    queryFn: async () => {
      const result = await getNodeSubmission(nodeId!);
      // Invalidate journey state when fetching submission to ensure node states are fresh
      qc.invalidateQueries({ queryKey: ["journey", "state"] });
      return result;
    },
    enabled,
    staleTime: 30 * 1000, // 30 seconds - match journey state cache
    placeholderData: (previousData) => previousData,
    refetchOnWindowFocus: true, // Refetch when returning to tab
  });
  const save = useMutation({
    mutationFn: async (payload: { form_data?: any; state?: string }) =>
      saveNodeSubmission(nodeId!, payload),
    onSuccess: () => {
      qc.invalidateQueries({
        queryKey: ["journey", "node", nodeId, "submission"],
      });
      qc.invalidateQueries({ queryKey: ["journey", "state"] });
      if (nodeId === "S1_profile") {
        qc.invalidateQueries({ queryKey: ["journey", "profile", "snapshot"] });
      }
    },
  });
  return {
    submission: query.data,
    isLoading: query.isLoading,
    save,
    refetch: query.refetch,
  };
}

export function useAdvance(playbook: any) {
  const nextIdOf = (node: { next?: string[] } | null) =>
    Array.isArray(node?.next) && node!.next.length > 0 ? node!.next[0] : null;
  return { nextIdOf };
}
