import { useAuth } from "@/contexts/AuthContext";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { EditProfileForm } from "./EditProfileForm";
import { useState, useRef } from "react";
import { presignAvatarUpload, updateProfile } from "@/api/user";
import { useToast } from "@/components/ui/use-toast";
import { useQueryClient } from "@tanstack/react-query";
import { Loader2, Upload, AlertCircle } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { getPendingEmailVerification } from "@/api/user";

export default function ProfilePage() {
  const { user } = useAuth();
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  const { data: pendingEmail } = useQuery({
    queryKey: ["me", "pending-email"],
    queryFn: getPendingEmailVerification,
    refetchInterval: 30000,
  });

  if (!user) return null;

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 5 * 1024 * 1024) {
      toast({
        title: "File too large",
        description: "Max 5MB",
        variant: "destructive",
      });
      return;
    }

    setIsUploading(true);
    try {
      // 1. Presign
      const { upload_url, public_url } = await presignAvatarUpload(
        file.name,
        file.type,
        file.size
      );

      // 2. Upload to S3
      const uploadRes = await fetch(upload_url, {
        method: "PUT",
        body: file,
        headers: { "Content-Type": file.type },
      });

      if (!uploadRes.ok) throw new Error("Upload failed");

      // 3. Update profile
      await updateProfile({ avatar_url: public_url });

      // 4. Refresh
      await queryClient.invalidateQueries({ queryKey: ["me"] });
      toast({ title: "Avatar updated" });
    } catch (err) {
      console.error(err);
      toast({ title: "Failed to upload avatar", variant: "destructive" });
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="container max-w-4xl py-8 space-y-8">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-gradient">Profile</h1>
          <p className="text-muted-foreground">
            Manage your personal information
          </p>
        </div>
        <Button variant="outline" onClick={() => setIsEditing(!isEditing)}>
          {isEditing ? "Cancel Editing" : "Edit Profile"}
        </Button>
      </div>

      {pendingEmail?.pending && (
        <div className="p-3 rounded-md bg-blue-100 text-blue-800 text-sm flex items-start gap-2">
          <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
          <div>
            <strong>Email verification pending</strong>
            <p className="text-xs mt-1">
              Please check <strong>{pendingEmail.new_email}</strong> for the
              verification link.
            </p>
          </div>
        </div>
      )}

      <div className="grid gap-8 md:grid-cols-[300px_1fr]">
        <Card className="card-enhanced">
          <CardContent className="pt-6 flex flex-col items-center text-center space-y-4">
            <div className="relative group">
              <Avatar className="h-32 w-32">
                <AvatarImage src={user.avatar_url} />
                <AvatarFallback className="text-4xl">
                  {user.first_name?.[0]}
                  {user.last_name?.[0]}
                </AvatarFallback>
              </Avatar>
              <div
                className="absolute inset-0 bg-black/60 rounded-full opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center cursor-pointer"
                onClick={() => fileInputRef.current?.click()}
              >
                {isUploading ? (
                  <Loader2 className="h-8 w-8 animate-spin text-white" />
                ) : (
                  <Upload className="h-8 w-8 text-white" />
                )}
              </div>
              <input
                type="file"
                ref={fileInputRef}
                className="hidden"
                accept="image/*"
                onChange={handleFileChange}
              />
            </div>
            <div>
              <h2 className="text-xl font-semibold">
                {user.first_name} {user.last_name}
              </h2>
              <p className="text-sm text-muted-foreground">{user.role}</p>
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          {isEditing ? (
            <Card>
              <CardHeader>
                <CardTitle>Edit Profile</CardTitle>
              </CardHeader>
              <CardContent>
                <EditProfileForm
                  user={user as any} // Cast to match UserProfile if needed, or update types
                  onSuccess={() => {
                    setIsEditing(false);
                    queryClient.invalidateQueries({ queryKey: ["me"] });
                  }}
                />
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardHeader>
                <CardTitle>Personal Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">
                      Email
                    </label>
                    <p>{user.email}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">
                      Phone
                    </label>
                    <p>{user.phone || "-"}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">
                      Date of Birth
                    </label>
                    <p>
                      {user.date_of_birth
                        ? new Date(user.date_of_birth).toLocaleDateString()
                        : "-"}
                    </p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-muted-foreground">
                      Address
                    </label>
                    <p>{user.address || "-"}</p>
                  </div>
                </div>
                <div>
                  <label className="text-sm font-medium text-muted-foreground">
                    Bio
                  </label>
                  <p className="whitespace-pre-wrap">{user.bio || "-"}</p>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
