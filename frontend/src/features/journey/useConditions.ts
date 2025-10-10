import { useSubmission } from "./hooks";

export function useConditions() {
  const { submission } = useSubmission("S1_profile");
  const years = Number(submission?.form?.data?.years_since_graduation ?? 0);
  const rp_required = years > 3;
  return { rp_required } as const;
}

