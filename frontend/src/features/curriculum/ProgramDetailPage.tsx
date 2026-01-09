import React, { useState, useEffect, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { 
  ArrowLeft, Calendar, Users, BarChart3, Settings, 
  ChevronRight, BookOpen, Clock, AlertCircle, CheckCircle2,
  FileText, ExternalLink, GraduationCap, Edit, Plus,
  MoreVertical, Filter, Download, LayoutTemplate, Layers,
  Edit2, Lock, Loader2, Sparkles, ChevronDown, ChevronUp,
  FileCheck, Upload, GitMerge, Hourglass, Trophy, FormInput
} from 'lucide-react';
import { Button, Badge, Card, IconButton, Tabs } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { getProgram, getProgramVersionNodes, getProgramVersionMap, updateProgram, deleteProgram } from './api';
import { ProgramModal } from './components/ProgramModal';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// Helper for localization
const parseLocalized = (val: any, lang: string): string => {
  if (!val) return '';
  if (typeof val === 'string') {
    try {
        if (val.startsWith('{')) {
            const parsed = JSON.parse(val);
            return parsed[lang] || parsed.en || parsed.ru || val;
        }
        return val;
    } catch { return val; }
  }
  if (typeof val === 'object') {
     return val[lang] || val.en || val.ru || JSON.stringify(val);
  }
  return String(val);
};

// Node type icons
const NODE_ICONS: Record<string, any> = {
  form: FormInput,
  upload: Upload,
  confirmTask: FileCheck,
  decision: GitMerge,
  meeting: Users,
  waiting: Hourglass,
  external: ExternalLink,
  milestone: Trophy,
  default: Layers
};

// Phase colors
const PHASE_COLORS = [
  { bg: 'bg-sky-50', border: 'border-sky-200', text: 'text-sky-700', accent: 'bg-sky-500' },
  { bg: 'bg-emerald-50', border: 'border-emerald-200', text: 'text-emerald-700', accent: 'bg-emerald-500' },
  { bg: 'bg-amber-50', border: 'border-amber-200', text: 'text-amber-700', accent: 'bg-amber-500' },
  { bg: 'bg-red-50', border: 'border-red-200', text: 'text-red-700', accent: 'bg-red-500' },
  { bg: 'bg-purple-50', border: 'border-purple-200', text: 'text-purple-700', accent: 'bg-purple-500' },
  { bg: 'bg-teal-50', border: 'border-teal-200', text: 'text-teal-700', accent: 'bg-teal-500' },
];

const StatCard = ({ label, value, sub, icon: Icon, color }: any) => (
  <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-4">
     <div className={cn("p-3 rounded-xl", color)}>
        <Icon size={24} />
     </div>
     <div>
        <div className="text-xs font-bold text-slate-400 uppercase tracking-wider">{label}</div>
        <div className="text-2xl font-black text-slate-900 leading-none mt-1">{value}</div>
        {sub && <div className="text-[10px] font-bold text-slate-500 mt-1">{sub}</div>}
     </div>
  </div>
);

// Journey Node Card
const JourneyNodeCard = ({ node, index, lang, colorScheme }: any) => {
  const Icon = NODE_ICONS[node.type] || NODE_ICONS.default;
  const title = parseLocalized(node.title, lang);
  const [isExpanded, setIsExpanded] = useState(false);
  
  // Parse config for details
  let config: any = {};
  try {
    if (node.config) {
      config = typeof node.config === 'string' ? JSON.parse(node.config) : node.config;
    }
  } catch {}

  const hasUploads = config.requirements?.uploads?.length > 0;
  const hasFields = config.requirements?.fields?.length > 0;

  return (
    <div className={cn(
      "relative bg-white rounded-xl border shadow-sm transition-all duration-200 group",
      colorScheme.border,
      "hover:shadow-md hover:border-indigo-300"
    )}>
      <div className="p-4 flex items-start gap-4">
        {/* Step Number */}
        <div className={cn(
          "w-8 h-8 rounded-lg flex items-center justify-center text-xs font-bold shrink-0",
          colorScheme.bg, colorScheme.text
        )}>
          {index + 1}
        </div>
        
        {/* Icon */}
        <div className={cn(
          "w-10 h-10 rounded-xl flex items-center justify-center shrink-0",
          colorScheme.bg, colorScheme.text
        )}>
          <Icon size={20} />
        </div>
        
        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="secondary" className="text-[10px] uppercase">
              {node.type}
            </Badge>
            {node.slug && (
              <span className="text-[10px] font-mono text-slate-400">{node.slug}</span>
            )}
          </div>
          <h4 className="font-bold text-slate-900 text-sm line-clamp-2">{title}</h4>
          
          {/* Quick Info */}
          <div className="flex items-center gap-3 mt-2 text-xs text-slate-500">
            {hasUploads && (
              <span className="flex items-center gap-1">
                <Upload size={12} /> {config.requirements.uploads.length} uploads
              </span>
            )}
            {hasFields && (
              <span className="flex items-center gap-1">
                <FormInput size={12} /> {config.requirements.fields.length} fields
              </span>
            )}
            {node.prerequisites?.length > 0 && (
              <span className="flex items-center gap-1">
                <GitMerge size={12} /> {node.prerequisites.length} prerequisites
              </span>
            )}
          </div>
        </div>
        
        {/* Expand Button */}
        <button 
          onClick={() => setIsExpanded(!isExpanded)}
          className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-50 rounded-lg transition-colors"
        >
          {isExpanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
        </button>
      </div>
      
      {/* Expanded Details */}
      {isExpanded && (
        <div className="px-4 pb-4 pt-2 border-t border-slate-100 animate-in slide-in-from-top-2">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
            {/* Uploads */}
            {hasUploads && (
              <div>
                <h5 className="font-bold text-slate-700 mb-2 flex items-center gap-1">
                  <Upload size={12} /> Required Uploads
                </h5>
                <ul className="space-y-1">
                  {config.requirements.uploads.map((u: any, i: number) => (
                    <li key={i} className="text-slate-600 flex items-start gap-1">
                      <span className="text-slate-400">•</span>
                      {parseLocalized(u.label, lang) || u.key}
                      {u.required && <span className="text-red-500">*</span>}
                    </li>
                  ))}
                </ul>
              </div>
            )}
            
            {/* Form Fields */}
            {hasFields && (
              <div>
                <h5 className="font-bold text-slate-700 mb-2 flex items-center gap-1">
                  <FormInput size={12} /> Form Fields
                </h5>
                <ul className="space-y-1">
                  {config.requirements.fields.slice(0, 5).map((f: any, i: number) => (
                    <li key={i} className="text-slate-600 flex items-start gap-1">
                      <span className="text-slate-400">•</span>
                      {parseLocalized(f.label, lang) || f.key}
                      {f.required && <span className="text-red-500">*</span>}
                    </li>
                  ))}
                  {config.requirements.fields.length > 5 && (
                    <li className="text-slate-400">+{config.requirements.fields.length - 5} more...</li>
                  )}
                </ul>
              </div>
            )}
            
            {/* Prerequisites */}
            {node.prerequisites?.length > 0 && (
              <div>
                <h5 className="font-bold text-slate-700 mb-2 flex items-center gap-1">
                  <GitMerge size={12} /> Prerequisites
                </h5>
                <div className="flex flex-wrap gap-1">
                  {node.prerequisites.map((p: string, i: number) => (
                    <span key={i} className="px-2 py-0.5 bg-slate-100 rounded text-slate-600 font-mono text-[10px]">
                      {p}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

// Phase Section
const PhaseSection = ({ phase, nodes, lang, colorIndex }: any) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const colorScheme = PHASE_COLORS[colorIndex % PHASE_COLORS.length];
  const phaseTitle = parseLocalized(phase.title, lang) || `Phase ${phase.id}`;

  return (
    <div className="mb-8">
      {/* Phase Header */}
      <button 
        onClick={() => setIsCollapsed(!isCollapsed)}
        className={cn(
          "w-full flex items-center gap-4 p-4 rounded-2xl transition-all",
          colorScheme.bg, colorScheme.border, "border",
          "hover:shadow-md"
        )}
      >
        <div className={cn("w-2 h-12 rounded-full", colorScheme.accent)} />
        <div className="flex-1 text-left">
          <h3 className={cn("text-lg font-black", colorScheme.text)}>{phaseTitle}</h3>
          <p className="text-sm text-slate-500">{nodes.length} steps</p>
        </div>
        <Badge className={cn(colorScheme.bg, colorScheme.text, "border", colorScheme.border)}>
          {phase.id}
        </Badge>
        {isCollapsed ? <ChevronRight size={20} className="text-slate-400" /> : <ChevronDown size={20} className="text-slate-400" />}
      </button>
      
      {/* Nodes */}
      {!isCollapsed && (
        <div className="mt-4 space-y-3 pl-4 border-l-2 border-slate-200 ml-5">
          {nodes.map((node: any, i: number) => (
            <JourneyNodeCard 
              key={node.id} 
              node={node} 
              index={i} 
              lang={lang}
              colorScheme={colorScheme}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export const ProgramDetailPage: React.FC = () => {
  const { t, i18n } = useTranslation('common');
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState('Overview');
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const queryClient = useQueryClient();

  const { data: program, isLoading: pLoading, isError } = useQuery({
    queryKey: ['curriculum', 'programs', id],
    queryFn: () => getProgram(id!),
    enabled: !!id,
  });

  // Mutations
  const updateMutation = useMutation({
    mutationFn: (data: any) => updateProgram(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'programs', id] });
      toast.success('Program updated successfully');
      setIsEditModalOpen(false);
    },
    onError: () => toast.error('Failed to update program')
  });

  const deleteMutation = useMutation({
    mutationFn: () => deleteProgram(id!),
    onSuccess: () => {
      toast.success('Program deleted');
      navigate('/admin/programs');
    },
    onError: () => toast.error('Failed to delete program')
  });

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this program? This action cannot be undone.')) {
      deleteMutation.mutate();
    }
  };

  // Fetch Journey Map (version metadata)
  const { data: journeyMap } = useQuery({
    queryKey: ['curriculum', 'programs', id, 'map'],
    queryFn: () => getProgramVersionMap(id!),
    enabled: !!id,
  });

  // Fetch Nodes
  const { data: fetchedNodes = [], isLoading: nLoading } = useQuery({
    queryKey: ['curriculum', 'programs', id, 'nodes'],
    queryFn: () => getProgramVersionNodes(id!),
    enabled: !!id,
  });

  // Process nodes into phases
  const { phases, nodesByPhase } = useMemo(() => {
    // nodesData could be a direct array or part of journeyMap.nodes
    const nodes = Array.isArray(fetchedNodes) && fetchedNodes.length > 0 
        ? fetchedNodes 
        : ((journeyMap as any)?.nodes || []);
    
    // Group by module_key
    const grouped: Record<string, any[]> = {};
    nodes.forEach((node: any) => {
      const key = node.module_key || 'default';
      if (!grouped[key]) grouped[key] = [];
      grouped[key].push(node);
    });
    
    // Create phase info from journeyMap.phases or default
    let phases: any[] = (journeyMap as any)?.phases || [];
    
    // If no phases in config, create default from module_keys
    if (phases.length === 0) {
      phases = Object.keys(grouped).sort().map((key, i) => ({
        id: key,
        title: { en: `Phase ${key}`, ru: `Фаза ${key}` },
        order: i
      }));
    }
    
    return { phases, nodesByPhase: grouped };
  }, [fetchedNodes, journeyMap]);

  const totalNodes = useMemo(() => {
    return Object.values(nodesByPhase).flat().length;
  }, [nodesByPhase]);

  if (pLoading || nLoading) {
    return <div className="flex h-96 items-center justify-center"><Loader2 className="animate-spin text-indigo-600" /></div>;
  }

  if (isError || !program) {
    return <div className="p-8 text-center text-red-500">Failed to load program</div>;
  }

  const title = parseLocalized(program.title, i18n.language);
  const description = parseLocalized(program.description, i18n.language);

  /* ---------------- Tabs Content ---------------- */

  const OverviewTab = () => (
    <div className="space-y-8 animate-in fade-in duration-300">
       <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <StatCard label="Total Phases" value={phases.length} icon={Layers} color="bg-indigo-50 text-indigo-600" />
          <StatCard label="Journey Steps" value={totalNodes} icon={FileCheck} color="bg-blue-50 text-blue-600" />
          <StatCard label="Form Tasks" value={Object.values(nodesByPhase).flat().filter((n: any) => n.type === 'form').length} icon={FormInput} color="bg-emerald-50 text-emerald-600" />
          <StatCard label="Uploads Required" value={Object.values(nodesByPhase).flat().filter((n: any) => n.type === 'upload' || n.type === 'confirmTask').length} icon={Upload} color="bg-amber-50 text-amber-600" />
       </div>

       <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
             <div className="flex justify-between items-center mb-6">
                <h3 className="font-bold text-slate-900">Phase Distribution</h3>
             </div>
             <div className="space-y-4">
               {phases.map((phase, i) => {
                 const count = (nodesByPhase[phase.id] || []).length;
                 const pct = totalNodes > 0 ? (count / totalNodes) * 100 : 0;
                 const colorScheme = PHASE_COLORS[i % PHASE_COLORS.length];
                 return (
                   <div key={phase.id} className="flex items-center gap-4">
                     <div className="w-32 text-sm font-medium text-slate-700 truncate">
                       {parseLocalized(phase.title, i18n.language)}
                     </div>
                     <div className="flex-1 h-4 bg-slate-100 rounded-full overflow-hidden">
                       <div 
                         className={cn("h-full rounded-full", colorScheme.accent)}
                         style={{ width: `${pct}%` }}
                       />
                     </div>
                     <div className="w-16 text-right text-sm font-bold text-slate-600">
                       {count} steps
                     </div>
                   </div>
                 );
               })}
             </div>
          </div>
          
          <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm flex flex-col">
             <h3 className="font-bold text-slate-900 mb-4">Quick Actions</h3>
             <div className="space-y-3 flex-1">
                <Button variant="primary" className="w-full justify-start" onClick={() => navigate(`/admin/studio/programs/${id}/builder`)}>
                   <Edit size={16} className="mr-2" /> Open Journey Builder
                </Button>
                <Button variant="secondary" className="w-full justify-start">
                   <Download size={16} className="mr-2" /> Export Playbook
                </Button>
                <Button variant="secondary" className="w-full justify-start">
                   <Users size={16} className="mr-2" /> View Enrollments
                </Button>
             </div>
          </div>
       </div>
    </div>
  );

  const JourneyMapTab = () => (
    <div className="animate-in fade-in duration-300">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h3 className="text-xl font-black text-slate-900">Journey Structure</h3>
          <p className="text-sm text-slate-500 mt-1">
            {phases.length} phases • {totalNodes} steps total
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" size="sm" onClick={() => navigate(`/admin/studio/programs/${id}/builder`)}>
            <Edit size={14} className="mr-2" /> Edit in Builder
          </Button>
        </div>
      </div>

      {phases.length === 0 ? (
        <div className="p-12 text-center border-2 border-dashed border-slate-200 rounded-2xl text-slate-400 bg-slate-50/50">
          <Layers size={48} className="mx-auto mb-4 text-slate-300" />
          <h4 className="font-bold text-slate-600">No journey structure defined</h4>
          <p className="text-sm mt-1 mb-4">This program has no phases or steps yet.</p>
          <Button size="sm" onClick={() => navigate(`/admin/studio/programs/${id}/builder`)}>Open Journey Builder</Button>
        </div>
      ) : (
        <div>
          {phases.map((phase, i) => (
            <PhaseSection 
              key={phase.id}
              phase={phase}
              nodes={nodesByPhase[phase.id] || []}
              lang={i18n.language}
              colorIndex={i}
            />
          ))}
        </div>
      )}
    </div>
  );

  const EnrollmentsTab = () => (
    <div className="p-12 text-center text-slate-400 bg-slate-50 rounded-2xl border border-dashed border-slate-200">
       <Users size={48} className="mx-auto mb-4 opacity-20" />
       <h3 className="font-bold text-slate-600">No Active Enrollments</h3>
       <p className="text-sm mt-2">Start a new cohort to enroll students.</p>
       <Button className="mt-6" variant="primary" icon={Plus}>Start Cohort</Button>
    </div>
  );

  return (
    <div className="max-w-7xl mx-auto space-y-8 animate-in fade-in duration-500 pb-20">
      {/* Header */}
      <div>
         <button onClick={() => navigate('/admin/programs')} className="text-slate-400 hover:text-slate-700 flex items-center gap-2 mb-4 transition-colors font-medium text-sm">
            <ArrowLeft size={16} /> Back to Programs
         </button>
         <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div>
               <div className="flex items-center gap-3 mb-2">
                 <Badge variant={program.is_active ? 'success' : 'secondary'}>{program.is_active ? 'Active' : 'Inactive'}</Badge>
                 <span className="text-xs font-bold text-slate-400 tracking-wider uppercase">{program.code || 'NO-CODE'}</span>
                 {journeyMap?.version && (
                   <Badge variant="secondary">v{journeyMap.version}</Badge>
                 )}
               </div>
               <h1 className="text-4xl font-black text-slate-900 tracking-tight mb-2">{title || 'Program'}</h1>
               <p className="text-lg text-slate-500 max-w-2xl leading-relaxed">{description}</p>
            </div>
            <div className="flex gap-3">
               <IconButton icon={Settings} onClick={() => setIsEditModalOpen(true)} className="bg-slate-50 border border-slate-200" />
               <Button variant="secondary" icon={ExternalLink}>Preview Portal</Button>
               <Button icon={Edit} onClick={() => navigate(`/admin/studio/programs/${id}/builder`)}>Edit Journey</Button>
            </div>
         </div>
      </div>

      {/* Program Metadata Modal */}
      <ProgramModal 
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        initialData={program}
        onSave={(data) => updateMutation.mutate(data)}
        onDelete={handleDelete}
        isLoading={updateMutation.isPending || deleteMutation.isPending}
      />

      {/* Navigation Tabs */}
      <Tabs 
         tabs={['Overview', 'Program Structure', 'Enrollments']} 
         activeTab={activeTab} 
         onChange={setActiveTab} 
      />

      {/* Content Area */}
      <div className="min-h-[400px]">
         {activeTab === 'Overview' && <OverviewTab />}
         {activeTab === 'Program Structure' && <JourneyMapTab />}
         {activeTab === 'Enrollments' && <EnrollmentsTab />}
      </div>
    </div>
  );
};
