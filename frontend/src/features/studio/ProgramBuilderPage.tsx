import React, { useState, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  getProgramVersionMap, getProgramVersionNodes, updateProgramVersionMap, 
  createProgramVersionNode, updateProgramVersionNode, deleteProgramVersionNode, deleteProgram 
} from '@/features/curriculum/api';
import { 
  FileText, Info, CheckSquare, Sparkles, Layout, Flag, Award, CreditCard, ClipboardList, 
  Stamp, Calendar, Undo, Redo, AlertCircle, Share, Minimize2, Maximize2, ZoomIn, ZoomOut, 
  GitMerge, GripVertical, Target, Map as MapIcon, Plus, Settings2, Trash2, ChevronRight, X, 
  Loader2, BookOpen, CheckCircle2, Zap, List, ExternalLink 
} from 'lucide-react';
import { toast } from 'sonner';
import { Reorder } from 'framer-motion';
import { Button, Input, IconButton, Switch } from '@/features/admin/components/AdminUI';
import { useHistory } from '@/features/admin/hooks/useHistory';
import { WorldSettingsModal } from '@/features/admin/components/WorldSettingsModal';
import { StepInspector } from './components/StepInspector';
import { JourneyMap } from '@/components/map/JourneyMap';
import { cn } from '@/lib/utils';
import { ProgramVersionNode, ProgramPhase, ProgramNodeType as NodeType, FlowEdge } from './types';



// Helper for localization (Moved to top for scope)
const parseLocalized = (val: any, lang: string = 'en'): string => {
  if (!val) return '';
  if (typeof val === 'string') {
    if (val.trim() === 'null') return '';
    try {
        if (val.trim().startsWith('{')) {
            const parsed = JSON.parse(val);
            return parsed[lang] || parsed.kk || parsed.kz || parsed.en || parsed.ru || val;
        }
        return val;
    } catch { return val; }
  }
  if (typeof val === 'object') {
     return val[lang] || val.kk || val.kz || val.en || val.ru || '';
  }
  return String(val);
};

// Helper to transform Builder State to Playbook for Preview
const builderStateToPlaybook = (worlds: World[], nodes: ProgramVersionNode[], edges: FlowEdge[]): any => {
  return {
    id: "preview-playbook",
    title: "Preview Program",
    worlds: worlds.map(w => ({
      id: w.id,
      title: parseLocalized(w.title),
      nodes: nodes
        .filter(n => n.module_key === w.id)
        .sort((a, b) => a.coordinates.y - b.coordinates.y) // Simple sort by Y for sequence
        .map(n => ({
          ...n.config, // Spread config first so specific props can override if needed, or vice versa. Usually specific props like 'id' are safer after.
          id: n.id,
          type: n.type,
          title: parseLocalized(n.title) || n.slug || "Untitled Step",
          status: "active", // Active state for better preview visibility
          description: parseLocalized(n.description) || "",
          points: n.points || 0,
          who_can_complete: ["student", "advisor", "admin"] // Allow all for preview
        }))
    }))
  };
};


// --- Types & Constants ---
// --- Types & Constants ---


interface JourneyBuilderProps {
  onNavigate?: (path: string) => void;
}

interface World {
  id: string;
  title: string;
  description?: string;
  order: number;
  color: string;
  x: number; // Added for drag support
  y: number; // Added for drag support
  condition?: {
    nodeId: string;
    fieldKey: string;
    operator: 'equals' | 'not_equals' | 'contains';
    value: string;
  } | null;
  collapsed?: boolean;
}

interface JourneyState {
  worlds: World[];
  nodes: ProgramVersionNode[];
  edges: FlowEdge[];
}

// --- Mock Initial Data ---
const INITIAL_WORLDS: World[] = [
  { id: 'w1', title: 'I — Preparation', order: 1, color: '#6366f1', description: 'Initial profile setup.', collapsed: false, condition: null, x: 50, y: 50 }, 
  { id: 'w2', title: 'II — Pre-examination', order: 2, color: '#10b981', collapsed: false, condition: null, x: 50, y: 450 }, 
  { id: 'w3', title: 'III — Defense', order: 3, color: '#f59e0b', collapsed: true, condition: null, x: 50, y: 850 } 
];

const INITIAL_NODES: ProgramVersionNode[] = [
  { id: 'n1', module_key: 'w1', coordinates: { x: 100, y: 120 }, title: "Student Profile", type: 'form', points: 50, config: { fields: [] }, slug: 'n1' },
  { id: 'n2', module_key: 'w1', coordinates: { x: 380, y: 120 }, title: "Research Methodology", type: 'course', points: 100, config: { courseId: 'c1' }, slug: 'n2' },
  { id: 'n3', module_key: 'w1', coordinates: { x: 660, y: 120 }, title: "Advisor Assignment", type: 'sync_ops', points: 0, config: { action: 'assign_advisor' }, slug: 'n3' },
  { id: 'n4', module_key: 'w2', coordinates: { x: 100, y: 520 }, title: "Anti-Plagiarism Fee", type: 'payment', points: 20, config: { amount: 5000, currency: 'KZT' }, slug: 'n4' },
  { id: 'n5', module_key: 'w2', coordinates: { x: 380, y: 520 }, title: "Department Approval", type: 'approval', points: 0, config: { role: 'advisor' }, slug: 'n5' },
];

