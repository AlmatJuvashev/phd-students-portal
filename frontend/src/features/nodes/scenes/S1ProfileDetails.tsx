import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";

import {
  FormTaskDetails,
  FormTaskDetailsProps,
} from "@/features/nodes/details/variants/FormTaskDetails";
import { getProfileSnapshot } from "@/api/journey";

const PROFILE_QUERY_KEY = ["journey", "profile", "snapshot"] as const;

export default function S1ProfileDetails({
  node,
  initial,
  ...rest
}: Omit<FormTaskDetailsProps, "renderActions">) {
  const { data } = useQuery({
    queryKey: PROFILE_QUERY_KEY,
    queryFn: getProfileSnapshot,
    staleTime: 5 * 60 * 1000,
  });

  const mergedInitial = useMemo(() => {
    if (!data) return initial;
    return { ...(initial ?? {}), ...data };
  }, [initial, data]);

  return (
    <FormTaskDetails
      node={node}
      initial={mergedInitial}
      {...rest}
    />
  );
}
