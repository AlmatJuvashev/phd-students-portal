import React from 'react';
import { NavLink, Outlet, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { BookOpen, GraduationCap, LayoutDashboard, ListChecks, Map } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAuth } from '@/contexts/AuthContext';

export const StudentLayout: React.FC = () => {
  const { t } = useTranslation('common');
  const { user } = useAuth();
  const location = useLocation();

  const active = (to: string) =>
    location.pathname === to ? 'bg-slate-900 text-white' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900';

  return (
    <div className="max-w-6xl mx-auto px-4 py-6 space-y-6">
      <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm flex flex-col md:flex-row gap-4 md:items-center md:justify-between">
        <div className="min-w-0">
          <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">
            {t('student.layout.title', { defaultValue: 'Student Portal' })}
          </div>
          <div className="mt-1 text-xl font-black text-slate-900 truncate">
            {t('student.layout.welcome', { defaultValue: 'Welcome, {{name}}', name: user?.first_name || '' })}
          </div>
          {user?.program && <div className="text-sm text-slate-500 mt-1 truncate">{user.program}</div>}
        </div>

        <nav className="flex flex-wrap gap-2">
          <NavLink to="/student/dashboard" className={cn('px-4 py-2 rounded-xl text-sm font-bold flex items-center gap-2', active('/student/dashboard'))}>
            <LayoutDashboard size={16} /> {t('student.nav.dashboard', { defaultValue: 'Dashboard' })}
          </NavLink>
          <NavLink to="/student/courses" className={cn('px-4 py-2 rounded-xl text-sm font-bold flex items-center gap-2', active('/student/courses'))}>
            <BookOpen size={16} /> {t('student.nav.courses', { defaultValue: 'My Courses' })}
          </NavLink>
          <NavLink
            to="/student/assignments"
            className={cn('px-4 py-2 rounded-xl text-sm font-bold flex items-center gap-2', active('/student/assignments'))}
          >
            <ListChecks size={16} /> {t('student.nav.assignments', { defaultValue: 'Assignments' })}
          </NavLink>
          <NavLink to="/student/grades" className={cn('px-4 py-2 rounded-xl text-sm font-bold flex items-center gap-2', active('/student/grades'))}>
            <GraduationCap size={16} /> {t('student.nav.grades', { defaultValue: 'Grades' })}
          </NavLink>
          <NavLink to="/journey" className={cn('px-4 py-2 rounded-xl text-sm font-bold flex items-center gap-2', active('/journey'))}>
            <Map size={16} /> {t('student.nav.journey', { defaultValue: 'Journey' })}
          </NavLink>
        </nav>
      </div>

      <Outlet />
    </div>
  );
};