const INITIAL_EDGES: FlowEdge[] = [
  { id: 'e1', from: 'n1', to: 'n2', type: 'solid' },
  { id: 'e2', from: 'n2', to: 'n3', type: 'solid' },
  { id: 'e3', from: 'n3', to: 'n4', type: 'dashed' }, // Phase transition
  { id: 'e4', from: 'n4', to: 'n5', type: 'solid' }
];

const generateId = (prefix: string) => `${prefix}_${Math.random().toString(36).substr(2, 5)}`;



  // --- Visual System ---

  const getNodeVisuals = (type: NodeType, t: any) => {
    switch(type) {
      case 'course': return { icon: BookOpen, color: 'text-purple-600', bg: 'bg-purple-50', border: 'border-purple-200', label: t('builder.nodeTypes.course', 'Learning Event') };
      case 'payment': return { icon: CreditCard, color: 'text-amber-600', bg: 'bg-amber-50', border: 'border-amber-200', label: t('builder.nodeTypes.payment', 'Financial Gate') };
      case 'sync_ops': return { icon: Sparkles, color: 'text-emerald-600', bg: 'bg-emerald-50', border: 'border-emerald-200', label: t('builder.nodeTypes.sync_ops', 'Ops Automation') };
      case 'approval': return { icon: Stamp, color: 'text-slate-600', bg: 'bg-slate-50', border: 'border-slate-300', label: t('builder.nodeTypes.approval', 'Admin Gate') };
      case 'form': return { icon: FileText, color: 'text-blue-600', bg: 'bg-blue-50', border: 'border-blue-200', label: t('builder.nodeTypes.form', 'Data Collection') };
      case 'meeting': return { icon: Calendar, color: 'text-pink-600', bg: 'bg-pink-50', border: 'border-pink-200', label: t('builder.nodeTypes.meeting', 'Sync Event') };
      case 'checklist': return { icon: CheckSquare, color: 'text-orange-600', bg: 'bg-orange-50', border: 'border-orange-200', label: t('builder.nodeTypes.checklist', 'Requirement') };
      case 'milestone': return { icon: Flag, color: 'text-indigo-600', bg: 'bg-indigo-50', border: 'border-indigo-200', label: t('builder.nodeTypes.milestone', 'Milestone') };
      case 'info': return { icon: Info, color: 'text-cyan-600', bg: 'bg-cyan-50', border: 'border-cyan-200', label: t('builder.nodeTypes.info', 'Information') };
      case 'survey': return { icon: ClipboardList, color: 'text-teal-600', bg: 'bg-teal-50', border: 'border-teal-200', label: t('builder.nodeTypes.survey', 'Feedback') };
      case 'confirmTask': return { icon: CheckCircle2, color: 'text-green-600', bg: 'bg-green-50', border: 'border-green-200', label: t('builder.nodeTypes.confirmTask', 'Confirmation') };
      default: return { icon: Layout, color: 'text-slate-500', bg: 'bg-white', border: 'border-slate-200', label: t('builder.nodeTypes.processStep', 'Process Step') };
    }
  };

const validateNode = (node: ProgramVersionNode): string[] => {
  const errors = [];
  if (!node.title.trim()) errors.push("Title is required");
  if (node.type === 'course' && !node.config?.courseId) errors.push("Linked Course is missing");
  if (node.type === 'payment' && !node.config?.amount) errors.push("Payment amount missing");
  return errors;
};

// --- Node Component (Map View) ---

