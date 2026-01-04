import React from 'react';
import { 
  ChevronDown, ChevronRight, Plus, GripVertical, 
  FileText, Video, CheckSquare, GraduationCap, 
  ClipboardList, Trash2, MoreHorizontal
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/utils';
import { Module, Lesson, Activity } from '../types';

interface ActivityListProps {
  modules: Module[];
  selectedActivityId: string | null;
  onSelectActivity: (id: string) => void;
  onAddModule: () => void;
  onAddLesson: (moduleId: string) => void;
  onAddActivity: (lessonId: string) => void;
}

const ActivityIcon = ({ type }: { type: string }) => {
  switch (type) {
    case 'video': return <Video size={14} />;
    case 'quiz': return <CheckSquare size={14} />;
    case 'survey': return <ClipboardList size={14} />;
    case 'assignment': return <GraduationCap size={14} />;
    default: return <FileText size={14} />;
  }
};

export const ActivityList: React.FC<ActivityListProps> = ({ 
  modules, 
  selectedActivityId, 
  onSelectActivity,
  onAddModule,
  onAddLesson,
  onAddActivity
}) => {
  const { t } = useTranslation('common');

  return (
    <div className="flex flex-col h-full bg-slate-50 border-r border-slate-200 w-80">
      <div className="p-4 border-b border-slate-200 flex items-center justify-between bg-white">
        <h3 className="font-bold text-slate-800 text-sm">{t('studio.sidebar.modules')}</h3>
        <button 
          onClick={onAddModule}
          className="p-1.5 hover:bg-slate-100 rounded-lg text-indigo-600 transition-colors"
          title={t('studio.sidebar.add_module')}
        >
          <Plus size={18} />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto p-2 space-y-4">
        {modules.length === 0 && (
          <div className="text-center py-8 text-slate-400 text-xs italic">
            {t('studio.sidebar.no_content')}
          </div>
        )}

        {modules.map((module) => (
          <div key={module.id} className="space-y-1">
            <div className="flex items-center gap-2 p-2 group">
              <div className="text-slate-400"><GripVertical size={14} /></div>
              <h4 className="font-bold text-xs text-slate-700 flex-1 uppercase tracking-wider">{module.title}</h4>
              <button 
                onClick={() => onAddLesson(module.id)}
                className="opacity-0 group-hover:opacity-100 p-1 hover:bg-slate-200 rounded text-slate-500 transition-all"
              >
                <Plus size={14} />
              </button>
            </div>

            <div className="pl-2 space-y-1">
              {module.lessons.map((lesson) => (
                <div key={lesson.id} className="space-y-1">
                  <div className="flex items-center gap-2 p-2 rounded-lg hover:bg-slate-200/50 group cursor-pointer">
                    <ChevronDown size={14} className="text-slate-400" />
                    <span className="text-sm font-medium text-slate-600 flex-1">{lesson.title}</span>
                    <button 
                      onClick={(e) => { e.stopPropagation(); onAddActivity(lesson.id); }}
                      className="opacity-0 group-hover:opacity-100 p-1 hover:bg-slate-200 rounded text-slate-500 transition-all"
                    >
                      <Plus size={14} />
                    </button>
                  </div>

                  <div className="pl-6 space-y-0.5">
                    {lesson.activities.map((activity) => (
                      <button
                        key={activity.id}
                        onClick={() => onSelectActivity(activity.id)}
                        className={cn(
                          "w-full flex items-center gap-3 p-2 rounded-lg text-sm transition-all text-left",
                          selectedActivityId === activity.id 
                            ? "bg-white shadow-sm border border-slate-200 text-indigo-600 font-semibold" 
                            : "text-slate-500 hover:bg-slate-200/50"
                        )}
                      >
                        <ActivityIcon type={activity.type} />
                        <span className="truncate">{activity.title}</span>
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
