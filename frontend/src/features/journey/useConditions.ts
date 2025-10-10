import { useSubmission } from "./hooks";

export function useConditions() {
  const { submission } = useSubmission("S1_profile");
  // Prefer date-based calculation
  const grad = submission?.form?.data?.graduation_date as string | undefined;
  let years = 0;
  if (grad) {
    const d = new Date(grad);
    if (!isNaN(d.getTime())) {
      const now = new Date();
      const diff = now.getTime() - d.getTime();
      years = diff / (1000 * 60 * 60 * 24 * 365.25);
    }
  } else {
    years = Number(submission?.form?.data?.years_since_graduation ?? 0);
  }
  const rp_required = years > 3;
  return { rp_required } as const;
}