const FlowNode = ({ id, title, type, active, onClick, onMouseDown, module_key, worlds, validationErrors, onAddNext, points }: any) => {
  const { t } = useTranslation();
  const world = worlds.find((w: World) => w.id === module_key);
  const visuals = getNodeVisuals(type as any, t);
  const Icon = visuals.icon;
  const hasErrors = validationErrors && validationErrors.length > 0;

  return (
    <div 
      onMouseDown={(e) => { e.stopPropagation(); onMouseDown(e, id); }}
      onClick={(e) => { e.stopPropagation(); onClick(id); }}
      className={cn(
        "absolute w-[260px] bg-white rounded-2xl shadow-sm flex flex-col cursor-grab active:cursor-grabbing transition-all hover:shadow-xl hover:-translate-y-1 z-10 group select-none",
        active ? "ring-2 ring-indigo-500 shadow-md" : "border border-slate-200",
        hasErrors && "ring-2 ring-red-400"
      )}
    >
      <div className="p-1">
         <div className={cn("rounded-xl p-3 flex items-start gap-3", visuals.bg)}>
            <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center bg-white shadow-sm border border-white/50", visuals.color)}>
               <Icon size={20} />
            </div>
            <div className="flex-1 min-w-0 pt-0.5">
               <div className="text-[10px] font-black uppercase tracking-wider opacity-60 mb-0.5">{visuals.label}</div>
               <div className="font-bold text-slate-900 text-sm leading-tight line-clamp-2">{parseLocalized(title) || "Untitled Step"}</div>
            </div>
         </div>
      </div>
      
      {/* Footer Info */}
      <div className="px-4 py-2 border-t border-slate-100 flex justify-between items-center text-[10px] text-slate-400 font-medium bg-white rounded-b-2xl">
         <span>{parseLocalized(world?.title).split('—')[0] || 'Unassigned'}</span>
         {points > 0 && <span className="flex items-center gap-1 text-amber-500 font-bold"><Award size={10} /> {points} XP</span>}
         {hasErrors ? (
            <span className="text-red-500 flex items-center gap-1 font-bold"><AlertCircle size={10} /> Invalid</span>
         ) : type === 'sync_ops' ? (
            <span className="flex items-center gap-1 text-emerald-600 font-bold"><Zap size={10} fill="currentColor" /> Auto-Trigger</span>
         ) : type === 'course' ? (
            <span className="flex items-center gap-1"><BookOpen size={10} /> Linked</span>
         ) : null}
      </div>

      {/* Connection Ports */}
      {active && (
        <>
          {/* Quick Add Action */}
          <div 
            className="absolute -right-3 top-1/2 -translate-y-1/2 w-6 h-6 bg-indigo-600 rounded-full flex items-center justify-center text-white shadow-lg cursor-pointer hover:scale-110 transition-transform z-20 border-2 border-white"
            onClick={(e) => { e.stopPropagation(); onAddNext(id); }}
            title="Add Next Step"
          >
             <Plus size={14} />
          </div>
        </>
      )}
    </div>
  );
};

// --- List View Component (Reorderable) ---
const ListView = ({ worlds, nodes, selectedNodeId, onSelectNode, onAddNode, onAddWorld, onReorderNodes, onEditWorld, onDeleteNode }: any) => {
  const { t } = useTranslation();
  return (
    <div className="flex-1 overflow-y-auto bg-slate-50 p-8 space-y-8">
      {worlds.map((world: World) => {
        // Get nodes for this world
        const worldNodes = nodes.filter((n: ProgramVersionNode) => n.module_key === world.id);
        
        return (
          <div key={world.id} className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
             <div className="p-4 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
                <div className="flex items-center gap-3">
                   <div className="w-3 h-3 rounded-full" style={{ backgroundColor: world.color }} />
                   <div>
                      <h3 className="font-bold text-slate-900 flex items-center gap-2">
                        {parseLocalized(world.title)}
                        {world.condition && (
                          <span className="bg-amber-100 text-amber-700 text-[10px] px-2 py-0.5 rounded-full flex items-center gap-1">
                            <GitMerge size={10} /> Conditional
                          </span>
                        )}
                      </h3>
                      {world.description && <p className="text-xs text-slate-500">{parseLocalized(world.description)}</p>}
                   </div>
                </div>
                <div className="flex items-center gap-2">
                   <IconButton icon={Settings2} size="sm" onClick={() => onEditWorld(world)} title="Configure Phase" />
                   <Button size="sm" variant="ghost" icon={Plus} onClick={() => onAddNode(world.id)}>Add Step</Button>
                </div>
             </div>
             
             {/* Draggable List */}
             <div className="divide-y divide-slate-100 bg-white">
                <Reorder.Group 
                  axis="y" 
                  values={worldNodes} 
                  onReorder={(newOrder) => onReorderNodes(world.id, newOrder)}
                >
                  {worldNodes.map((node: ProgramVersionNode) => {
                     const visuals = getNodeVisuals(node.type as any, t);
                     const Icon = visuals.icon;
                     return (
                        <Reorder.Item key={node.id} value={node}>
                          <div 
                            onClick={() => onSelectNode(node.id)}
                            className={cn(
                              "p-4 flex items-center gap-4 hover:bg-slate-50 transition-colors cursor-pointer relative bg-white group",
                              selectedNodeId === node.id ? "bg-indigo-50/50" : ""
                            )}
                          >
                             <div className="text-slate-300 cursor-grab active:cursor-grabbing hover:text-slate-500">
                               <GripVertical size={16} />
                             </div>
                             <div className={cn("p-2 rounded-lg flex-shrink-0", visuals.bg, visuals.color)}>
                                <Icon size={20} />
                             </div>
                             <div className="flex-1 min-w-0">
                                <div className="font-bold text-slate-900 truncate">{parseLocalized(node.title)}</div>
                                <div className="text-xs text-slate-500 truncate">{parseLocalized(node.description) || visuals.label}</div>
                             </div>
                             
                             {node.points && node.points > 0 && (
                               <div className="text-xs font-bold text-amber-500 bg-amber-50 px-2 py-1 rounded flex items-center gap-1">
                                 <Award size={12} /> {node.points} XP
                               </div>
                             )}
                             
                             <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                <IconButton 
                                  icon={Trash2} 
                                  size="sm" 
                                  className="text-slate-400 hover:text-red-500 hover:bg-red-50"
                                  onClick={(e: React.MouseEvent) => { e.stopPropagation(); onDeleteNode(node.id); }} 
                                />
                                <ChevronRight size={16} className="text-slate-300" />
                             </div>
                          </div>
                        </Reorder.Item>
                     );
                  })}
                </Reorder.Group>
                
                {worldNodes.length === 0 && (
                   <div className="p-8 text-center text-slate-400 text-sm italic">
                      No steps in this phase yet.
                   </div>
                )}
             </div>
          </div>
        );
      })}
      
      {/* Add Phase Button */}
      <button 
        onClick={onAddWorld}
        className="w-full py-4 border-2 border-dashed border-slate-300 rounded-2xl text-slate-400 font-bold hover:text-indigo-600 hover:border-indigo-300 hover:bg-indigo-50 transition-all flex items-center justify-center gap-2"
      >
         <Plus size={20} /> Add New Phase
      </button>
    </div>
  );
};

