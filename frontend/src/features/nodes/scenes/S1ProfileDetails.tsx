import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";

import {
  FormTaskDetails,
  FormTaskDetailsProps,
} from "@/features/nodes/details/variants/FormTaskDetails";
import { getProfileSnapshot } from "@/api/journey";
import { getPrograms, getSpecialties, getDepartments } from "@/api/dictionaries";

const PROFILE_QUERY_KEY = ["journey", "profile", "snapshot"] as const;

import { useAuth } from "@/contexts/AuthContext";

export default function S1ProfileDetails({
  node,
  initial,
  ...rest
}: Omit<FormTaskDetailsProps, "renderActions">) {
  const { user } = useAuth();
  
  // 1. Fetch Profile Snapshot
  const { data: snapshot } = useQuery({
    queryKey: PROFILE_QUERY_KEY,
    queryFn: getProfileSnapshot,
    staleTime: 5 * 60 * 1000,
  });

  // 2. Fetch Dictionaries
  const { data: programs } = useQuery({ queryKey: ["dicts", "programs"], queryFn: getPrograms, staleTime: 60 * 60 * 1000 });
  const { data: specialties } = useQuery({ queryKey: ["dicts", "specialties"], queryFn: getSpecialties, staleTime: 60 * 60 * 1000 });
  const { data: departments } = useQuery({ queryKey: ["dicts", "departments"], queryFn: getDepartments, staleTime: 60 * 60 * 1000 });

  // 3. Inject Options into Node Fields
  const enrichedNode = useMemo(() => {
    if (!node.requirements?.fields) return node;

    // Create a deep copy to avoid mutating the original playbook definition in memory
    const newNode = JSON.parse(JSON.stringify(node));

    newNode.requirements.fields = newNode.requirements.fields.map((f: any) => {
      if (f.key === "program" && programs) {
        return {
          ...f,
          options: programs.map(p => ({ value: p.name, label: p.name })) // Storing Name as value for now as per schema logic
        };
      }
      if (f.key === "specialty" && specialties) {
        return {
          ...f,
          options: specialties.map(s => ({ value: s.name, label: s.name }))
        };
      }
      if (f.key === "department" && departments) {
        return {
          ...f,
          options: departments.map(d => ({ value: d.name, label: d.name }))
        };
      }
      return f;
    });

    return newNode;
  }, [node, programs, specialties, departments]);


  const mergedInitial = useMemo(() => {
    // 1. Start with data from Auth User (lowest priority)
    const defaults = user ? {
        full_name: user.full_name || `${user.first_name || ''} ${user.last_name || ''}`.trim(),
        email: user.email,
        phone: user.phone,
        specialty: user.specialty,
        program: user.program,
        department: user.department,
    } : {};

    // 2. Override with Snapshot (submitted data)
    // 3. Override with "Initial" prop (node state)
    return { 
        ...defaults, 
        ...(snapshot ?? {}), 
        ...(initial ?? {}) 
    };
  }, [initial, snapshot, user]);

  return (
    <FormTaskDetails
      node={enrichedNode}
      initial={mergedInitial}
      {...rest}
    />
  );
}
