import { api } from "./client";

export interface UserProfile {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  bio?: string;
  address?: string;
  date_of_birth?: string; // ISO string
  avatar_url?: string;
  program?: string;
  specialty?: string;
  department?: string;
  cohort?: string;
  role: string;
  is_active: boolean;
}

export async function updateProfile(data: Partial<UserProfile> & { current_password?: string }) {
  return api.patch("/me", data);
}

export async function updateAvatar(url: string) {
  return api.patch("/me/avatar", { avatar_url: url });
}

export async function presignAvatarUpload(filename: string, contentType: string, sizeBytes: number) {
  return api<{ upload_url: string; object_key: string; public_url: string }>("/me/avatar/presign", {
    method: "POST",
    body: JSON.stringify({
      filename,
      content_type: contentType,
      size_bytes: sizeBytes,
    }),
  });
}
export async function getPendingEmailVerification() {
  return api<{ pending: boolean; new_email?: string }>("/me/pending-email");
}
