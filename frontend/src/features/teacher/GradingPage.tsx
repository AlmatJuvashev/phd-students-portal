
import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  CheckSquare, Filter, Search, ChevronDown, 
  FileText, ExternalLink, Clock, CheckCircle2,
  X, AlertCircle, Loader2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { motion, AnimatePresence } from 'framer-motion';
import { getTeacherSubmissions, submitGradeForSubmission } from './api';
import { ActivitySubmission } from './types';
import { format } from 'date-fns';

export const GradingPage = () => {
  const queryClient = useQueryClient();
  const [filter, setFilter] = useState<'all' | 'ungraded' | 'graded'>('ungraded');
  const [search, setSearch] = useState('');
  const [activeSubmission, setActiveSubmission] = useState<ActivitySubmission | null>(null);
  const [currentScore, setCurrentScore] = useState<number>(0);
  const [feedback, setFeedback] = useState('');

  // Fetch submissions from API
  const { data: submissions = [], isLoading } = useQuery({
    queryKey: ['teacher-submissions'],
    queryFn: getTeacherSubmissions,
  });

  // Submit grade mutation
  const { mutate: submitGrade, isPending } = useMutation({
    mutationFn: submitGradeForSubmission,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teacher-submissions'] });
      setActiveSubmission(null);
    }
  });

  const filtered = submissions.filter(s => {
    // Basic search on student name (if available) or activity title
    const searchLower = search.toLowerCase();
    const studentName = s.student_name || 'Unknown Student';
    const matchesSearch = studentName.toLowerCase().includes(searchLower) || (s.activity_title || '').toLowerCase().includes(searchLower);
    
    const isGraded = s.status === 'graded';
    const matchesFilter = filter === 'all' || (filter === 'ungraded' ? !isGraded : isGraded);
    
    return matchesSearch && matchesFilter;
  });

  const openSubmission = (sub: ActivitySubmission) => {
    setActiveSubmission(sub);
    // Initialize score appropriately (assuming content or separate grade field holds it)
    // For now defaulting to 0 or existing logic if fields existed
    setCurrentScore(0); 
    setFeedback('');
  };

  const handleSubmitGrade = () => {
    if (!activeSubmission) return;
    submitGrade({
      course_offering_id: activeSubmission.course_offering_id,
      activity_id: activeSubmission.activity_id,
      student_id: activeSubmission.student_id,
      score: currentScore,
      max_score: 100, // Assuming default or fetch from activity details if available
      feedback: feedback
    });
  };

  if (isLoading) {
    return <div className="h-full flex items-center justify-center"><Loader2 className="animate-spin text-slate-400" /></div>;
  }

  return (
    <div className="h-full flex flex-col space-y-6 p-6">
       <div className="flex justify-between items-end flex-shrink-0">
          <div>
             <h1 className="text-2xl font-black text-slate-900 tracking-tight">Grading Hub</h1>
             <p className="text-slate-500 text-sm mt-1">Review and grade student submissions.</p>
          </div>
       </div>

       {/* Toolbar */}
       <div className="bg-white p-2 rounded-2xl border border-slate-200 shadow-sm flex flex-col sm:flex-row gap-4 flex-shrink-0 items-center">
          <div className="relative flex-1">
             <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
             <input 
               value={search}
               onChange={(e) => setSearch(e.target.value)}
               placeholder="Search student or assignment..." 
               className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus:outline-none focus:ring-0 text-sm"
             />
          </div>
          <div className="h-10 w-px bg-slate-200 hidden sm:block" />
          <div className="flex gap-1 bg-slate-100 p-1 rounded-xl">
             {(['all', 'ungraded', 'graded'] as const).map(f => (
               <button
                 key={f}
                 onClick={() => setFilter(f)}
                 className={cn(
                   "px-4 py-1.5 rounded-lg text-xs font-bold capitalize transition-all",
                   filter === f ? "bg-white text-slate-900 shadow-sm" : "text-slate-500 hover:text-slate-700"
                 )}
               >
                 {f}
               </button>
             ))}
          </div>
       </div>

       {/* Content */}
       <div className="flex-1 bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm flex flex-col">
          <div className="overflow-y-auto flex-1">
             <table className="w-full text-sm text-left">
                <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase sticky top-0 z-10 w-full">
                   <tr>
                      <th className="px-6 py-4">Student</th>
                      <th className="px-6 py-4">Assignment</th>
                      <th className="px-6 py-4">Submitted</th>
                      <th className="px-6 py-4">Status</th>
                      <th className="px-6 py-4 text-right">Action</th>
                   </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                   {filtered.map(sub => (
                      <tr key={sub.id} className="hover:bg-slate-50 transition-colors group cursor-pointer" onClick={() => openSubmission(sub)}>
                         <td className="px-6 py-4 font-bold text-slate-900">{sub.student_name || 'Unknown'}</td>
                         <td className="px-6 py-4">
                            <div className="font-medium text-slate-800">{sub.activity_title || 'Untitled Assignment'}</div>
                         </td>
                         <td className="px-6 py-4 text-slate-600">
                            {format(new Date(sub.submitted_at), 'MMM d, h:mm a')}
                         </td>
                         <td className="px-6 py-4">
                            {sub.status === 'graded' ? (
                               <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 text-xs font-bold border border-emerald-100">
                                  <CheckCircle2 size={12} /> Graded
                               </span>
                            ) : (
                               <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-indigo-50 text-indigo-700 text-xs font-bold border border-indigo-100">
                                  <Clock size={12} /> Needs Grading
                               </span>
                            )}
                         </td>
                         <td className="px-6 py-4 text-right">
                            <Button size="sm" variant="secondary" className="opacity-0 group-hover:opacity-100 transition-opacity">
                               {sub.status === 'graded' ? 'Review' : 'Grade'}
                            </Button>
                         </td>
                      </tr>
                   ))}
                   {filtered.length === 0 && (
                      <tr>
                         <td colSpan={5} className="px-6 py-12 text-center text-slate-400 italic">No submissions found.</td>
                      </tr>
                   )}
                </tbody>
             </table>
          </div>
       </div>

       {/* Grading Modal Overlay */}
       <AnimatePresence>
          {activeSubmission && (
             <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" onClick={() => setActiveSubmission(null)}>
                <motion.div 
                  initial={{ opacity: 0, scale: 0.95 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.95 }}
                  className="bg-white w-full max-w-4xl h-[80vh] rounded-3xl shadow-2xl flex flex-col overflow-hidden"
                  onClick={(e) => e.stopPropagation()}
                >
                   {/* Header */}
                   <div className="px-8 py-5 border-b border-slate-200 flex justify-between items-center bg-slate-50">
                      <div>
                         <h2 className="text-lg font-black text-slate-900">{activeSubmission.activity_title}</h2>
                         <p className="text-xs text-slate-500 font-bold mt-1 uppercase tracking-wide">
                            {activeSubmission.student_name}
                         </p>
                      </div>
                      <div className="flex items-center gap-4">
                         <Button variant="ghost" size="icon" onClick={() => setActiveSubmission(null)} className="rounded-full">
                            <X size={20} />
                         </Button>
                      </div>
                   </div>

                   {/* Body */}
                   <div className="flex-1 flex overflow-hidden">
                      {/* Left: Document Viewer Placeholder */}
                      <div className="flex-1 bg-slate-100 flex flex-col items-center justify-center border-r border-slate-200 p-8">
                         <div className="bg-white p-12 rounded-2xl shadow-sm border border-slate-200 text-center max-w-md">
                            <FileText size={48} className="mx-auto text-slate-300 mb-4" />
                            <h3 className="font-bold text-slate-700 text-lg mb-2">Submission Content</h3>
                            <div className="text-sm text-slate-500 mb-6 bg-slate-50 p-4 rounded text-left overflow-auto max-h-40">
                                <pre>{JSON.stringify(activeSubmission.content, null, 2)}</pre>
                            </div>
                         </div>
                      </div>

                      {/* Right: Grading Tools */}
                      <div className="w-96 bg-white flex flex-col overflow-y-auto">
                         <div className="p-6 space-y-6">
                            {/* Score Input */}
                            <div className="p-4 bg-slate-50 rounded-2xl border border-slate-200 space-y-3">
                               <label className="text-xs font-bold text-slate-500 uppercase">Grade</label>
                               <div className="flex gap-2">
                                  <Input 
                                    type="number" 
                                    className="text-lg font-bold" 
                                    placeholder="0" 
                                    value={currentScore}
                                    onChange={(e) => setCurrentScore(parseInt(e.target.value) || 0)} 
                                  />
                                  <div className="flex items-center justify-center bg-white border border-slate-200 rounded-lg px-3 font-bold text-slate-400 text-sm">/ 100</div>
                               </div>
                            </div>

                            {/* Rubric (Mock for now, needs backend support for structure) */}
                            <div className="space-y-4">
                               <label className="text-xs font-bold text-slate-500 uppercase flex items-center gap-2">
                                  <CheckSquare size={14} /> Rubric
                               </label>
                               <div className="p-4 bg-slate-50 rounded-xl text-xs text-slate-500 italic text-center">
                                  Rubric data not yet linked.
                               </div>
                            </div>

                            {/* Feedback */}
                            <div className="space-y-2">
                               <label className="text-xs font-bold text-slate-500 uppercase">Feedback</label>
                               <Textarea 
                                 className="w-full h-32 resize-none"
                                 placeholder="Enter comments for the student..."
                                 value={feedback}
                                 onChange={(e) => setFeedback(e.target.value)}
                               />
                            </div>
                         </div>
                         
                         <div className="mt-auto p-6 border-t border-slate-100">
                            <Button onClick={handleSubmitGrade} disabled={isPending} className="w-full py-6 text-lg font-bold">
                               {isPending ? <Loader2 className="animate-spin mr-2" /> : null}
                               Submit Grade
                            </Button>
                         </div>
                      </div>
                   </div>
                </motion.div>
             </div>
          )}
       </AnimatePresence>
    </div>
  );
};
