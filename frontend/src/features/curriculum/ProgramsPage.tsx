import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Plus, Search, MoreVertical, FileText, Clock, Users, 
  GraduationCap, BookOpen, Layers, ArrowUpDown, LayoutTemplate, 
  CheckCircle2, User as UserIcon, BarChart3, Loader2, Edit2, Trash2
} from 'lucide-react';
import { toast } from 'sonner';
import { Button, Badge, Card, IconButton } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { getPrograms, getCourses, createProgram, createCourse, updateProgram, deleteProgram, updateCourse, deleteCourse } from './api';
import { Program, Course } from './types';
import { ProgramModal } from './components/ProgramModal';
import { CourseModal } from './components/CourseModal';

// Templates Data
const TEMPLATES = [
  { id: 't1', title: 'PhD Onboarding', icon: GraduationCap, color: 'bg-indigo-50 text-indigo-600', description: 'Ethics, research, and defense' },
  { id: 't2', title: 'ENT Prep', icon: BookOpen, color: 'bg-emerald-50 text-emerald-600', description: 'Medical entrance exam flow' },
  { id: 't3', title: 'School Semester', icon: Layers, color: 'bg-orange-50 text-orange-600', description: 'Weekly modular structure' },
];

type ItemType = 'program' | 'course';

interface LibraryItem {
  id: string;
  type: ItemType;
  title: string;
  status: 'published' | 'draft' | 'archived' | 'active';
  students: number;
  updated: string;
  updatedAt: Date;
  owner: string;
  itemsCount?: number;
  description: string;
  originalId: string;
}

// Helper to parse localized strings
const parseLocalized = (val: any, lang: string): string => {
    if (!val) return '';
    if (typeof val === 'string') {
        try {
            if (val.startsWith('{')) {
                const parsed = JSON.parse(val);
                return parsed[lang] || parsed.en || val;
            }
            return val;
        } catch { return val; }
    }
    if (typeof val === 'object') {
       return val[lang] || val.en || JSON.stringify(val);
    }
    return String(val);
};

interface ProgramsPageProps {
  initialTab?: 'all' | 'program' | 'course';
}

