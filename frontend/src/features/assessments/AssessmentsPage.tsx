import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus, Search, Archive, Trash2, FileText } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Badge } from '../../components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '../../components/ui/card';
import { getAssessments, deleteAssessment } from './api';
import { Link } from 'react-router-dom';
import { format } from 'date-fns';

export function AssessmentsPage() {
  const [search, setSearch] = useState('');
  const queryClient = useQueryClient();

  const { data: assessments, isLoading } = useQuery({
    queryKey: ['assessments', { search }],
    queryFn: () => getAssessments({ search }),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteAssessment,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assessments'] });
    },
  });

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this assessment?')) {
      await deleteMutation.mutateAsync(id);
    }
  };

  return (
    <div className="space-y-6 container mx-auto py-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Assessments</h1>
          <p className="text-muted-foreground">Manage quizzes, exams, and surveys.</p>
        </div>
        <Button asChild>
          <Link to="/admin/assessments/new/builder">
            <Plus className="mr-2 h-4 w-4" />
            Create Assessment
          </Link>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
             <div className="space-y-1">
                <CardTitle>All Assessments</CardTitle>
                <p className="text-sm text-muted-foreground">View and manage all assessments across courses.</p>
             </div>
             <div className="relative w-64">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search assessments..."
                  className="pl-8"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
             </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Title</TableHead>
                <TableHead>Course</TableHead>
                <TableHead>Grading</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8">
                    Loading assessments...
                  </TableCell>
                </TableRow>
              ) : assessments?.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                    No assessments found. Create one to get started.
                  </TableCell>
                </TableRow>
              ) : (
                assessments?.map((assessment) => (
                  <TableRow key={assessment.id}>
                    <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                             <FileText className="h-4 w-4 text-primary" />
                             {assessment.title}
                        </div>
                    </TableCell>
                    <TableCell>
                      {assessment.course_offering?.course?.code || 'N/A'} - {assessment.course_offering?.section || ''}
                    </TableCell>
                    <TableCell>
                      <Badge variant={assessment.grading_policy === 'AUTOMATIC' ? 'secondary' : 'outline'}>
                        {assessment.grading_policy}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {format(new Date(assessment.created_at), 'MMM d, yyyy')}
                    </TableCell>
                    <TableCell className="text-right space-x-2">
                      <Button variant="ghost" size="icon" title="Edit" asChild>
                        <Link to={`/admin/assessments/${assessment.id}/builder`}>
                             <FileText className="h-4 w-4" />
                        </Link>
                      </Button>
                      <Button variant="ghost" size="icon" title="Preview" asChild>
                        <Link to={`/admin/assessments/${assessment.id}/preview`}>
                             <Search className="h-4 w-4" />
                        </Link>
                      </Button>
                      <Button variant="ghost" size="icon" title="Archive">
                        <Archive className="h-4 w-4" />
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="icon" 
                        className="text-destructive hover:text-destructive"
                        onClick={() => handleDelete(assessment.id)}
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
