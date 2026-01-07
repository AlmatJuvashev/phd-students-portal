
import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { 
  Search, Filter, Mail, Loader2,
  TrendingDown, TrendingUp, AlertCircle, Eye
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { motion, AnimatePresence } from 'framer-motion';
// Assuming we might have a specific course context, or we fetch for all active courses
// For this polished view, let's fetch for the first active course or allow selection.
// To keep it simple for the initial port, we'll fetch 'at risk' for a known course or just mock if course ID isn't available easily globally.
// Ideally, we accept `courseId` as a prop or context.
import { getTeacherCourses, getAtRiskStudents, getCourseStudents } from './api';
import { StudentRiskProfile } from './types';

export const StudentTracker = () => {
  const [search, setSearch] = useState('');
  const [selectedCourseId, setSelectedCourseId] = useState<string | null>(null);

  // 1. Fetch Courses to select context
  const { data: courses = [] } = useQuery({
    queryKey: ['teacher-courses'],
    queryFn: getTeacherCourses,
  });

  // Default to first course if not selected
  const effectiveCourseId = selectedCourseId || courses[0]?.id;

  // 2. Fetch Students for the course
  const { data: students = [], isLoading } = useQuery({
    queryKey: ['course-students', effectiveCourseId],
    queryFn: () => effectiveCourseId ? getCourseStudents(effectiveCourseId) : Promise.resolve([]),
    enabled: !!effectiveCourseId,
  });

  const filtered = students.filter(s => 
    s.student_name.toLowerCase().includes(search.toLowerCase())
  );

  const getStatusColor = (status: string) => {
    switch(status.toLowerCase()) {
      case 'low': return 'bg-emerald-100 text-emerald-700';
      case 'medium': return 'bg-amber-100 text-amber-700';
      case 'high': 
      case 'critical': return 'bg-red-100 text-red-700';
      default: return 'bg-slate-100 text-slate-600';
    }
  };

  return (
    <div className="h-full flex flex-col space-y-6 p-6">
       <div className="flex justify-between items-end flex-shrink-0">
          <div>
             <h1 className="text-2xl font-black text-slate-900 tracking-tight">Student Tracker</h1>
             <p className="text-slate-500 text-sm mt-1">Monitor engagement and academic performance.</p>
          </div>
          <div className="flex gap-2">
             <Button variant="outline" className="gap-2"><Mail size={16} /> Message Course</Button>
          </div>
       </div>

       {/* Toolbar */}
       <div className="bg-white p-2 rounded-2xl border border-slate-200 shadow-sm flex flex-col sm:flex-row gap-4 flex-shrink-0 items-center">
          {/* Course Selector */}
          <select 
            className="h-10 px-3 bg-slate-50 border-none rounded-lg text-sm font-bold text-slate-700 focus:ring-0 cursor-pointer"
            value={effectiveCourseId || ''}
            onChange={(e) => setSelectedCourseId(e.target.value)}
          >
            {courses.map(c => (
                <option key={c.id} value={c.id}>
                    {c.section} ({(c as any).code || 'Course'})
                </option>
            ))}
            {courses.length === 0 && <option>No Active Courses</option>}
          </select>

          <div className="h-6 w-px bg-slate-200 hidden sm:block" />

          <div className="relative flex-1">
             <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
             <input 
               value={search}
               onChange={(e) => setSearch(e.target.value)}
               placeholder="Search students..." 
               className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus:outline-none focus:ring-0 text-sm"
             />
          </div>
       </div>

       {/* List */}
       <div className="flex-1 bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm flex flex-col">
          {isLoading ? (
             <div className="flex-1 flex items-center justify-center">
                <Loader2 className="animate-spin text-indigo-600" />
             </div>
          ) : (
            <div className="overflow-y-auto flex-1">
               <table className="w-full text-sm text-left">
                  <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase sticky top-0 z-10">
                     <tr>
                        <th className="px-6 py-4">Student</th>
                        <th className="px-6 py-4">Status</th>
                        <th className="px-6 py-4">Performance</th>
                        <th className="px-6 py-4">Engagement</th>
                        <th className="px-6 py-4 text-right">Actions</th>
                     </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                     {filtered.map(student => (
                        <tr key={student.student_id} className="hover:bg-slate-50 transition-colors group">
                           <td className="px-6 py-4">
                              <div className="flex items-center gap-3">
                                 <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-xs font-bold text-slate-500 border border-slate-200">
                                    {student.student_name.charAt(0)}
                                 </div>
                                 <div className="font-bold text-slate-900">{student.student_name}</div>
                              </div>
                           </td>
                           <td className="px-6 py-4">
                                <Badge className={cn("text-[10px] px-2 py-0.5 rounded-full font-bold uppercase border-none hover:bg-opacity-80", getStatusColor(student.risk_level))}>
                                  {student.risk_level}
                                </Badge>
                           </td>
                           <td className="px-6 py-4">
                              <div className="flex items-center gap-4">
                                 <div className="flex-1 min-w-[80px]">
                                    <div className="flex justify-between text-[10px] font-bold text-slate-400 mb-1">
                                       <span>GPA {student.average_grade}</span>
                                       <span>{student.overall_progress}%</span>
                                    </div>
                                    <div className="w-full h-1.5 bg-slate-100 rounded-full overflow-hidden">
                                       <div className={cn("h-full rounded-full", student.overall_progress < 40 ? "bg-red-500" : "bg-emerald-500")} style={{ width: `${student.overall_progress}%` }} />
                                    </div>
                                 </div>
                              </div>
                           </td>
                           <td className="px-6 py-4 text-slate-500 text-xs">
                              {student.assignments_completed}/{student.assignments_total} completed
                              {student.days_inactive > 7 && <span className="ml-2 text-red-500 font-bold">{student.days_inactive}d inactive</span>}
                           </td>
                           <td className="px-6 py-4 text-right">
                              <div className="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                 <Button variant="ghost" size="sm" className="hover:text-indigo-600 gap-2">
                                    <Eye size={14} /> View Profile
                                 </Button>
                              </div>
                           </td>
                        </tr>
                     ))}
                     {filtered.length === 0 && (
                        <tr><td colSpan={5} className="text-center py-10 text-slate-400">No students found.</td></tr>
                     )}
                  </tbody>
               </table>
            </div>
          )}
       </div>
    </div>
  );
};