export const ProgramsPage: React.FC<ProgramsPageProps> = ({ initialTab = 'program' }) => {
  const { i18n } = useTranslation('common');
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'all' | 'program' | 'course'>(initialTab);
  const [searchTerm, setSearchTerm] = useState('');
  
  const [filterOwner, setFilterOwner] = useState<string>('all');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [sortBy, setSortBy] = useState<'updated' | 'title'>('updated');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  const [isProgramModalOpen, setIsProgramModalOpen] = useState(false);
  const [isCourseModalOpen, setIsCourseModalOpen] = useState(false);
  const [editingProgram, setEditingProgram] = useState<Program | null>(null);
  const [editingCourse, setEditingCourse] = useState<Course | null>(null);
  
  const queryClient = useQueryClient();

  // Load Data
  const { data: programs = [], isLoading: pLoading } = useQuery({ queryKey: ['curriculum', 'programs'], queryFn: getPrograms });
  const { data: courses = [], isLoading: cLoading } = useQuery({ queryKey: ['curriculum', 'courses'], queryFn: () => getCourses() });

  // Mutations
  const createProgramMutation = useMutation({
    mutationFn: createProgram,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'programs'] });
      toast.success('Program created successfully');
      setIsProgramModalOpen(false);
    },
    onError: () => toast.error('Failed to create program')
  });

  const updateProgramMutation = useMutation({
    mutationFn: (data: any) => updateProgram(editingProgram!.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'programs'] });
      toast.success('Program updated');
      setEditingProgram(null);
    },
    onError: () => toast.error('Failed to update program')
  });

  const deleteProgramMutation = useMutation({
    mutationFn: (id: string) => deleteProgram(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'programs'] });
      toast.success('Program deleted');
    },
    onError: () => toast.error('Failed to delete program')
  });

  const createCourseMutation = useMutation({
    mutationFn: createCourse,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'courses'] });
      toast.success('Course registered successfully');
      setIsCourseModalOpen(false);
    },
    onError: () => toast.error('Failed to create course')
  });

  const updateCourseMutation = useMutation({
    mutationFn: (data: any) => updateCourse(editingCourse!.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'courses'] });
      toast.success('Course updated');
      setEditingCourse(null);
    },
    onError: () => toast.error('Failed to update course')
  });

  const deleteCourseMutation = useMutation({
    mutationFn: (id: string) => deleteCourse(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['curriculum', 'courses'] });
      toast.success('Course deleted');
    },
    onError: () => toast.error('Failed to delete course')
  });

  const handleDelete = (item: LibraryItem) => {
    if (window.confirm(`Are you sure you want to delete this ${item.type}?`)) {
      if (item.type === 'program') deleteProgramMutation.mutate(item.originalId);
      else deleteCourseMutation.mutate(item.originalId);
    }
  };

  const handleEdit = (item: LibraryItem) => {
      if (item.type === 'program') {
          const p = programs.find((p: Program) => p.id === item.originalId);
          if (p) setEditingProgram(p);
      } else {
          const c = courses.find((c: Course) => c.id === item.originalId);
          if (c) setEditingCourse(c);
      }
  };

  const items = useMemo(() => {
    // ... rest of useMemo logic unchanged
    // Map Programs
    const pItems: LibraryItem[] = programs.map((p: Program) => ({
        id: p.id,
        type: 'program',
        title: parseLocalized(p.title, i18n.language),
        status: p.status === 'active' ? 'published' : (p.status as any),
        students: 0, // Placeholder
        updated: new Date(p.updated_at || new Date()).toLocaleDateString(),
        updatedAt: new Date(p.updated_at || new Date()),
        owner: 'Admin',
        itemsCount: 0, 
        description: parseLocalized(p.description, i18n.language),
        originalId: p.id
    }));

    // Map Courses
    const cItems: LibraryItem[] = courses.map((c: Course) => ({
        id: c.id,
        type: 'course',
        title: parseLocalized(c.title, i18n.language),
        status: c.is_active ? 'published' : 'archived',
        students: 0,
        updated: new Date(c.updated_at || new Date()).toLocaleDateString(),
        updatedAt: new Date(c.updated_at || new Date()),
        owner: 'Admin',
        description: `${c.code} • ${c.credits} Credits${c.workload_hours ? ` • ${c.workload_hours}h` : ''}`,
        originalId: c.id
    }));

    return [...pItems, ...cItems];
  }, [programs, courses, i18n.language]);

  const owners = useMemo(() => Array.from(new Set(items.map(i => i.owner))), [items]);

  const filteredAndSortedItems = useMemo(() => {
    let result = items.filter(item => {
      const matchesTab = activeTab === 'all' || item.type === activeTab;
      const matchesSearch = item.title.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesOwner = filterOwner === 'all' || item.owner === filterOwner;
      const matchesStatus = filterStatus === 'all' || item.status === filterStatus || 
                            (filterStatus === 'published' && (item.status === 'active' || item.status === 'published'));
      return matchesTab && matchesSearch && matchesOwner && matchesStatus;
    });

    result.sort((a, b) => {
      if (sortBy === 'updated') {
        return sortOrder === 'desc' 
          ? b.updatedAt.getTime() - a.updatedAt.getTime()
          : a.updatedAt.getTime() - b.updatedAt.getTime();
      } else {
        return sortOrder === 'desc'
          ? b.title.localeCompare(a.title)
          : a.title.localeCompare(b.title);
      }
    });

    return result;
  }, [items, activeTab, searchTerm, filterOwner, filterStatus, sortBy, sortOrder]);

  const handleCreate = () => {
    if (activeTab === 'course') {
      setIsCourseModalOpen(true);
    } else {
      setIsProgramModalOpen(true);
    }
  };

  const handleCreateFromTemplate = (tpl: any) => {
      console.log("Create from template", tpl);
  };

  if (pLoading || cLoading) {
      return <div className="flex h-96 items-center justify-center"><Loader2 className="animate-spin text-indigo-600" /></div>;
  }

  return (
    <div className="max-w-7xl mx-auto space-y-10 animate-in fade-in duration-500">
      
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-6">
        <div>
          <h1 className="text-3xl font-black text-slate-900 tracking-tight">
            {activeTab === 'program' ? 'Programs' : activeTab === 'course' ? 'Courses' : 'Asset Library'}
          </h1>
          <p className="text-slate-500 mt-2 text-lg">Manage and scale your educational strategy.</p>
        </div>
        <div className="flex gap-3">
          <Button variant="secondary" icon={FileText}>Import Data</Button>
          <Button icon={Plus} onClick={handleCreate}>
            Create Item
          </Button>
        </div>
      </div>

      {/* Quick Templates Panel */}
      <section className="space-y-4">
        <div className="flex items-center gap-2 text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">
          <LayoutTemplate size={14} className="text-indigo-500" /> Quick Templates
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {TEMPLATES.map(t => (
            <div key={t.id} className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-4 hover:border-indigo-300 transition-all group">
              <div className={cn("w-12 h-12 rounded-xl flex items-center justify-center transition-transform group-hover:scale-110", t.color)}>
                <t.icon size={24} />
              </div>
              <div className="flex-1 min-w-0">
                <h4 className="font-bold text-slate-900 text-sm truncate">{t.title}</h4>
                <p className="text-xs text-slate-400 truncate font-medium">{t.description}</p>
              </div>
              <IconButton 
                icon={Plus} 
                className="bg-slate-50 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50" 
                onClick={() => handleCreateFromTemplate(t)}
                title="Create from template"
              />
            </div>
          ))}
        </div>
      </section>

      {/* Controls Bar */}
      <div className="space-y-4">
        <div className="bg-white p-2 rounded-2xl border border-slate-200 shadow-sm flex flex-col xl:flex-row gap-4">
          {/* Main Search & Tabs */}
          <div className="flex flex-col sm:flex-row flex-1 gap-2">
            <div className="flex p-1 bg-slate-100 rounded-xl flex-shrink-0">
              {(['all', 'program', 'course'] as const).map(tab => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={cn(
                    "px-4 py-2 rounded-lg text-sm font-bold transition-all capitalize",
                    activeTab === tab ? "bg-white text-slate-900 shadow-sm" : "text-slate-500 hover:text-slate-900"
                  )}
                >
                  {tab === 'all' ? 'All Items' : tab + 's'}
                </button>
              ))}
            </div>

            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={16} />
              <input 
                type="text" 
                placeholder="Search by title..." 
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full h-full pl-10 pr-4 rounded-xl bg-slate-50 border-none focus:bg-white focus:ring-2 focus:ring-indigo-100 text-sm font-medium transition-all"
              />
            </div>
          </div>

          {/* Advanced Filters */}
          <div className="flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-1.5 bg-slate-50 rounded-xl border border-slate-100">
              <UserIcon size={14} className="text-slate-400" />
              <select 
                value={filterOwner}
                onChange={(e) => setFilterOwner(e.target.value)}
                className="bg-transparent text-xs font-bold text-slate-600 outline-none cursor-pointer"
              >
                <option value="all">All Owners</option>
                {owners.map(o => <option key={o} value={o}>{o}</option>)}
              </select>
            </div>

            <div className="flex items-center gap-2 px-3 py-1.5 bg-slate-50 rounded-xl border border-slate-100">
              <CheckCircle2 size={14} className="text-slate-400" />
              <select 
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
                className="bg-transparent text-xs font-bold text-slate-600 outline-none cursor-pointer"
              >
                <option value="all">All Status</option>
                <option value="published">Published</option>
                <option value="draft">Drafts</option>
                <option value="archived">Archived</option>
              </select>
            </div>

            <div className="h-8 w-px bg-slate-200 mx-1 hidden xl:block" />

            {/* Sorting */}
            <div className="flex items-center bg-slate-100 p-1 rounded-xl">
               <button 
                onClick={() => {
                  if (sortBy === 'updated') setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  else setSortBy('updated');
                }}
                className={cn(
                  "px-3 py-1.5 rounded-lg text-xs font-bold flex items-center gap-1.5 transition-all",
                  sortBy === 'updated' ? "bg-white shadow-sm text-slate-900" : "text-slate-500 hover:text-slate-700"
                )}
               >
                 <Clock size={14} /> 
                 {sortBy === 'updated' && <ArrowUpDown size={10} />}
               </button>
               <button 
                onClick={() => {
                  if (sortBy === 'title') setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
                  else setSortBy('title');
                }}
                className={cn(
                  "px-3 py-1.5 rounded-lg text-xs font-bold flex items-center gap-1.5 transition-all",
                  sortBy === 'title' ? "bg-white shadow-sm text-slate-900" : "text-slate-500 hover:text-slate-700"
                )}
               >
                 <FileText size={14} /> 
                 {sortBy === 'title' && <ArrowUpDown size={10} />}
               </button>
            </div>
          </div>
        </div>
        
        {/* Results Info */}
        <div className="flex justify-between items-center px-2">
           <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">
             Found {filteredAndSortedItems.length} educational units
           </span>
           <div className="flex items-center gap-4 text-[10px] font-bold text-slate-400">
             <span className="flex items-center gap-1"><div className="w-2 h-2 rounded-full bg-emerald-500" /> Published</span>
             <span className="flex items-center gap-1"><div className="w-2 h-2 rounded-full bg-amber-500" /> Draft</span>
             <span className="flex items-center gap-1"><div className="w-2 h-2 rounded-full bg-slate-400" /> Archived</span>
           </div>
        </div>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {filteredAndSortedItems.map(item => (
          <Card 
            key={item.id} 
            className={cn(
              "group relative overflow-hidden border-slate-200 hover:border-indigo-300 hover:shadow-lg transition-all duration-300 flex flex-col",
              item.status === 'archived' && "opacity-75"
            )}
          >
            {/* Top Stripe */}
            <div className={cn(
              "h-1.5 w-full",
              item.status === 'archived' ? 'bg-slate-300' : (item.type === 'program' ? 'bg-indigo-500' : 'bg-emerald-500')
            )} />
            
            <div className="p-6 flex-1 flex flex-col">
              <div className="flex justify-between items-start mb-4">
                 <div className={cn(
                   "w-12 h-12 rounded-2xl flex items-center justify-center text-xl shadow-sm",
                   item.type === 'program' ? "bg-indigo-50 text-indigo-600" : "bg-emerald-50 text-emerald-600"
                 )}>
                    {item.type === 'program' ? <Layers size={24} /> : <BookOpen size={24} />}
                 </div>
                 <div className="flex gap-2">
                   <Badge variant={item.status === 'published' ? 'success' : item.status === 'draft' ? 'warning' : 'secondary'}>
                     {item.status}
                   </Badge>
                   <IconButton icon={MoreVertical} className="w-8 h-8 -mr-2" />
                 </div>
              </div>

              <div className="mb-1 flex items-center gap-2">
                <span className="text-[10px] font-black uppercase tracking-wider text-slate-400">
                  {item.type}
                </span>
                <span className="text-[10px] font-bold text-slate-300">•</span>
                <span className="text-[10px] font-bold text-slate-400">BY {item.owner.toUpperCase()}</span>
              </div>
              
              <h3 className="text-xl font-bold text-slate-900 mb-2 group-hover:text-indigo-600 transition-colors line-clamp-1">
                {item.title}
              </h3>
              
              <p className="text-sm text-slate-500 mb-6 line-clamp-2 flex-1">
                {item.description}
              </p>

              {/* Metrics */}
              <div className="flex items-center gap-4 pt-4 border-t border-slate-100">
                <div className="flex items-center gap-1.5 text-xs font-medium text-slate-500">
                  <Users size={14} className="text-slate-400" />
                  {item.students} students
                </div>
                {(item.itemsCount !== undefined) && (
                  <div className="flex items-center gap-1.5 text-xs font-medium text-slate-500">
                    <BookOpen size={14} className="text-slate-400" />
                    {item.itemsCount} courses
                  </div>
                )}
                <div className="flex-1 text-right text-[10px] text-slate-400 font-bold uppercase">
                  {item.updated}
                </div>
              </div>
            </div>

            {/* Hover Action */}
            <div className={cn(
               "absolute inset-0 bg-white/60 backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-3 p-6 z-10",
               item.status === 'archived' && "pointer-events-none"
            )}>
               <div className="transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 flex flex-col gap-3 w-full max-w-[200px]">
                 {item.type === 'program' ? (
                    <>
                      <Button onClick={() => navigate(`/admin/studio/programs/${item.originalId}/builder`)} className="w-full shadow-lg">
                        Edit Program
                      </Button>
                      <Button variant="secondary" onClick={() => navigate(`/admin/programs/${item.originalId}`)} className="w-full bg-white flex items-center justify-center gap-2">
                        <BarChart3 size={16} /> View Ops
                      </Button>
                    </>
                 ) : (
                    <>
                      <Button onClick={() => navigate(`/admin/studio/courses/${item.originalId}/builder`)} className="w-full shadow-lg">
                        Open Course Builder
                      </Button>
                      <Button variant="secondary" onClick={() => navigate(`/admin/studio/courses/${item.originalId}/preview`)} className="w-full bg-white">
                        Simulation Mode
                      </Button>
                    </>
                 )}
                 
                 <div className="h-px bg-slate-200 my-1 opacity-50" />
                 
                 <Button variant="secondary" onClick={() => handleEdit(item)} className="w-full bg-white flex items-center justify-center gap-2">
                   <Edit2 size={16} /> Edit Metadata
                 </Button>
                 <Button variant="secondary" onClick={() => handleDelete(item)} className="w-full bg-white text-red-500 hover:text-red-700 hover:bg-red-50 flex items-center justify-center gap-2">
                   <Trash2 size={16} /> Delete Unit
                 </Button>
               </div>
            </div>
          </Card>
        ))}

        {/* Empty State */}
        {filteredAndSortedItems.length === 0 && (
          <div className="col-span-full py-20 bg-slate-50 border-2 border-dashed border-slate-200 rounded-[3rem] flex flex-col items-center justify-center text-slate-400">
             <div className="w-20 h-20 bg-white rounded-full flex items-center justify-center shadow-sm mb-6">
                <Search size={32} strokeWidth={1.5} />
             </div>
             <h3 className="text-xl font-bold text-slate-800">No matching assets</h3>
             <p className="mt-2 text-slate-500">Try adjusting your filters or create a new unit.</p>
             <Button className="mt-8" icon={Plus} onClick={handleCreate}>New Authoring Session</Button>
          </div>
        )}
      </div>

      <ProgramModal 
        isOpen={isProgramModalOpen || !!editingProgram}
        onClose={() => { setIsProgramModalOpen(false); setEditingProgram(null); }}
        initialData={editingProgram}
        onSave={(data: any) => editingProgram ? updateProgramMutation.mutate(data) : createProgramMutation.mutate(data)}
        onDelete={() => editingProgram && handleDelete({ type: 'program', originalId: editingProgram.id } as any)}
        isLoading={createProgramMutation.isPending || updateProgramMutation.isPending || deleteProgramMutation.isPending}
      />

      <CourseModal 
        isOpen={isCourseModalOpen || !!editingCourse}
        onClose={() => { setIsCourseModalOpen(false); setEditingCourse(null); }}
        initialData={editingCourse}
        onSave={(data: any) => editingCourse ? updateCourseMutation.mutate(data) : createCourseMutation.mutate(data)}
        onDelete={() => editingCourse && handleDelete({ type: 'course', originalId: editingCourse.id } as any)}
        isLoading={createCourseMutation.isPending || updateCourseMutation.isPending || deleteCourseMutation.isPending}
      />

    </div>
  );
};
