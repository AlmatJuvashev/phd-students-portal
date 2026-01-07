import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { CheckCircle2, Clock, Loader2, XCircle, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { getAssessment } from './api';
import { useToast } from "@/components/ui/use-toast";
import type { QuizQuestion } from '../studio/types';

const supported: Record<string, boolean> = {
  MCQ: true,
  TRUE_FALSE: true,
  TEXT: true,
  // Mapping studio type names if backend returns those
  multiple_choice: true, 
  short_text: true,
};

const formatSeconds = (totalSeconds: number) => {
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
};

export const AssessmentPreviewPage: React.FC = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['assessment', id],
    queryFn: () => getAssessment(id!),
    enabled: Boolean(id),
  });

  const assessment = data?.assessment;
  const questions = (data?.questions || []) as any[]; // Using any to adapt between backend Question and studio QuizQuestion if needed

  const [answers, setAnswers] = useState<Record<string, any>>({});
  const [timeLeft, setTimeLeft] = useState(0);

  useEffect(() => {
    if (assessment?.time_limit_minutes) {
        setTimeLeft(assessment.time_limit_minutes * 60);
    }
  }, [assessment]);

  // Timer simulation
  useEffect(() => {
    if (timeLeft <= 0) return;
    const interval = setInterval(() => setTimeLeft(prev => Math.max(0, prev - 1)), 1000);
    return () => clearInterval(interval);
  }, [timeLeft]);

  const handleSelectOption = (questionId: string, optionId: string) => {
    setAnswers(prev => ({ ...prev, [questionId]: { optionId } }));
  };

  const handleTextChange = (questionId: string, text: string) => {
    setAnswers(prev => ({ ...prev, [questionId]: { text } }));
  };

  const handleSubmit = () => {
    toast({
      title: "Preview Mode",
      description: "This is just a preview. No attempt was created.",
      duration: 3000,
    });
  };

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto p-6 flex justify-center py-20">
         <div className="flex flex-col items-center gap-4 text-slate-500">
            <Loader2 className="animate-spin h-8 w-8 text-indigo-500" />
            <p className="font-medium animate-pulse">Loading preview...</p>
         </div>
      </div>
    );
  }

  if (isError || !assessment) {
    return (
      <div className="max-w-4xl mx-auto p-6"> 
        <div className="bg-red-50 border border-red-100 p-6 rounded-2xl text-red-800">
            <h3 className="font-bold text-lg mb-2">Failed to load assessment</h3>
            <p>{(error as Error)?.message || 'Unknown error'}</p>
            <Button variant="outline" className="mt-4 bg-white" onClick={() => navigate('/admin/assessments')}>
                Go Back
            </Button>
        </div>
      </div>
    );
  }

  const answeredCount = questions.filter(q => {
      const a = answers[q.id];
      if (!a) return false;
      if (q.type === 'TEXT' || q.type === 'short_text') return !!a.text?.trim();
      return !!a.optionId;
  }).length;

  return (
    <div className="min-h-screen bg-slate-50/50">
        
      {/* Simulation Banner */}
      <div className="bg-indigo-600 text-white px-4 py-2 text-center text-sm font-bold shadow-md sticky top-0 z-50">
          PREVIEW MODE — Answers are not saved
      </div>

      <div className="max-w-4xl mx-auto p-6 space-y-8">
          
          <Button variant="ghost" className="mb-4 pl-0 hover:bg-transparent hover:text-indigo-600" onClick={() => navigate('/admin/assessments')}>
              <ArrowLeft className="mr-2 h-4 w-4"/> Back to Assessments
          </Button>

          {/* Header Card */}
          <div className="bg-white border border-slate-200 rounded-3xl p-8 shadow-sm relative overflow-hidden">
               <div className="absolute top-0 right-0 p-3 opacity-10">
                   <Clock size={120} />
               </div>
               
               <div className="relative z-10 flex flex-col md:flex-row justify-between gap-6">
                   <div className="space-y-4">
                       <div>
                           <Badge variant="outline" className="mb-2 bg-indigo-50 text-indigo-700 border-indigo-100">Assessment Preview</Badge>
                           <h1 className="text-3xl font-black text-slate-900 tracking-tight">{assessment.title}</h1>
                       </div>
                       {assessment.description && (
                           <p className="text-slate-600 max-w-2xl text-lg leading-relaxed">{assessment.description}</p>
                       )}
                       <div className="flex gap-4 pt-2">
                           {assessment.time_limit_minutes ? (
                               <div className="flex items-center gap-2 text-sm font-bold text-slate-500">
                                   <Clock size={16} />
                                   Time Limit: {assessment.time_limit_minutes} mins
                               </div>
                           ) : (
                               <div className="flex items-center gap-2 text-sm font-bold text-slate-500">
                                   <Clock size={16} />
                                   No time limit
                               </div>
                           )}
                           <div className="flex items-center gap-2 text-sm font-bold text-slate-500">
                               <CheckCircle2 size={16} />
                               Passing Score: {assessment.passing_score}%
                           </div>
                       </div>
                   </div>

                   <div className="bg-slate-50 rounded-2xl p-6 min-w-[200px] border border-slate-100 flex flex-col items-center justify-center text-center">
                       <div className="text-xs font-bold text-slate-400 uppercase mb-1">Simulated Timer</div>
                       <div className="text-3xl font-black text-slate-900 font-mono mb-1">
                           {formatSeconds(timeLeft)}
                       </div>
                       <div className="text-xs text-slate-400">remaining</div>
                   </div>
               </div>
          </div>

          {/* Progress */}
          <div className="sticky top-14 z-40 bg-slate-50/95 backdrop-blur-sm py-4 border-b border-slate-200/50 mb-8 transition-all">
             <div className="flex items-center justify-between text-xs text-slate-500 mb-2">
                 <span>Simulated Progress</span>
                 <span className="font-bold text-slate-700">{Math.round((answeredCount / questions.length) * 100)}%</span>
             </div>
             <Progress value={(answeredCount / questions.length) * 100} className="h-2" />
          </div>

          {/* Questions */}
          <div className="space-y-6 pb-20">
              {questions.map((q, idx) => (
                  <div key={q.id} className="p-8 bg-white rounded-3xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
                      <div className="flex gap-4">
                          <div className="flex-none">
                              <div className="w-8 h-8 rounded-full bg-slate-100 text-slate-500 flex items-center justify-center font-bold text-sm">
                                  {idx + 1}
                              </div>
                          </div>
                          <div className="flex-1 space-y-6">
                              <div>
                                  <div className="text-lg font-bold text-slate-900 leading-snug">{q.stem || q.text}</div>
                                  <div className="mt-1 flex items-center gap-2">
                                      <Badge variant="secondary" className="text-[10px] uppercase">{q.type}</Badge>
                                      <span className="text-xs text-slate-400 font-bold">{q.points_default || q.points} Points</span>
                                  </div>
                              </div>

                              {/* MCQ / TrueFalse */}
                              {(q.type === 'MCQ' || q.type === 'TRUE_FALSE' || q.type === 'multiple_choice') && (
                                  <div className="grid gap-3">
                                      {q.options?.map((opt: any) => {
                                          const isSelected = answers[q.id]?.optionId === opt.id;
                                          return (
                                              <button
                                                  key={opt.id}
                                                  onClick={() => handleSelectOption(q.id, opt.id)}
                                                  className={cn(
                                                      "relative text-left p-4 rounded-xl border-2 transition-all flex items-center justify-between group",
                                                      isSelected 
                                                          ? "border-indigo-600 bg-indigo-50/50" 
                                                          : "border-slate-100 hover:border-slate-300 hover:bg-slate-50"
                                                  )}
                                              >
                                                  <div className="flex items-center gap-4">
                                                      <div className={cn(
                                                          "w-6 h-6 rounded-full border-2 flex items-center justify-center transition-colors",
                                                          isSelected ? "border-indigo-600 bg-indigo-600" : "border-slate-300 group-hover:border-slate-400"
                                                      )}>
                                                          {isSelected && <div className="w-2.5 h-2.5 bg-white rounded-full" />}
                                                      </div>
                                                      <span className={cn("font-medium", isSelected ? "text-indigo-900" : "text-slate-700")}>{opt.text}</span>
                                                  </div>
                                              </button>
                                          );
                                      })}
                                  </div>
                              )}

                              {/* Text */}
                              {(q.type === 'TEXT' || q.type === 'short_text') && (
                                  <Textarea 
                                      placeholder="Type your answer here..."
                                      className="min-h-[140px] text-lg p-4 resize-none bg-slate-50 border-slate-200 focus:bg-white transition-colors"
                                      value={answers[q.id]?.text || ''}
                                      onChange={(e) => handleTextChange(q.id, e.target.value)}
                                  />
                              )}
                              
                              {!supported[q.type] && !supported[q.type.replace('_', ' ')] && ( // Handle type mismatch if any
                                  <div className="p-4 bg-orange-50 text-orange-800 rounded-xl text-sm border border-orange-100">
                                      Preview for this question type is not fully styled yet.
                                  </div>
                              )}
                          </div>
                      </div>
                  </div>
              ))}
          </div>

          {/* Footer */}
          <div className="flex justify-end pt-8 border-t border-slate-200">
              <Button size="lg" className="w-full md:w-auto px-12 text-lg h-14 rounded-2xl shadow-lg shadow-indigo-500/20" onClick={handleSubmit}>
                  Submit Preview
              </Button>
          </div>
      
      </div>
    </div>
  );
};
