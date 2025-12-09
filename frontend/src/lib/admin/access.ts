import type { User } from "@/contexts/AuthContext";

export const isSuperAdmin = (u?: User | null) => u?.role === "superadmin";
export const isAdmin = (u?: User | null) => u?.role === "admin" || u?.role === "superadmin";
export const requireAny = (u: User | null, roles: string[]) => !!u && roles.includes(u.role);

