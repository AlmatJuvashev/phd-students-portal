import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, AlertCircle, CheckCircle2 } from "lucide-react";

export const ProfilePage = () => {
  const queryClient = useQueryClient();
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const { data: me, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: () => api("/me"),
  });

  const { data: pendingEmail } = useQuery({
    queryKey: ["me", "pending-email"],
    queryFn: () => api("/me/pending-email"),
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  const updateMutation = useMutation({
    mutationFn: (data: { email: string; phone: string; current_password: string }) =>
      api.patch("/me", data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["me"] });
      queryClient.invalidateQueries({ queryKey: ["me", "pending-email"] });
      
      if (data.message === "verification_email_sent") {
        setSuccessMessage("Verification email sent! Please check your new email address to complete the change.");
      } else if (data.message === "verification_email_pending") {
        setSuccessMessage("Email change requested, but verification email could not be sent. Contact administrator.");
      } else if (data.message === "phone updated successfully") {
        setSuccessMessage("Phone number updated successfully!");
      } else {
        setSuccessMessage("Profile updated successfully!");
      }
      
      setErrorMessage("");
      setTimeout(() => setSuccessMessage(""), 5000);
    },
    onError: (err: any) => {
      const message = err.message || "Failed to update profile";
      setErrorMessage(message);
      setSuccessMessage("");
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const password = formData.get("current_password") as string;
    
    if (!password) {
      setErrorMessage("Current password is required");
      return;
    }
    
    updateMutation.mutate({
      email: formData.get("email") as string,
      phone: formData.get("phone") as string,
      current_password: password,
    });
    
    // Clear password field after submission
    e.currentTarget.reset();
    // Restore email and phone values
    setTimeout(() => {
      const emailInput = e.currentTarget.querySelector('[name="email"]') as HTMLInputElement;
      const phoneInput = e.currentTarget.querySelector('[name="phone"]') as HTMLInputElement;
      if (emailInput) emailInput.value = me?.email || "";
      if (phoneInput) phoneInput.value = me?.phone || "";
    }, 100);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto py-8">
      <Card>
        <CardHeader>
          <CardTitle>My Profile</CardTitle>
        </CardHeader>
        <CardContent>
          {pendingEmail?.pending && (
            <div className="mb-4 p-3 rounded-md bg-blue-100 text-blue-800 text-sm flex items-start gap-2">
              <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
              <div>
                <strong>Email verification pending</strong>
                <p className="text-xs mt-1">
                  Please check <strong>{pendingEmail.new_email}</strong> for the verification link.
                </p>
              </div>
            </div>
          )}
          
          {successMessage && (
            <div className="mb-4 p-3 rounded-md bg-green-100 text-green-800 text-sm flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 mt-0.5 flex-shrink-0" />
              <span>{successMessage}</span>
            </div>
          )}
          
          {errorMessage && (
            <div className="mb-4 p-3 rounded-md bg-red-100 text-red-800 text-sm flex items-start gap-2">
              <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
              <span>{errorMessage}</span>
            </div>
          )}
          
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="firstName">First Name</Label>
              <Input
                id="firstName"
                value={me?.first_name || ""}
                disabled
                className="bg-muted"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="lastName">Last Name</Label>
              <Input
                id="lastName"
                value={me?.last_name || ""}
                disabled
                className="bg-muted"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="email">Email (Required)</Label>
              <Input
                id="email"
                name="email"
                type="email"
                defaultValue={me?.email || ""}
                required
              />
              <p className="text-xs text-muted-foreground">
                This email will be used for all notifications. Changes require verification.
              </p>
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="phone">Phone</Label>
              <Input
                id="phone"
                name="phone"
                defaultValue={me?.phone || ""}
              />
            </div>
            
            <div className="border-t pt-4 mt-4">
              <div className="space-y-2">
                <Label htmlFor="current_password">
                  Current Password <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="current_password"
                  name="current_password"
                  type="password"
                  required
                  placeholder="Enter your password to confirm changes"
                />
                <p className="text-xs text-muted-foreground">
                  Required for security purposes
                </p>
              </div>
            </div>
            
            <div className="pt-2">
              <Button type="submit" disabled={updateMutation.isPending} className="w-full">
                {updateMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Save Changes
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

