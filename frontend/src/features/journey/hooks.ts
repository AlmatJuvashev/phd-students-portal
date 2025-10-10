import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api/client";
import { getNodeSubmission, saveNodeSubmission } from "@/api/journey";

export function useJourneyState() {
  const qc = useQueryClient();
  const query = useQuery({
    queryKey: ["journey", "state"],
    queryFn: () => api("/journey/state"),
  });
  const reset = useMutation({
    mutationFn: async () => api("/journey/reset", { method: "POST" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["journey", "state"] }),
  });

  const setNodeState = useMutation({
    mutationFn: async ({ node_id, state }: { node_id: string; state: string }) =>
      api("/journey/state", {
        method: "PUT",
        body: JSON.stringify({ node_id, state }),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["journey", "state"] }),
  });

  return {
    state: query.data as Record<string, string> | undefined,
    isLoading: query.isLoading,
    refetch: query.refetch,
    reset: () => reset.mutate(),
    setNodeState: (node_id: string, state: string) => setNodeState.mutate({ node_id, state }),
  };
}

export function useSubmission(nodeId?: string | null) {
  const enabled = !!nodeId;
  const query = useQuery({
    queryKey: ["journey", "node", nodeId, "submission"],
    queryFn: () => getNodeSubmission(nodeId!),
    enabled,
  });
  const qc = useQueryClient();
  const save = useMutation({
    mutationFn: async (payload: { form_data?: any; state?: string }) =>
      saveNodeSubmission(nodeId!, payload),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["journey", "node", nodeId, "submission"] }),
  });
  return { submission: query.data, isLoading: query.isLoading, save };
}

export function useAdvance(playbook: any) {
  const nextIdOf = (node: { next?: string[] } | null) =>
    Array.isArray(node?.next) && node!.next.length > 0 ? node!.next[0] : null;
  return { nextIdOf };
}

