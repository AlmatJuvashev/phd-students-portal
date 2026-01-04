import React, { useState } from 'react';
import { 
  Plus, Trash2, Video, Mail, Trash,
  AlertCircle, Link as LinkIcon, Paperclip, FormInput, ClipboardList, CheckCircle2
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';
import { Activity, ActivityType, FormField } from '../types';
import { MarkdownEditor } from './MarkdownEditor';
import { FormBuilderModal } from './FormBuilderModal';
import { ChecklistBuilderModal } from './ChecklistBuilderModal';
import { QuizBuilderModal } from './QuizBuilderModal';
import { SurveyBuilderModal } from './SurveyBuilderModal';
import { ConfirmTaskBuilderModal } from './ConfirmTaskBuilderModal';
import { EduStudioHub } from './EduStudioHub';

interface ActivityDetailsProps {
  activity: Activity;
  onUpdate: (updates: Partial<Activity>) => void;
}

const ACTIVITY_TYPES: { id: ActivityType; label: string; icon: any }[] = [
  { id: 'text', label: 'Text / Article', icon: LinkIcon },
  { id: 'video', label: 'Video Lesson', icon: Video },
  { id: 'quiz', label: 'Quiz', icon: AlertCircle },
  { id: 'survey', label: 'Survey', icon: ClipboardList },
  { id: 'assignment', label: 'Assignment / Task', icon: FormInput },
];

export const ActivityDetails: React.FC<ActivityDetailsProps> = ({ activity, onUpdate }) => {
  const { t } = useTranslation('common');
  const [showFormBuilder, setShowFormBuilder] = useState(false);
  const [showChecklistBuilder, setShowChecklistBuilder] = useState(false);
  const [showQuizBuilder, setShowQuizBuilder] = useState(false);
  const [showSurveyBuilder, setShowSurveyBuilder] = useState(false);

  return (
    <div className="flex-1 flex flex-col h-full bg-white overflow-y-auto">
      {/* Activity Header */}
      <div className="p-6 border-b border-slate-100 flex items-center justify-between sticky top-0 bg-white/80 backdrop-blur-md z-10">
        <div className="flex-1 flex items-center gap-4">
          <div className="p-2 bg-indigo-50 text-indigo-600 rounded-xl">
             {ACTIVITY_TYPES.find(t => t.id === activity.type)?.icon && React.createElement(ACTIVITY_TYPES.find(t => t.id === activity.type)!.icon, { size: 24 })}
          </div>
          <div className="flex-1">
             <input 
               value={activity.title}
               onChange={(e) => onUpdate({ title: e.target.value })}
               className="text-xl font-black text-slate-900 bg-transparent border-none outline-none focus:ring-0 w-full p-0"
               placeholder={t('studio.editor.activity_title')}
             />
             <div className="flex items-center gap-4 mt-1">
                <div className="flex items-center gap-2">
                   <Label className="text-[10px] font-bold text-slate-400 uppercase">{t('studio.editor.points')}</Label>
                   <input 
                     type="number" 
                     value={activity.points}
                     onChange={(e) => onUpdate({ points: parseInt(e.target.value) || 0 })}
                     className="w-12 h-6 text-xs font-bold text-indigo-600 bg-slate-100 rounded border-none outline-none focus:ring-2 focus:ring-indigo-100 text-center"
                   />
                </div>
                <div className="flex items-center gap-2">
                   <Label className="text-[10px] font-bold text-slate-400 uppercase">{t('studio.editor.optional')}</Label>
                   <Switch 
                     checked={activity.is_optional}
                     onCheckedChange={(checked) => onUpdate({ is_optional: checked })}
                   />
                </div>
             </div>
          </div>
        </div>
      </div>

      <div className="p-6 space-y-8 max-w-4xl mx-auto w-full">
        {/* Type-Specific Configs */}
        {activity.type === 'video' && (
          <div className="space-y-4">
            <h3 className="text-sm font-bold text-slate-900 flex items-center gap-2">
               <Video size={16} className="text-indigo-500" />
               {t('studio.editor.video_urls')}
            </h3>
            <div className="space-y-2">
              {(activity.video_urls || []).map((url, i) => (
                <div key={i} className="flex gap-2">
                  <Input 
                    value={url} 
                    onChange={(e) => {
                      const newUrls = [...(activity.video_urls || [])];
                      newUrls[i] = e.target.value;
                      onUpdate({ video_urls: newUrls });
                    }}
                    placeholder="https://youtube.com/..." 
                  />
                  <Button variant="ghost" size="icon" onClick={() => {
                    onUpdate({ video_urls: (activity.video_urls || []).filter((_, idx) => idx !== i) });
                  }}>
                    <Trash2 size={16} />
                  </Button>
                </div>
              ))}
              <Button variant="outline" size="sm" onClick={() => onUpdate({ video_urls: [...(activity.video_urls || []), ''] })}>
                <Plus size={14} className="mr-1" /> {t('studio.editor.add_video')}
              </Button>
            </div>
          </div>
        )}

        {activity.type === 'quiz' && (
           <div className="p-6 bg-indigo-50 rounded-2xl border border-indigo-100 border-dashed flex flex-col items-center justify-center text-center gap-4">
               <div className="w-12 h-12 bg-white rounded-full flex items-center justify-center text-indigo-600 shadow-sm">
                   <AlertCircle size={24} />
               </div>
               <div>
                   <h3 className="font-bold text-slate-900">Quiz Configuration</h3>
                   <p className="text-xs text-slate-500">{activity.quiz_config?.questions?.length || 0} questions configured</p>
               </div>
               <Button onClick={() => setShowQuizBuilder(true)} variant="default" className="bg-indigo-600 hover:bg-indigo-700">Open Quiz Builder</Button>
           </div>
        )}

        {activity.type === 'survey' && (
           <div className="p-6 bg-rose-50 rounded-2xl border border-rose-100 border-dashed flex flex-col items-center justify-center text-center gap-4">
               <div className="w-12 h-12 bg-white rounded-full flex items-center justify-center text-rose-600 shadow-sm">
                   <ClipboardList size={24} />
               </div>
               <div>
                   <h3 className="font-bold text-slate-900">Survey Configuration</h3>
                   <p className="text-xs text-slate-500">{activity.survey_config?.questions?.length || 0} items configured</p>
               </div>
               <Button onClick={() => setShowSurveyBuilder(true)} variant="default" className="bg-rose-600 hover:bg-rose-700">Open Survey Builder</Button>
           </div>
        )}

        {activity.type === 'assignment' && (
          <div className="space-y-4">
             {/* Sub-type Selection */}
             <div className="p-4 bg-indigo-50 rounded-xl border border-indigo-100 space-y-4">
                 <div className="flex gap-2">
                    <Button 
                      size="sm" 
                      variant={(activity.assignment_config?.submission_types || []).includes('form') ? 'default' : 'outline'}
                      onClick={() => onUpdate({ 
                        assignment_config: { ...activity.assignment_config, submission_types: ['form'], group_assignment: false, peer_review: false } 
                      })}
                    >
                      Custom Form
                    </Button>
                    <Button 
                      size="sm" 
                      variant={(activity.assignment_config?.submission_types || []).includes('checklist') ? 'default' : 'outline'}
                      onClick={() => onUpdate({ 
                        assignment_config: { ...activity.assignment_config, submission_types: ['checklist'], group_assignment: false, peer_review: false } 
                      })}
                    >
                      Checklist
                    </Button>
                    <Button 
                      size="sm" 
                      variant={(activity.assignment_config?.submission_types || []).includes('file_upload') ? 'default' : 'outline'}
                      onClick={() => onUpdate({ 
                        assignment_config: { ...activity.assignment_config, submission_types: ['file_upload'], group_assignment: false, peer_review: false } 
                      })}
                    >
                      File Upload
                    </Button>
                 </div>

                 {(activity.assignment_config?.submission_types || []).includes('form') && (
                   <div className="pt-2 border-t border-indigo-200/50 flex items-center justify-between">
                     <span className="text-xs font-bold text-indigo-700">
                        {activity.assignment_config?.form_fields?.length || 0} fields configured
                     </span>
                     <Button onClick={() => setShowFormBuilder(true)} size="sm" variant="secondary">Configure Form</Button>
                   </div>
                 )}

                 {(activity.assignment_config?.submission_types || []).includes('checklist') && (
                   <div className="pt-2 border-t border-indigo-200/50 flex items-center justify-between">
                     <span className="text-xs font-bold text-indigo-700">
                        {activity.assignment_config?.checklist_config?.items.length || 0} tasks configured
                     </span>
                     <Button onClick={() => setShowChecklistBuilder(true)} size="sm" variant="secondary">Configure Checklist</Button>
                   </div>
                 )}
             </div>
          </div>
        )}

        {/* Content Editor */}
        <div className="space-y-4">
          <h3 className="text-sm font-bold text-slate-900 flex items-center gap-2">
             <LinkIcon size={16} className="text-indigo-500" />
             {t('studio.editor.content')}
          </h3>
          <MarkdownEditor value={activity.content} onChange={(v) => onUpdate({ content: v })} />
        </div>

        {/* Attachments & Citations */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
           <div className="space-y-4">
              <h3 className="text-sm font-bold text-slate-900 flex items-center gap-2">
                 <Paperclip size={16} className="text-indigo-500" />
                 {t('studio.editor.attachments')}
              </h3>
              <div className="space-y-2">
                {activity.attachments.map(att => (
                   <div key={att.id} className="flex items-center justify-between p-3 rounded-xl border border-slate-100 bg-slate-50">
                      <div className="flex items-center gap-3 overflow-hidden">
                        <Paperclip size={14} className="text-slate-400 flex-shrink-0" />
                        <span className="text-xs font-medium text-slate-600 truncate">{att.name}</span>
                      </div>
                      <button onClick={() => onUpdate({ attachments: activity.attachments.filter(a => a.id !== att.id) })} className="text-slate-400 hover:text-red-500 transition-colors">
                        <Trash2 size={14} />
                      </button>
                   </div>
                ))}
                <div className="text-xs text-slate-400 italic text-center p-2">Drag & drop files here functionality pending...</div>
              </div>
           </div>

           <div className="space-y-4">
              <h3 className="text-sm font-bold text-slate-900 flex items-center gap-2">
                 <LinkIcon size={16} className="text-indigo-500" />
                 {t('studio.editor.citations')}
              </h3>
              <div className="space-y-2">
                {activity.citations.map(cit => (
                   <div key={cit.id} className="flex items-center justify-between p-3 rounded-xl border border-slate-100 bg-slate-50">
                      <div className="flex items-center gap-3 overflow-hidden">
                        <LinkIcon size={14} className="text-slate-400 flex-shrink-0" />
                        <span className="text-xs font-medium text-slate-600 truncate">{cit.text}</span>
                      </div>
                      <button onClick={() => onUpdate({ citations: activity.citations.filter(c => c.id !== cit.id) })} className="text-slate-400 hover:text-red-500 transition-colors">
                        <Trash2 size={14} />
                      </button>
                   </div>
                ))}
                <Button variant="ghost" size="sm" className="w-full border-2 border-dashed border-slate-200 text-slate-400 hover:border-indigo-200 hover:bg-indigo-50">
                   <Plus size={14} className="mr-1" /> {t('studio.editor.add_citation')}
                </Button>
              </div>
           </div>
        </div>
      </div>

      {showFormBuilder && (
        <FormBuilderModal 
          initialFields={activity.assignment_config?.form_fields}
          onClose={() => setShowFormBuilder(false)}
          onSave={(fields: FormField[]) => {
            onUpdate({ 
              assignment_config: { 
                ...activity.assignment_config, 
                submission_types: ['form'],
                group_assignment: false, 
                peer_review: false,
                form_fields: fields 
              } 
            });
            setShowFormBuilder(false);
          }}
        />
      )}

      {showChecklistBuilder && (
        <ChecklistBuilderModal 
          initialConfig={activity.assignment_config?.checklist_config}
          onClose={() => setShowChecklistBuilder(false)}
          onSave={(config) => {
            onUpdate({ 
              assignment_config: { 
                ...activity.assignment_config, 
                submission_types: ['checklist'],
                group_assignment: false, 
                peer_review: false,
                checklist_config: config 
              } 
            });
            setShowChecklistBuilder(false);
          }}
        />
      )}

      {showQuizBuilder && (
          <QuizBuilderModal
            isOpen={showQuizBuilder}
            onClose={() => setShowQuizBuilder(false)}
            initialQuestions={activity.quiz_config?.questions}
            initialConfig={activity.quiz_config}
            onSave={(questions, config) => {
                onUpdate({ quiz_config: { ...config, questions } });
                setShowQuizBuilder(false);
            }}
          />
      )}

      {showSurveyBuilder && (
          <SurveyBuilderModal
            isOpen={showSurveyBuilder}
            onClose={() => setShowSurveyBuilder(false)}
            initialQuestions={activity.survey_config?.questions}
            initialConfig={activity.survey_config}
            onSave={(questions, config) => {
                onUpdate({ survey_config: { ...config, questions } });
                setShowSurveyBuilder(false);
            }}
          />
      )}
    </div>
  );
};
