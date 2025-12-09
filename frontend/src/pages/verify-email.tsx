import React, { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, CheckCircle2, XCircle } from "lucide-react";
import { api } from "@/api/client";

export const VerifyEmailPage = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const [message, setMessage] = useState("");

  useEffect(() => {
    const token = searchParams.get("token");
    
    if (!token) {
      setStatus("error");
      setMessage("Invalid verification link - token missing");
      return;
    }

    // Call verification endpoint
    api.get(`/me/verify-email?token=${token}`)
      .then((data) => {
        setStatus("success");
        setMessage(data.message || "Email verified successfully!");
        // Redirect to profile after 3 seconds
        setTimeout(() => navigate("/profile"), 3000);
      })
      .catch((err) => {
        setStatus("error");
        setMessage(err.message || "Failed to verify email - link may be expired or invalid");
      });
  }, [searchParams, navigate]);

  return (
    <div className="max-w-md mx-auto py-16">
      <Card>
        <CardHeader>
          <CardTitle className="text-center">Email Verification</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center gap-4 py-8">
            {status === "loading" && (
              <>
                <Loader2 className="h-12 w-12 animate-spin text-blue-500" />
                <p className="text-sm text-muted-foreground">Verifying your email...</p>
              </>
            )}
            
            {status === "success" && (
              <>
                <CheckCircle2 className="h-12 w-12 text-green-500" />
                <p className="text-sm font-medium text-green-800">{message}</p>
                <p className="text-xs text-muted-foreground">
                  Redirecting to profile...
                </p>
              </>
            )}
            
            {status === "error" && (
              <>
                <XCircle className="h-12 w-12 text-red-500" />
                <p className="text-sm font-medium text-red-800">{message}</p>
                <button
                  onClick={() => navigate("/profile")}
                  className="text-sm text-blue-600 hover:underline mt-2"
                >
                  Go to Profile
                </button>
              </>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
