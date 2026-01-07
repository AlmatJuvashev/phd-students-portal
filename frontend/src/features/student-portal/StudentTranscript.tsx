import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Loader2, GraduationCap, Award, BookOpen, AlertCircle } from 'lucide-react';
import { getStudentTranscript } from './api';
import { format } from 'date-fns';

export function StudentTranscript() {
  const { t } = useTranslation('common');
  
  const { data: transcript, isLoading, error } = useQuery({
    queryKey: ['student', 'transcript'],
    queryFn: getStudentTranscript,
  });

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-full flex-col items-center justify-center gap-4 p-8 text-center">
        <AlertCircle className="h-10 w-10 text-destructive" />
        <p className="text-lg font-semibold">{t('transcript.error_loading', 'Failed to load transcript')}</p>
        <p className="text-sm text-muted-foreground">
          {t('transcript.try_later', 'Please try again later or contact support.')}
        </p>
      </div>
    );
  }

  if (!transcript) return null;

  return (
    <div className="space-y-6 container mx-auto max-w-5xl py-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold tracking-tight">{t('transcript.title', 'Academic Transcript')}</h1>
        <p className="text-muted-foreground">
          {t('transcript.subtitle', 'Your complete academic record and GPA summary.')}
        </p>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('transcript.cumulative_gpa', 'Cumulative GPA')}
            </CardTitle>
            <Award className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{transcript.cumulative_gpa.toFixed(2)}</div>
            <p className="text-xs text-muted-foreground">
              {t('transcript.total_points', 'Total Quality Points')}: {transcript.total_points.toFixed(1)}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('transcript.total_credits', 'Credits Earned')}
            </CardTitle>
            <GraduationCap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{transcript.total_credits}</div>
            <p className="text-xs text-muted-foreground">
              {t('transcript.credits_description', 'Total credits passed')}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('transcript.generated_at', 'Generated At')}
            </CardTitle>
            <BookOpen className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-sm pt-2">
              {format(new Date(transcript.generated_at), 'PPP')}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Terms List */}
      <div className="space-y-6">
        {transcript.terms.length === 0 ? (
           <Card>
             <CardContent className="py-8 text-center text-muted-foreground">
               {t('transcript.no_records', 'No academic records found.')}
             </CardContent>
           </Card>
        ) : (
          transcript.terms.map((term: any) => (
            <Card key={term.term_id}>
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>{term.term_name}</CardTitle>
                    <CardDescription className="mt-1">
                      {t('transcript.term_gpa', 'Term GPA')}: {term.term_gpa.toFixed(2)} • {t('transcript.term_credits', 'Credits')}: {term.term_credits}
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t('transcript.course_code', 'Code')}</TableHead>
                      <TableHead>{t('transcript.course_title', 'Course Title')}</TableHead>
                      <TableHead className="text-right">{t('transcript.credits', 'Credits')}</TableHead>
                      <TableHead className="text-center">{t('transcript.grade', 'Grade')}</TableHead>
                      <TableHead className="text-right">{t('transcript.points', 'Points')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {term.grades.map((grade: any) => (
                      <TableRow key={grade.id}>
                        <TableCell className="font-medium">{grade.course_code}</TableCell>
                        <TableCell>{grade.course_title}</TableCell>
                        <TableCell className="text-right">{grade.credits}</TableCell>
                        <TableCell className="text-center">
                          <Badge variant={grade.is_passed ? 'secondary' : 'destructive'} className="w-8 justify-center">
                             {grade.grade}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">{grade.grade_points.toFixed(1)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  );
}
