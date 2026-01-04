import React, { useEffect, useState, useRef } from 'react';
import { 
  Map as MapIcon, Plus, BookOpen, CreditCard, Sparkles, Stamp, FileText, Calendar, 
  CheckSquare, Flag, Info, ClipboardList, CheckCircle2, Layout, Award, Settings2, GripVertical, ChevronRight, X, List
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { ProgramPhase, ProgramVersion, ProgramVersionNode, FlowEdge } from './types';
import { useTranslation } from 'react-i18next';
import { CoursePickerDialog } from './components/CoursePickerDialog';

interface ProgramJourneyBuilderProps {
  initialMap?: ProgramVersion;
  onSave: (map: ProgramVersion) => void;
  onNavigate: (path: string) => void;
}

// Map Node Types to Activity Types for UI Logic
const NODE_VISUALS: Record<string, { icon: any, color: string, bg: string, border: string, label: string }> = {
  'course': { icon: BookOpen, color: 'text-purple-600', bg: 'bg-purple-50', border: 'border-purple-200', label: 'Learning Event' },
  'payment': { icon: CreditCard, color: 'text-amber-600', bg: 'bg-amber-50', border: 'border-amber-200', label: 'Financial Gate' },
  'sync_ops': { icon: Sparkles, color: 'text-emerald-600', bg: 'bg-emerald-50', border: 'border-emerald-200', label: 'Ops Automation' },
  'approval': { icon: Stamp, color: 'text-slate-600', bg: 'bg-slate-50', border: 'border-slate-300', label: 'Admin Gate' },
  'form': { icon: FileText, color: 'text-blue-600', bg: 'bg-blue-50', border: 'border-blue-200', label: 'Data Collection' },
  'meeting': { icon: Calendar, color: 'text-pink-600', bg: 'bg-pink-50', border: 'border-pink-200', label: 'Sync Event' },
  'checklist': { icon: CheckSquare, color: 'text-orange-600', bg: 'bg-orange-50', border: 'border-orange-200', label: 'Requirement' },
  'milestone': { icon: Flag, color: 'text-indigo-600', bg: 'bg-indigo-50', border: 'border-indigo-200', label: 'Milestone' },
  'info': { icon: Info, color: 'text-cyan-600', bg: 'bg-cyan-50', border: 'border-cyan-200', label: 'Information' },
  'survey': { icon: ClipboardList, color: 'text-teal-600', bg: 'bg-teal-50', border: 'border-teal-200', label: 'Feedback' },
  'confirmTask': { icon: CheckCircle2, color: 'text-green-600', bg: 'bg-green-50', border: 'border-green-200', label: 'Confirmation' },
  'default': { icon: Layout, color: 'text-slate-500', bg: 'bg-white', border: 'border-slate-200', label: 'Process Step' },
};

const generateId = (prefix: string) => `${prefix}_${Math.random().toString(36).substr(2, 5)}`;

const DEFAULT_PHASES: ProgramPhase[] = [
  { id: 'I', title: 'Phase I', order: 1, color: '#6366f1', position: { x: 50, y: 50 } },
];

export const ProgramJourneyBuilder: React.FC<ProgramJourneyBuilderProps> = ({ initialMap, onSave }) => {
  const { t } = useTranslation('common');
  
  // State
  const [phases, setPhases] = useState<ProgramPhase[]>(initialMap?.phases || DEFAULT_PHASES);
  const [nodes, setNodes] = useState<ProgramVersionNode[]>(initialMap?.nodes || []);
  const [edges, setEdges] = useState<FlowEdge[]>(initialMap?.edges || []);
  
  const [viewMode, setViewMode] = useState<'list' | 'map'>('map');
  const [transform, setTransform] = useState({ x: 0, y: 0, scale: 1 });
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [showCoursePicker, setShowCoursePicker] = useState(false);
  
  // Dragging
  const [dragNodeId, setDragNodeId] = useState<string | null>(null);
  const [dragPhaseId, setDragPhaseId] = useState<string | null>(null);
  const [isPanning, setIsPanning] = useState(false);
  const canvasRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!initialMap) return;
    setPhases(initialMap.phases?.length ? initialMap.phases : DEFAULT_PHASES);
    setNodes(initialMap.nodes || []);
    setEdges(initialMap.edges || []);
  }, [initialMap]);

  const selectedNode = nodes.find(n => n.id === selectedNodeId);

  const updateNode = (id: string, updates: Partial<ProgramVersionNode>) => {
    setNodes(prev => prev.map(n => n.id === id ? { ...n, ...updates } : n));
  };

  const handleCreateNode = (phaseId: string) => {
    const parentPhase = phases.find(p => p.id === phaseId);
    if (!parentPhase) return;

    const newNode: ProgramVersionNode = {
        id: generateId('node'),
        slug: generateId('slug'),
        type: 'form',
        title: 'New Step',
        module_key: phaseId,
        coordinates: { 
            x: parentPhase.position.x + 50, 
            y: parentPhase.position.y + 100 
        },
        config: {},
        points: 0
    };

    setNodes(prev => [...prev, newNode]);
    setSelectedNodeId(newNode.id);
  };
  
  const getCurrentVisuals = (type: string) => NODE_VISUALS[type] || NODE_VISUALS['default'];

  // --- Canvas Logic ---
  const handleMouseDown = (e: React.MouseEvent) => {
    if (e.button === 0) setIsPanning(true);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (dragNodeId) {
       updateNode(dragNodeId, { 
           coordinates: { 
               x: (selectedNode?.coordinates?.x || 0) + e.movementX/transform.scale, 
               y: (selectedNode?.coordinates?.y || 0) + e.movementY/transform.scale 
           } 
       });
    } else if (dragPhaseId) {
       // Drag Phase
       setPhases(prev => prev.map(p => p.id === dragPhaseId ? {
           ...p,
           position: {
               x: p.position.x + e.movementX/transform.scale,
               y: p.position.y + e.movementY/transform.scale
           }
       } : p));
    } else if (isPanning) {
      setTransform(prev => ({ ...prev, x: prev.x + e.movementX, y: prev.y + e.movementY }));
    }
  };

  const handleMouseUp = () => {
    setIsPanning(false);
    setDragNodeId(null);
    setDragPhaseId(null);
  };

  return (
    <div className="flex flex-col h-[calc(100vh-8rem)] -m-8 font-sans transition-all duration-300">
       {/* Toolbar */}
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
          </div>
          <div className="flex items-center gap-4">
             <Button
               onClick={() =>
                 onSave({
                   id: initialMap?.id || '',
                   program_id: initialMap?.program_id || '',
                   title: initialMap?.title || 'Program Version',
                   version: initialMap?.version || '0.0.0',
                   nodes,
                   edges,
                   phases,
                 })
               }
               size="sm"
             >
               Save Journey
             </Button>
          </div>
       </div>

       <div className="flex-1 flex overflow-hidden relative">
          {viewMode === 'map' && (
             <div className="flex-1 relative overflow-hidden bg-slate-50 bg-[radial-gradient(#e5e7eb_1px,transparent_1px)] [background-size:16px_16px]">
                <div 
                  ref={canvasRef}
                  className="absolute inset-0 cursor-grab active:cursor-grabbing select-none"
                  onMouseDown={handleMouseDown}
                  onMouseMove={handleMouseMove}
                  onMouseUp={handleMouseUp}
                  onMouseLeave={handleMouseUp}
                >
                   <div style={{ transform: `translate(${transform.x}px, ${transform.y}px) scale(${transform.scale})`, transformOrigin: '0 0', width: '100%', height: '100%' }}>
                      
                      {/* Phases (Modules) */}
                      {phases.map(phase => (
                         <div 
                            key={phase.id} 
                            style={{ left: phase.position.x, top: phase.position.y, width: 2000, height: 400 }}
                            className="absolute border-2 border-dashed border-slate-200 rounded-3xl -z-10 group/world"
                         >
                            <div 
                              className="absolute -top-4 left-10 px-4 py-1.5 bg-white border border-slate-200 rounded-full font-bold text-xs text-slate-600 shadow-sm flex items-center gap-2 cursor-move"
                              onMouseDown={(e) => { e.stopPropagation(); setDragPhaseId(phase.id); }}
                            >
                               <div className="w-3 h-3 rounded-full" style={{ backgroundColor: phase.color }} />
                               {phase.title}
                               <button onClick={(e) => { e.stopPropagation(); handleCreateNode(phase.id); }} className="ml-2 p-1 hover:bg-slate-100 rounded-full text-indigo-600"><Plus size={12} /></button>
                            </div>
                         </div>
                      ))}

                      {/* Nodes */}
                      {nodes.map(node => {
                          const visuals = getCurrentVisuals(node.type);
                          const Icon = visuals.icon;
                          const isActive = selectedNodeId === node.id;
                          
                          return (
                              <div 
                                key={node.id} 
                                style={{ left: node.coordinates.x, top: node.coordinates.y }}
                                className={cn(
                                    "absolute w-[260px] bg-white rounded-2xl shadow-sm flex flex-col cursor-grab active:cursor-grabbing transition-all hover:shadow-xl hover:-translate-y-1 z-10 group select-none",
                                    isActive ? "ring-2 ring-indigo-500 shadow-md" : "border border-slate-200"
                                )}
                                onMouseDown={(e) => { e.stopPropagation(); setDragNodeId(node.id); }}
                                onClick={(e) => { e.stopPropagation(); setSelectedNodeId(node.id); }}
                              >
                                  <div className="p-1">
                                     <div className={cn("rounded-xl p-3 flex items-start gap-3", visuals.bg)}>
                                         <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center bg-white shadow-sm border border-white/50", visuals.color)}>
                                            <Icon size={20} />
                                         </div>
                                         <div className="flex-1 min-w-0 pt-0.5">
                                            <div className="text-[10px] font-black uppercase tracking-wider opacity-60 mb-0.5">{visuals.label}</div>
                                            <div className="font-bold text-slate-900 text-sm leading-tight line-clamp-2">{node.title || "Untitled"}</div>
                                         </div>
                                     </div>
                                  </div>
                              </div>
                          );
                      })}
                   </div>
                </div>
             </div>
          )}
          
          {selectedNode && (
              <div className="w-80 bg-white border-l border-slate-200 p-6 flex flex-col z-20 shadow-xl">
                  <div className="flex justify-between items-center mb-6">
                      <h3 className="font-bold text-slate-800">Properties</h3>
                      <button onClick={() => setSelectedNodeId(null)}><X size={16} /></button>
                  </div>
                  <div className="space-y-4">
                      <div>
                          <label className="text-xs font-bold text-slate-500">Title</label>
                          <Input 
                            value={selectedNode.title} 
                            onChange={(e) => updateNode(selectedNode.id, { title: e.target.value })} 
                          />
                      </div>
                      <div>
                         <label className="text-xs font-bold text-slate-500">Type</label>
                         <select 
                           value={selectedNode.type} 
                           onChange={(e) => updateNode(selectedNode.id, { type: e.target.value as any })}
                           className="w-full text-sm border rounded p-2"
                         >
                            {Object.keys(NODE_VISUALS).filter(t => t !== 'default').map(t => (
                              <option key={t} value={t}>{NODE_VISUALS[t].label}</option>
                            ))}
                         </select>
                      </div>
                      
                      {/* Dynamic Config Based on Type */}
                      {selectedNode.type === 'course' && (
                         <div className="p-4 bg-purple-50 rounded-lg border border-purple-100">
                            <label className="text-xs font-bold text-purple-700 block mb-2">Linked Course</label>
                            {selectedNode.config?.course_id ? (
                                <div className="space-y-2">
                                    <div className="text-sm font-bold text-slate-900 p-2 bg-white rounded border border-purple-200">
                                        {selectedNode.config.course_title || 'Linked Course'}
                                    </div>
                                    <Button 
                                        variant="outline" 
                                        size="sm" 
                                        className="w-full bg-white text-xs"
                                        onClick={() => setShowCoursePicker(true)}
                                    >
                                        Change Course
                                    </Button>
                                </div>
                            ) : (
                                <Button 
                                    variant="outline" 
                                    size="sm" 
                                    className="w-full bg-white"
                                    onClick={() => setShowCoursePicker(true)}
                                >
                                    Select Course
                                </Button>
                            )}
                         </div>
                      )}
                  </div>
              </div>
          )}
       </div>

       {showCoursePicker && selectedNode && (
          <CoursePickerDialog
            isOpen={showCoursePicker}
            onClose={() => setShowCoursePicker(false)}
            onSelect={(course) => {
                updateNode(selectedNode.id, { 
                    config: { 
                        ...selectedNode.config, 
                        course_id: course.id,
                        course_title: course.title,
                        course_code: course.code
                    } 
                });
                setShowCoursePicker(false);
            }}
            programId={initialMap?.program_id}
          />
       )}
    </div>
  );
};
