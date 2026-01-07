import React, { useState, useMemo, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  ArrowLeft, Eye, Share, Undo, Redo, 
  Loader2, Save, CheckCircle2, Wand2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { ActivityList } from './components/ActivityList';
import { ActivityDetails } from './components/ActivityDetails';
import { getCourseContent, updateCourseContent } from './api';
import { CourseContent, Activity, Module, Lesson } from './types';

import { AIAssistantModal } from './components/AIAssistantModal';

export const CourseBuilder: React.FC = () => {
  const { t } = useTranslation('common');
  const { courseId } = useParams<{ courseId: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  
  const [selectedActivityId, setSelectedActivityId] = useState<string | null>(null);
  const [courseStatus, setCourseStatus] = useState<'draft' | 'published'>('draft');
  const [isAIModalOpen, setIsAIModalOpen] = useState(false);
  
  // Queries
  const { data: content, isLoading, isError } = useQuery({
    queryKey: ['studio', 'course', courseId, 'content'],
    queryFn: () => getCourseContent(courseId!),
    enabled: !!courseId,
  });

  // Local state for editing (buffered)
  const [localContent, setLocalContent] = useState<CourseContent | null>(null);

  useEffect(() => {
    if (content) {
      setLocalContent(content);
      // Auto-select first activity if none selected
      if (!selectedActivityId && content.modules?.[0]?.lessons?.[0]?.activities?.[0]) {
        setSelectedActivityId(content.modules[0].lessons[0].activities[0].id);
      }
    }
  }, [content]);

  // Mutations
  const saveMutation = useMutation({
    mutationFn: (newContent: CourseContent) => updateCourseContent(courseId!, newContent),
    onSuccess: () => {
      queryClient.setQueryData(['studio', 'course', courseId, 'content'], localContent);
    }
  });

  // Autosave simulation
  const autosaveTimerRef = useRef<any>(null);
  const notifyChange = (newContent: CourseContent) => {
    setLocalContent(newContent);
    if (autosaveTimerRef.current) clearTimeout(autosaveTimerRef.current);
    autosaveTimerRef.current = setTimeout(() => {
      saveMutation.mutate(newContent);
    }, 2000);
  };

  // --- Handlers ---

  const handleUpdateActivity = (updates: Partial<Activity>) => {
    if (!localContent || !selectedActivityId) return;

    const newContent = {
      ...localContent,
      modules: localContent.modules.map(m => ({
        ...m,
        lessons: m.lessons.map(l => ({
          ...l,
          activities: l.activities.map(a => a.id === selectedActivityId ? { ...a, ...updates } : a)
        }))
      }))
    };
    notifyChange(newContent);
  };

  const handleAddModule = () => {
    if (!localContent) return;
    const newModule: Module = {
      id: `m_${Date.now()}`,
      title: 'New Module',
      order: localContent.modules.length + 1,
      lessons: []
    };
    notifyChange({ ...localContent, modules: [...localContent.modules, newModule] });
  };

  const handleAddLesson = (moduleId: string) => {
    if (!localContent) return;
    const newLesson: Lesson = {
      id: `l_${Date.now()}`,
      title: 'New Lesson',
      order: 1, // Should calculate based on existing
      activities: []
    };
    const newContent = {
      ...localContent,
      modules: localContent.modules.map(m => m.id === moduleId ? { ...m, lessons: [...m.lessons, newLesson] } : m)
    };
    notifyChange(newContent);
  };

  const handleAddActivity = (lessonId: string) => {
    if (!localContent) return;
    const newActivity: Activity = {
      id: `a_${Date.now()}`,
      title: 'New Activity',
      type: 'text',
      points: 0,
      is_optional: false,
      content: '',
      attachments: [],
      citations: []
    };
    const newContent = {
      ...localContent,
      modules: localContent.modules.map(m => ({
        ...m,
        lessons: m.lessons.map(l => l.id === lessonId ? { ...l, activities: [...l.activities, newActivity] } : l)
      }))
    };
    setSelectedActivityId(newActivity.id);
    notifyChange(newContent);
  };

  const handleApplyAIStructure = (generatedModules: Module[]) => {
    if (!localContent) return;
    
    // Append generated modules to existing ones
    // Adjust order
    const startingOrder = localContent.modules.length + 1;
    const orderedModules = generatedModules.map((m, idx) => ({
      ...m,
      order: startingOrder + idx
    }));

    const newContent = {
      ...localContent,
      modules: [...localContent.modules, ...orderedModules]
    };
    notifyChange(newContent);
  };

  const activeActivity = useMemo(() => {
    if (!localContent || !selectedActivityId) return null;
    for (const m of localContent.modules) {
      for (const l of m.lessons) {
        const a = l.activities.find(act => act.id === selectedActivityId);
        if (a) return a;
      }
    }
    return null;
  }, [localContent, selectedActivityId]);

  if (isLoading) return <div className="h-full flex items-center justify-center"><Loader2 className="animate-spin" /></div>;
  if (isError || !localContent) return <div className="p-8 text-center text-red-500">Error loading course content.</div>;

  return (
    <div className="flex flex-col h-[calc(100vh-8rem)] -m-8 bg-slate-50">
      {/* Header */}
      <div className="h-16 border-b border-slate-200 px-6 flex items-center justify-between bg-white z-20 flex-shrink-0">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
            <ArrowLeft size={18} />
          </Button>
          <div className="h-6 w-px bg-slate-200" />
          <div className="flex flex-col">
            <h2 className="font-bold text-slate-800 leading-none">{t('studio.title')}</h2>
            <div className="flex items-center gap-2 mt-1">
              <span className="text-[10px] font-bold uppercase text-slate-400 tracking-wider">Builder</span>
              <div className="flex items-center gap-1 bg-slate-100 rounded px-1.5 py-0.5">
                <span className={cn("w-1.5 h-1.5 rounded-full", courseStatus === 'published' ? 'bg-emerald-500' : 'bg-amber-500')} />
                <span className="text-[9px] uppercase font-bold text-slate-500">{t(`studio.${courseStatus}`)}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4">
           {saveMutation.isPending && <span className="text-xs text-slate-400 flex items-center gap-2"><Loader2 size={14} className="animate-spin" /> Saving...</span>}
           {saveMutation.isSuccess && !saveMutation.isPending && <span className="text-xs text-emerald-500 flex items-center gap-2"><CheckCircle2 size={14} /> Saved</span>}
           <div className="h-6 w-px bg-slate-200" />
           <Button variant="ghost" size="sm" onClick={() => setIsAIModalOpen(true)} className="text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50">
             <Wand2 size={16} className="mr-2" /> AI Assistant
           </Button>
           <Button variant="outline" size="sm">
              <Eye size={16} className="mr-2" /> {t('studio.preview')}
           </Button>
           <Button size="sm">
              <Share size={16} className="mr-2" /> {t('studio.publish')}
           </Button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        <ActivityList 
          modules={localContent.modules}
          selectedActivityId={selectedActivityId}
          onSelectActivity={setSelectedActivityId}
          onAddModule={handleAddModule}
          onAddLesson={handleAddLesson}
          onAddActivity={handleAddActivity}
        />
        
        <div className="flex-1 overflow-hidden relative">
          {activeActivity ? (
            <ActivityDetails activity={activeActivity} onUpdate={handleUpdateActivity} />
          ) : (
            <div className="h-full flex flex-col items-center justify-center text-slate-400">
               <Loader2 size={48} className="mb-4 opacity-10" />
               <p>{t('studio.sidebar.no_content')}</p>
            </div>
          )}
        </div>
      </div>

      <AIAssistantModal 
        isOpen={isAIModalOpen}
        onClose={() => setIsAIModalOpen(false)}
        onApply={handleApplyAIStructure}
      />
    </div>
  );
};
