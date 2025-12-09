import { useQuery } from "@tanstack/react-query";
import { getProfileSnapshot } from "@/api/journey";

export function useProfileSnapshot(enabled = true) {
  return useQuery({
    queryKey: ["journey", "profile", "snapshot"],
    queryFn: () => getProfileSnapshot(),
    enabled,
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
  });
}