export const ProgramBuilderPage: React.FC<JourneyBuilderProps> = ({ onNavigate }) => {
  const { t } = useTranslation();
  const { programId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  
  // Use provided onNavigate or fall back to react-router's navigate
  const handleNavigate = (path: string) => {
    if (onNavigate) {
      onNavigate(path);
    } else {
      navigate(path);
    }
  };

  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [transform, setTransform] = useState({ x: 0, y: 0, scale: 1 });
  const [isPanning, setIsPanning] = useState(false);
  const [viewMode, setViewMode] = useState<'list' | 'map'>('map'); 
  const [isZenMode, setIsZenMode] = useState(false);
  const [editingWorld, setEditingWorld] = useState<World | null>(null);
  
  // Preview Mode State
  const [isPreview, setIsPreview] = useState(false);

  // --- Data Fetching ---
  const { data: mapData, isLoading: isLoadingMap } = useQuery({
    queryKey: ['programMap', programId],
    queryFn: () => getProgramVersionMap(programId!),
    enabled: !!programId
  });

  const { data: nodesData, isLoading: isLoadingNodes } = useQuery({
    queryKey: ['programNodes', programId],
    queryFn: () => getProgramVersionNodes(programId!),
    enabled: !!programId
  });

  // --- Mutations ---
  const createNodeMutation = useMutation({
    mutationFn: (data: any) => createProgramVersionNode(programId!, data),
    onSuccess: (newNode: any) => {
      queryClient.invalidateQueries({ queryKey: ['programNodes', programId] });
      toast.success("Step created");
      setSelectedNodeId(newNode.id);
    },
    onError: () => toast.error("Failed to create step")
  });

  const updateNodeMutation = useMutation({
    mutationFn: ({ id, data }: { id: string, data: any }) => updateProgramVersionNode(programId!, id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programNodes', programId] });
    },
    onError: () => toast.error("Failed to update step")
  });

  const deleteNodeMutation = useMutation({
    mutationFn: (nodeId: string) => deleteProgramVersionNode(programId!, nodeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programNodes', programId] });
      toast.success("Step removed");
      setSelectedNodeId(null);
    },
    onError: () => toast.error("Failed to delete step")
  });

  const updateMapMutation = useMutation({
    mutationFn: (data: any) => updateProgramVersionMap(programId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programMap', programId] });
      toast.success("Phase configuration saved");
    },
    onError: () => toast.error("Failed to save phases")
  });

  // History Management
  const { 
    state: journeyState, 
    set: setJourneyState,
    undo,
    redo,
    canUndo,
    canRedo
  } = useHistory<JourneyState>({
    worlds: INITIAL_WORLDS,
    nodes: INITIAL_NODES,
    edges: []
  });

  // Hydrate State
  useEffect(() => {
    if (mapData && nodesData) {
        console.log("Hydrating Builder State...", { mapData, nodesData });
        const safeMapData = mapData as any;
        
        let loadedWorlds: World[] = [];
        const rawWorlds = safeMapData.phases || []; 
        
        loadedWorlds = rawWorlds.map((w: any) => ({
            id: String(w.id || `w_${Math.random()}`),
            title: parseLocalized(w.title),
            order: Number(w.order) || 0,
            color: w.color || '#6366f1',
            description: parseLocalized(w.description),
            collapsed: !!w.collapsed,
            condition: w.condition,
            x: Number(w.position?.x || w.x) || 0, 
            y: Number(w.position?.y || w.y) || 0
        }));
        
        const rawNodes = Array.isArray(nodesData) ? nodesData : (safeMapData.nodes || []);
        
        const loadedNodes: ProgramVersionNode[] = (rawNodes as any[]).map(n => ({
            id: n.id,
            slug: n.slug || n.id,
            module_key: n.module_key || loadedWorlds[0]?.id || 'w1', 
            title: parseLocalized(n.title), 
            description: parseLocalized(n.description),
            type: n.type as NodeType,
            coordinates: typeof n.coordinates === 'string' ? JSON.parse(n.coordinates || '{}') : (n.coordinates || { x: 100, y: 100 }),
            config: typeof n.config === 'string' ? JSON.parse(n.config || '{}') : (n.config || {}),
            points: n.points || 0
        }));

        if (loadedWorlds.length === 0 && loadedNodes.length > 0) {
            const uniqueKeys = Array.from(new Set(loadedNodes.map(n => n.module_key)));
            loadedWorlds = uniqueKeys.map((key, idx) => ({
                id: String(key || `generated_${idx}`),
                title: `Phase ${key}`, 
                order: idx + 1,
                color: '#6366f1',
                x: 50,
                y: 50 + (idx * 400)
            }));
        }

        setJourneyState({
            worlds: loadedWorlds.length > 0 ? loadedWorlds : INITIAL_WORLDS,
            nodes: loadedNodes.length > 0 ? loadedNodes : INITIAL_NODES,
            edges: [] 
        });
    }
  }, [mapData, nodesData, setJourneyState]);

  const { worlds, nodes, edges } = journeyState;
  
  // Dragging State
  const [dragNodeId, setDragNodeId] = useState<string | null>(null);
  const [dragWorldId, setDragWorldId] = useState<string | null>(null);
  
  const canvasRef = useRef<HTMLDivElement>(null);
  const selectedNode = nodes.find(n => n.id === selectedNodeId);

  // Keyboard Shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const target = e.target as HTMLElement;
      if (['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)) return;
      if (e.key === 'Delete' || e.key === 'Backspace') {
        if (selectedNodeId) handleDeleteNode(selectedNodeId);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectedNodeId]);

  // Zoom via wheel
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const handleWheel = (e: WheelEvent) => {
      e.preventDefault();
      const rect = canvas.getBoundingClientRect();
      const mouseX = e.clientX - rect.left;
      const mouseY = e.clientY - rect.top;
      setTransform(prev => {
        const sensitivity = e.ctrlKey ? 0.01 : 0.002; 
        const delta = -e.deltaY * sensitivity;
        const zoomFactor = 1 + delta;
        const newScale = Math.min(Math.max(0.2, prev.scale * zoomFactor), 4);
        const newX = mouseX - (mouseX - prev.x) * (newScale / prev.scale);
        const newY = mouseY - (mouseY - prev.y) * (newScale / prev.scale);
        return { x: newX, y: newY, scale: newScale };
      });
    };
    canvas.addEventListener('wheel', handleWheel, { passive: false });
    return () => canvas.removeEventListener('wheel', handleWheel);
  }, []);

  const handleAddNode = (targetWorldId?: string, insertAfterNodeId?: string) => {
    const module_key = targetWorldId || worlds[0].id;
    const world = worlds.find(w => w.id === module_key);
    let coords = { x: (world?.x || 50) + 50, y: (world?.y || 50) + 70 };

    if (insertAfterNodeId) {
      const prevNode = nodes.find(n => n.id === insertAfterNodeId);
      if (prevNode) { coords = { x: prevNode.coordinates.x + 300, y: prevNode.coordinates.y }; }
    } else {
      const worldNodes = nodes.filter(n => n.module_key === module_key);
      const lastNode = worldNodes[worldNodes.length - 1];
      if (lastNode) { coords = { x: lastNode.coordinates.x + 300, y: lastNode.coordinates.y }; }
    }

    createNodeMutation.mutate({
      slug: generateId('node'),
      module_key: module_key,
      coordinates: JSON.stringify(coords),
      title: JSON.stringify({ en: 'New Step' }),
      type: 'form',
      points: 0,
      config: JSON.stringify({ fields: [] })
    });
  };

  const handleDeleteNode = (id: string) => {
    if (confirm("Delete this step?")) {
      deleteNodeMutation.mutate(id);
    }
  };

  const handleAddWorld = () => {
    const newWorldId = generateId('w');
    // Find bottom-most world to place new one below
    const maxY = worlds.length > 0 ? Math.max(...worlds.map(w => w.y)) + 400 : 50;
    
    const newWorld: World = {
      id: newWorldId,
      title: 'New Phase',
      description: '',
      order: worlds.length + 1,
      color: '#64748b',
      collapsed: false,
      condition: null,
      x: 50,
      y: maxY
    };
    setJourneyState(prev => ({
      ...prev,
      worlds: [...prev.worlds, newWorld]
    }));
    setEditingWorld(newWorld);
  };

  const handleUpdateWorld = (updatedWorld: World) => {
    const nextWorlds = worlds.map(w => w.id === updatedWorld.id ? updatedWorld : w);
    setJourneyState(prev => ({
      ...prev,
      worlds: nextWorlds
    }));
    updateMapMutation.mutate({
        config: JSON.stringify(nextWorlds)
    });
    setEditingWorld(null);
  };

  const handleDeleteWorld = (worldId: string) => {
    if (confirm("Delete this phase and all its steps?")) {
      setJourneyState(prev => ({
        ...prev,
        worlds: prev.worlds.filter(w => w.id !== worldId),
        nodes: prev.nodes.filter(n => n.module_key !== worldId),
      }));
      updateMapMutation.mutate({
          config: JSON.stringify(worlds.filter(w => w.id !== worldId))
      });
      setEditingWorld(null);
    }
  };

  const handleReorderNodes = (worldId: string, newOrder: ProgramVersionNode[]) => {
    const worldNodes = nodes.filter(n => n.module_key === worldId);
    const otherNodes = nodes.filter(n => n.module_key !== worldId);
    
    // In a vertical list, the X stays same, but we could auto-space Y if desired.
    // For now we just update the order for the List View.
    setJourneyState(prev => ({
       ...prev,
       nodes: [...otherNodes, ...newOrder]
    }));
  };

  // --- Canvas Logic ---
  const handleMouseDown = (e: React.MouseEvent) => {
    if (e.button === 0) setIsPanning(true);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (isPanning) {
      setTransform(prev => ({
        ...prev,
        x: prev.x + e.movementX,
        y: prev.y + e.movementY,
      }));
    } else if (dragNodeId) {
       setJourneyState(prev => ({
         ...prev,
         nodes: prev.nodes.map(n => n.id === dragNodeId ? {
           ...n,
           coordinates: {
                x: n.coordinates.x + (e.movementX / transform.scale),
                y: n.coordinates.y + (e.movementY / transform.scale)
           }
         } : n)
       }));
    } else if (dragWorldId) {
        setJourneyState(prev => ({
            ...prev,
            worlds: prev.worlds.map(w => w.id === dragWorldId ? {
                ...w,
                x: w.x + (e.movementX / transform.scale),
                y: w.y + (e.movementY / transform.scale)
            } : w),
            // Move nodes with the world
            nodes: prev.nodes.map(n => n.module_key === dragWorldId ? {
                ...n,
                coordinates: {
                    x: n.coordinates.x + (e.movementX / transform.scale),
                    y: n.coordinates.y + (e.movementY / transform.scale)
                }
            } : n)
        }));
    }
  };

  const handleMouseUp = () => {
    if (dragNodeId) {
        const node = nodes.find(n => n.id === dragNodeId);
        if (node) {
            updateNodeMutation.mutate({ 
                id: node.id, 
                data: { coordinates: JSON.stringify(node.coordinates) } 
            });
        }
    }
    if (dragWorldId) {
        const world = worlds.find(w => w.id === dragWorldId);
        if (world) {
            updateMapMutation.mutate({
                config: JSON.stringify(worlds)
            });
        }
    }
    setIsPanning(false);
    setDragNodeId(null);
    setDragWorldId(null);
  };

  const zoomIn = () => setTransform(prev => ({ ...prev, scale: Math.min(prev.scale + 0.1, 2) }));
  const zoomOut = () => setTransform(prev => ({ ...prev, scale: Math.max(prev.scale - 0.1, 0.5) }));
  const resetView = () => setTransform({ x: 0, y: 0, scale: 1 });

  const getEdgePath = (edge: FlowEdge) => {
    const fromNode = nodes.find(n => n.id === edge.from);
    const toNode = nodes.find(n => n.id === edge.to);
    if (!fromNode || !toNode) return '';

    const startX = fromNode.coordinates.x + 260; 
    const startY = fromNode.coordinates.y + 40; 
    const endX = toNode.coordinates.x; 
    const endY = toNode.coordinates.y + 40;
    
    const controlPointX = (startX + endX) / 2;
    return `M ${startX} ${startY} C ${controlPointX} ${startY}, ${controlPointX} ${endY}, ${endX} ${endY}`;
  };

  const NODE_TYPE_GROUPS = [
    {
      label: "Input & Actions",
      types: ['form', 'checklist', 'survey', 'confirmTask']
    },
    {
      label: "Content & Events",
      types: ['info', 'course', 'meeting']
    },
    {
      label: "Operations",
      types: ['approval', 'payment', 'sync_ops', 'milestone']
    }
  ];

  return (
    <div className={cn("flex flex-col h-[calc(100vh-8rem)] -m-8 font-sans transition-all duration-300", isZenMode && "fixed inset-0 z-[100] m-0 bg-slate-50")}>
      {/* 1. Toolbar */}
      <div className="h-16 bg-white border-b border-slate-200 px-6 flex items-center justify-between flex-shrink-0 z-30 relative shadow-sm">
        <div className="flex items-center gap-4">
          <div className="flex bg-slate-100 p-1 rounded-lg">
            <button onClick={() => setViewMode('map')} className={cn("px-3 py-1.5 text-xs font-bold rounded-md flex items-center gap-2 transition-all", viewMode === 'map' ? 'bg-white shadow-sm text-slate-900' : 'text-slate-500 hover:text-slate-700')}>
              <MapIcon size={14} /> Map
            </button>
            <button onClick={() => setViewMode('list')} className={cn("px-3 py-1.5 text-xs font-bold rounded-md flex items-center gap-2 transition-all", viewMode === 'list' ? 'bg-white shadow-sm text-slate-900' : 'text-slate-500 hover:text-slate-700')}>
              <List size={14} /> List
            </button>
          </div>
          <div className="h-6 w-px bg-slate-200" />
          <div className="flex flex-col">
             <h2 className="text-sm font-bold text-slate-900 leading-none">
                {mapData?.title ? parseLocalized(mapData.title) : 'Loading Program...'}
             </h2>
             <span className="text-[10px] text-slate-500">Journey Builder v2.2</span>
          </div>
        </div>
        
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1">
             <IconButton icon={Undo} onClick={undo} disabled={!canUndo} />
             <IconButton icon={Redo} onClick={redo} disabled={!canRedo} />
          </div>
          <div className="h-6 w-px bg-slate-200" />
          <IconButton 
            icon={isZenMode ? Minimize2 : Maximize2} 
            onClick={() => setIsZenMode(!isZenMode)}
            className={cn(isZenMode && "bg-indigo-100 text-indigo-700")}
            title="Toggle Zen Mode"
          />
          <div className="px-3 py-1.5 bg-slate-50 border border-slate-200 rounded-lg flex items-center gap-2">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
            <span className="text-[10px] font-bold text-slate-500 uppercase tracking-tight">Live Sync</span>
          </div>

          <div className="h-6 w-px bg-slate-200" />

          {/* Preview Toggle */}
           <Button 
            variant={isPreview ? "primary" : "secondary"} 
            icon={isPreview ? X : ExternalLink} 
            onClick={() => setIsPreview(!isPreview)}
          >
            {isPreview ? 'Close Preview' : 'Preview Portal'}
          </Button>
        </div>
      </div>

      {/* PREVIEW OVERLAY (Portal to Body) */}
      {isPreview && createPortal(
        <div className="fixed inset-0 z-[9999] bg-white overflow-y-auto animate-in fade-in slide-in-from-bottom-4 duration-300">
           {/* Close Button Floating */}
           <div className="fixed top-4 right-4 z-[10000]">
              <Button onClick={() => setIsPreview(false)} variant="secondary" className="shadow-xl border-slate-200">
                <X size={16} className="mr-2" /> Close Preview
              </Button>
           </div>
           
           <div className="py-12 min-h-screen bg-slate-50/50">
              <JourneyMap 
                playbook={builderStateToPlaybook(worlds, nodes, edges)}
                locale="en"
                onStateChanged={() => {}} 
                viewerRole="student"
                stateByNodeId={nodes.reduce((acc, n) => ({ ...acc, [n.id]: 'active' }), {})}
              />
           </div>
        </div>,
        document.body
      )}

      {/* EDIT MODE CONTENT (Hidden when preview is active to preserve state) */}
      <div className={cn("flex-1 flex overflow-hidden relative", isPreview && "hidden")}>
         
        {/* VIEW: MAP */}
        {viewMode === 'map' && (
          <div className="flex-1 relative overflow-hidden bg-slate-50">
            <div 
              ref={canvasRef}
              className="absolute inset-0 cursor-grab active:cursor-grabbing select-none"
              onMouseDown={handleMouseDown}
              onMouseMove={handleMouseMove}
              onMouseUp={handleMouseUp}
              onMouseLeave={handleMouseUp}
            >
               {/* Grid Pattern */}
               <div 
                  className="absolute inset-0 pointer-events-none opacity-[0.03]"
                  style={{
                    backgroundImage: `linear-gradient(#000 1px, transparent 1px), linear-gradient(90deg, #000 1px, transparent 1px)`,
                    backgroundSize: `${20 * transform.scale}px ${20 * transform.scale}px`,
                    backgroundPosition: `${transform.x}px ${transform.y}px`
                  }}
               />

               <div style={{ transform: `translate(${transform.x}px, ${transform.y}px) scale(${transform.scale})`, transformOrigin: '0 0', width: '100%', height: '100%' }}>
                  {/* World Backgrounds */}
                  {worlds.map((world, idx) => {
                     const worldNodes = nodes.filter(n => n.module_key === world.id);
                     // Calculate height based on nodes or default
                     let height = 200;
                     if (worldNodes.length > 0) {
                        const minY = Math.min(...worldNodes.map(n => n.coordinates.y));
                        const maxY = Math.max(...worldNodes.map(n => n.coordinates.y)) + 120;
                        height = Math.max(200, maxY - world.y); // Keep minimum height
                     }

                     return (
                        <div key={world.id} className="absolute border-2 border-dashed border-slate-200 rounded-3xl -z-10 transition-all group/world" 
                             style={{ 
                               left: world.x, 
                               top: world.y, 
                               width: 2000, 
                               height: height
                             }}>
                           {/* World Header/Handle */}
                           <div 
                              className="absolute -top-4 left-10 px-4 py-1.5 bg-white border border-slate-200 rounded-full font-bold text-xs text-slate-600 shadow-sm flex items-center gap-2 group-hover/world:border-indigo-300 transition-colors cursor-move" 
                              onMouseDown={(e) => { e.stopPropagation(); setDragWorldId(world.id); }}
                           >
                              <div className="w-3 h-3 rounded-full" style={{ backgroundColor: world.color }} />
                              {world.title}
                              {world.condition && <GitMerge size={10} className="text-amber-500" />}
                              <div 
                                className="flex gap-1 ml-1 pl-2 border-l border-slate-100" 
                                onMouseDown={(e) => e.stopPropagation()} // Prevent drag when clicking buttons
                              >
                                <button onClick={() => setEditingWorld(world)} className="text-slate-400 hover:text-indigo-600 transition-colors">
                                  <Settings2 size={12} />
                                </button>
                              </div>
                           </div>
                        </div>
                     );
                  })}

                  {/* Edges */}
                  <svg className="absolute top-0 left-0 w-[5000px] h-[5000px] overflow-visible pointer-events-none z-0">
                    {edges.map(edge => {
                      const path = getEdgePath(edge);
                      if (!path) return null;
                      return <path key={edge.id} d={path} stroke="#94a3b8" strokeWidth="2" strokeDasharray={edge.type === 'dashed' ? '5,5' : '0'} fill="none" />;
                    })}
                  </svg>

                  {/* Nodes */}
                  {nodes.map(node => (
                     <div key={node.id} style={{ position: 'absolute', left: node.coordinates.x, top: node.coordinates.y }}>
                        <FlowNode 
                          {...node} 
                          active={selectedNodeId === node.id} 
                          onClick={setSelectedNodeId} 
                          onMouseDown={(e: any) => setDragNodeId(node.id)}
                          worlds={worlds}
                          validationErrors={validateNode(node)}
                          onAddNext={() => handleAddNode(node.module_key, node.id)} 
                          points={node.points}
                        />
                     </div>
                  ))}
               </div>
            </div>

            {/* Canvas Controls Toolbar */}
            <div className="absolute bottom-6 right-6 flex flex-col gap-2 z-20">
               <div className="bg-white rounded-xl shadow-lg border border-slate-200 p-1 flex flex-col gap-1">
                  <button onClick={zoomIn} className="p-2 hover:bg-slate-100 rounded-lg text-slate-500" title="Zoom In">
                     <ZoomIn size={18} />
                  </button>
                  <button onClick={zoomOut} className="p-2 hover:bg-slate-100 rounded-lg text-slate-500" title="Zoom Out">
                     <ZoomOut size={18} />
                  </button>
                  <button onClick={resetView} className="p-2 hover:bg-slate-100 rounded-lg text-slate-500" title="Reset View">
                     <Target size={18} />
                  </button>
               </div>
               <div className="bg-white rounded-xl shadow-lg border border-slate-200 p-1 flex flex-col gap-1">
                  <button onClick={handleAddWorld} className="p-2 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg text-slate-500" title="Add Phase (Module)">
                     <Layout size={18} />
                  </button>
                  <button onClick={() => handleAddNode()} className="p-2 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg text-slate-500" title="Add Step">
                     <Plus size={18} />
                  </button>
               </div>
            </div>
         </div>
        )}

        {/* VIEW: LIST */}
        {viewMode === 'list' && (
          <ListView 
            worlds={worlds} 
            nodes={nodes} 
            selectedNodeId={selectedNodeId} 
            onSelectNode={setSelectedNodeId}
            onAddNode={handleAddNode}
            onAddWorld={handleAddWorld}
            onReorderNodes={handleReorderNodes}
            onEditWorld={setEditingWorld}
            onDeleteNode={handleDeleteNode}
          />
        )}

        {/* Step Inspector Sidebar */}
        {selectedNode && (
          <StepInspector 
            node={selectedNode}
            phases={worlds as any}
            onUpdate={(updates) => {
              const updatedNode = { ...selectedNode, ...updates };
              setJourneyState(prev => ({
                ...prev,
                nodes: prev.nodes.map(n => n.id === selectedNode.id ? updatedNode : n)
              }));
              updateNodeMutation.mutate({ id: selectedNode.id, data: updates });
            }}
            onDelete={handleDeleteNode}
            onClose={() => setSelectedNodeId(null)}
            onNavigate={handleNavigate}
          />
        )}

        {/* World Settings Modal */}
        {editingWorld && (
          <WorldSettingsModal 
            world={editingWorld}
            allNodes={nodes}
            onSave={handleUpdateWorld}
            onClose={() => setEditingWorld(null)}
            onDelete={() => handleDeleteWorld(editingWorld.id)}
          />
        )}
      </div>
    </div>
  );
};

export default ProgramBuilderPage;
