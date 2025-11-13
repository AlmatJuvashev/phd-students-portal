import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { fetchStudentDetails } from "../api";
import { StudentJourneyPanel } from "../components/StudentJourneyPanel";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { stageLabel } from "../utils";

export function StudentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading, error } = useQuery({
    queryKey: ["student", id],
    queryFn: () => fetchStudentDetails(id!),
    enabled: !!id,
  });

  return (
    <div className="max-w-5xl mx-auto space-y-6 p-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" onClick={() => navigate(-1)}>
          Back
        </Button>
        <div>
          <h1 className="text-2xl font-semibold">Student Information</h1>
          <p className="text-sm text-muted-foreground">
            Detailed journey overview and document status
          </p>
        </div>
      </div>
      {isLoading && <div className="text-sm text-muted-foreground">Loading…</div>}
      {error && (
        <div className="text-sm text-red-600">Failed to load student details.</div>
      )}
      {data && (
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>{data.name}</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-3 md:grid-cols-2">
              <div>
                <div className="text-xs text-muted-foreground">Program</div>
                <p className="font-medium">{data.program || "—"}</p>
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Department</div>
                <p className="font-medium">{data.department || "—"}</p>
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Cohort</div>
                <p className="font-medium">{data.cohort || "—"}</p>
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Advisors</div>
                <div className="flex flex-wrap gap-2">
                  {(data.advisors || []).map((a) => (
                    <Badge key={a.id} variant="outline">
                      {a.name}
                    </Badge>
                  ))}
                </div>
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Current stage</div>
                <Badge className="bg-primary/10 text-primary">
                  {stageLabel(data.current_stage)}
                </Badge>
              </div>
              <div>
                <div className="text-xs text-muted-foreground">Progress</div>
                <p className="font-semibold">
                  {Math.round(data.overall_progress_pct || 0)}%
                </p>
              </div>
            </CardContent>
          </Card>
          <StudentJourneyPanel studentId={data.id} />
        </div>
      )}
    </div>
  );
}

export default StudentDetailPage;
