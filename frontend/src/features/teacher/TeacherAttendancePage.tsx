import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Calendar, CheckCircle2, Clock, MapPin, Save, UserX, Users } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { getCourseRoster } from './api';
import { format } from 'date-fns';

export const TeacherAttendancePage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { courseId } = useParams<{ courseId: string }>();
  // Mock date for now, or use live date
  const [date, setDate] = useState(new Date());

  const rosterQuery = useQuery({
    queryKey: ['teacher', 'courses', courseId, 'roster'],
    queryFn: () => getCourseRoster(courseId!),
    enabled: !!courseId,
  });

  const [attendance, setAttendance] = useState<Record<string, 'present' | 'absent' | 'excused'>>({});

  const handleMark = (studentId: string, status: 'present' | 'absent' | 'excused') => {
    setAttendance(prev => ({ ...prev, [studentId]: status }));
  };

  const students = rosterQuery.data || [];

  // Initialize defaults if needed
  React.useEffect(() => {
    if (students.length > 0 && Object.keys(attendance).length === 0) {
       const initial: Record<string, 'present' | 'absent' | 'excused'> = {};
       students.forEach(s => initial[s.student_id] = 'present');
       setAttendance(initial);
    }
  }, [students]);

  const stats = {
    present: Object.values(attendance).filter(s => s === 'present').length,
    absent: Object.values(attendance).filter(s => s === 'absent').length,
    excused: Object.values(attendance).filter(s => s === 'excused').length,
    total: students.length
  };

  return (
    <div className="max-w-5xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <button 
             onClick={() => navigate(`/admin/teacher/courses/${courseId}`)}
             className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors mb-2"
          >
             <ArrowLeft size={16} /> {t('teacher.attendance.back_to_course')}
          </button>
          <h1 className="text-2xl font-black text-slate-900 tracking-tight">{t('teacher.attendance.title')}</h1>
          <div className="flex items-center gap-2 text-slate-500 font-medium text-sm mt-1">
             <Calendar size={14} /> {format(date, 'MMMM d, yyyy')}
             <span className="text-slate-300">•</span>
             <Clock size={14} /> 09:00 AM - 10:30 AM
             <span className="text-slate-300">•</span>
             <MapPin size={14} /> Room 301
          </div>
        </div>
        <div className="flex gap-2">
            <Button>{t('teacher.attendance.save_records')}</Button>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
         <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm">
            <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t('teacher.attendance.stats.present')}</div>
            <div className="text-2xl font-black text-emerald-600 mt-1">{stats.present}</div>
         </div>
         <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm">
            <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t('teacher.attendance.stats.absent')}</div>
            <div className="text-2xl font-black text-red-600 mt-1">{stats.absent}</div>
         </div>
         <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm">
            <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t('teacher.attendance.stats.excused')}</div>
            <div className="text-2xl font-black text-amber-600 mt-1">{stats.excused}</div>
         </div>
         <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm">
            <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t('teacher.attendance.stats.attendance_rate')}</div>
            <div className="text-2xl font-black text-slate-900 mt-1">
               {stats.total > 0 ? Math.round((stats.present / stats.total) * 100) : 0}%
            </div>
         </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm">
        <table className="w-full text-sm text-left">
          <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
             <tr>
               <th className="px-6 py-4">{t('teacher.attendance.table.student')}</th>
               <th className="px-6 py-4">{t('teacher.attendance.table.status')}</th>
               <th className="px-6 py-4 text-right">{t('teacher.attendance.table.action')}</th>
             </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
             {students.map(student => (
                <tr key={student.student_id} className="hover:bg-slate-50 transition-colors">
                   <td className="px-6 py-4 mx-auto font-bold text-slate-900">
                      <div className="flex items-center gap-3">
                         <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-xs font-bold text-slate-500">
                            {(student.student_name || student.student_id).charAt(0)}
                         </div>
                         <div>
                            <div>{student.student_name || student.student_id}</div>
                            <div className="text-xs text-slate-400 font-normal">{student.student_email}</div>
                         </div>
                      </div>
                   </td>
                   <td className="px-6 py-4">
                      {attendance[student.student_id] === 'present' && <Badge className="bg-emerald-100 text-emerald-700 hover:bg-emerald-200 border-none">Present</Badge>}
                      {attendance[student.student_id] === 'absent' && <Badge className="bg-red-100 text-red-700 hover:bg-red-200 border-none">Absent</Badge>}
                      {attendance[student.student_id] === 'excused' && <Badge className="bg-amber-100 text-amber-700 hover:bg-amber-200 border-none">Excused</Badge>}
                   </td>
                   <td className="px-6 py-4 text-right">
                      <div className="flex justify-end gap-1">
                         <Button 
                           size="sm" 
                           variant={attendance[student.student_id] === 'present' ? 'default' : 'ghost'} 
                           className={cn("h-8 w-8 p-0 rounded-full", attendance[student.student_id] === 'present' ? "bg-emerald-500 hover:bg-emerald-600" : "text-emerald-600 hover:bg-emerald-50")}
                           onClick={() => handleMark(student.student_id, 'present')}
                         >
                            <CheckCircle2 size={16} />
                         </Button>
                         <Button 
                           size="sm" 
                           variant={attendance[student.student_id] === 'absent' ? 'default' : 'ghost'}
                           className={cn("h-8 w-8 p-0 rounded-full", attendance[student.student_id] === 'absent' ? "bg-red-500 hover:bg-red-600" : "text-red-500 hover:bg-red-50")}
                           onClick={() => handleMark(student.student_id, 'absent')}
                         >
                            <UserX size={16} />
                         </Button>
                         <Button 
                           size="sm" 
                           variant={attendance[student.student_id] === 'excused' ? 'default' : 'ghost'}
                           className={cn("h-8 w-8 p-0 rounded-full", attendance[student.student_id] === 'excused' ? "bg-amber-500 hover:bg-amber-600" : "text-amber-500 hover:bg-amber-50")}
                           onClick={() => handleMark(student.student_id, 'excused')}
                         >
                            <Clock size={16} />
                         </Button>
                      </div>
                   </td>
                </tr>
             ))}
             {students.length === 0 && (
                <tr><td colSpan={3} className="p-12 text-center text-slate-400 italic">No students enrolled.</td></tr>
             )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
