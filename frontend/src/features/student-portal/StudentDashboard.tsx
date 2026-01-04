import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { AlertCircle, ArrowRight, Calendar, Clock, GraduationCap, Play } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import { cn } from '@/lib/utils';

export const StudentDashboard: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { user } = useAuth();

  const activeProgram = {
    title: user?.program || t('student.dashboard.default_program'),
    type: t('student.dashboard.program_type'),
    progress: 0,
    overdue: 0,
  };

  const tasks = [
    { title: t('student.dashboard.tasks.0.title'), due: t('student.dashboard.tasks.0.due'), type: 'urgent' as const },
    { title: t('student.dashboard.tasks.1.title'), due: t('student.dashboard.tasks.1.due'), type: 'normal' as const },
    { title: t('student.dashboard.tasks.2.title'), due: t('student.dashboard.tasks.2.due'), type: 'normal' as const },
  ];

  return (
    <div className="max-w-6xl mx-auto space-y-10 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div>
        <h1 className="text-3xl font-black text-slate-900 tracking-tight">
          {t('student.dashboard.welcome', { name: user?.first_name || '' })}
        </h1>
        <p className="text-slate-500 font-medium mt-1">{t('student.dashboard.subtitle')}</p>
      </div>

      <div className="relative bg-slate-900 rounded-[2.5rem] p-8 md:p-12 overflow-hidden shadow-2xl shadow-indigo-900/20 text-white">
        <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-indigo-600 rounded-full blur-[120px] opacity-30 -mr-20 -mt-20 pointer-events-none" />
        <div className="absolute bottom-0 left-0 w-[300px] h-[300px] bg-emerald-600 rounded-full blur-[100px] opacity-20 -ml-10 -mb-10 pointer-events-none" />

        <div className="relative z-10 flex flex-col md:flex-row gap-8 items-start md:items-center justify-between">
          <div className="space-y-4 max-w-2xl">
            <div className="flex items-center gap-2">
              <span className="bg-white/10 backdrop-blur-md px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider text-indigo-200 border border-white/10">
                {t('student.dashboard.current_focus')}
              </span>
              <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">{activeProgram.type}</span>
            </div>
            <h2 className="text-3xl md:text-4xl font-black leading-tight">{activeProgram.title}</h2>

            <div className="flex items-center gap-6 pt-2">
              <div>
                <div className="text-2xl font-black text-emerald-400">{activeProgram.progress}%</div>
                <div className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{t('student.dashboard.complete')}</div>
              </div>
              <div className="w-px h-8 bg-white/10" />
              <div>
                <div className="text-2xl font-black text-white">{activeProgram.overdue}</div>
                <div className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{t('student.dashboard.overdue')}</div>
              </div>
            </div>

            <div className="w-full bg-white/10 h-2 rounded-full overflow-hidden max-w-md">
              <div className="h-full bg-emerald-500 rounded-full" style={{ width: `${activeProgram.progress}%` }} />
            </div>
          </div>

          <div className="bg-white/5 backdrop-blur-md border border-white/10 p-6 rounded-3xl w-full md:w-80 flex flex-col gap-4">
            <div className="flex items-start gap-3">
              <div className="w-10 h-10 rounded-full bg-indigo-500 flex items-center justify-center flex-shrink-0 shadow-lg shadow-indigo-500/30">
                <Play size={20} fill="currentColor" className="ml-1" />
              </div>
              <div>
                <div className="text-xs font-bold text-indigo-300 uppercase tracking-wider mb-1">{t('student.dashboard.up_next')}</div>
                <div className="font-bold text-sm leading-tight">{t('student.dashboard.up_next_title')}</div>
                <div className="text-xs text-slate-400 mt-1 flex items-center gap-1">
                  <Clock size={12} /> {t('student.dashboard.up_next_time')}
                </div>
              </div>
            </div>
            <button
              onClick={() => navigate('/journey')}
              className="w-full py-3 bg-white text-slate-900 rounded-xl font-black text-sm hover:bg-indigo-50 transition-colors flex items-center justify-center gap-2"
            >
              {t('student.dashboard.continue')} <ArrowRight size={16} />
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          <h3 className="font-bold text-slate-900 text-lg">{t('student.dashboard.your_programs')}</h3>
          <div className="bg-white p-6 rounded-3xl border border-slate-100 shadow-sm flex items-center gap-6">
            <div className="w-16 h-16 bg-slate-100 rounded-2xl flex items-center justify-center text-slate-400">
              <GraduationCap size={32} />
            </div>
            <div className="flex-1">
              <h4 className="font-bold text-slate-900 text-lg">{activeProgram.title}</h4>
              <p className="text-xs text-slate-500 font-medium mt-1">{t('student.dashboard.program_hint')}</p>
            </div>
            <button
              onClick={() => navigate('/journey')}
              className="w-10 h-10 rounded-full border border-slate-200 flex items-center justify-center text-slate-400 hover:bg-slate-900 hover:text-white hover:border-transparent transition-all"
            >
              <ArrowRight size={20} />
            </button>
          </div>
        </div>

        <div className="space-y-6">
          <h3 className="font-bold text-slate-900 text-lg">{t('student.dashboard.upcoming_tasks')}</h3>
          <div className="bg-white rounded-3xl border border-slate-100 shadow-sm p-2">
            <div className="space-y-1">
              {tasks.map((task, i) => (
                <div key={i} className="p-4 hover:bg-slate-50 rounded-2xl flex items-start gap-3 transition-colors cursor-pointer group">
                  <div
                    className={cn(
                      'w-5 h-5 rounded-full border-2 flex-shrink-0 mt-0.5 group-hover:bg-indigo-500 group-hover:border-indigo-500 transition-colors',
                      task.type === 'urgent' ? 'border-red-400' : 'border-slate-300'
                    )}
                  />
                  <div>
                    <div className="font-bold text-sm text-slate-800 leading-snug">{task.title}</div>
                    <div
                      className={cn(
                        'text-xs font-bold mt-1 flex items-center gap-1',
                        task.type === 'urgent' ? 'text-red-500' : 'text-slate-400'
                      )}
                    >
                      <Calendar size={12} /> {task.due}
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <button className="w-full py-3 text-xs font-bold text-slate-400 hover:text-indigo-600 transition-colors border-t border-slate-50 mt-2">
              {t('student.dashboard.view_all_tasks')}
            </button>
          </div>

          <div className="bg-amber-50 border border-amber-100 rounded-2xl p-4 text-amber-900 text-sm flex gap-3">
            <AlertCircle className="mt-0.5" size={18} />
            <div>
              <div className="font-bold">{t('student.dashboard.notice_title')}</div>
              <div className="text-xs text-amber-800 mt-1">{t('student.dashboard.notice_body')}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

