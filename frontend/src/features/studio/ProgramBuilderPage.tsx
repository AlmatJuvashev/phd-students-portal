import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ProgramJourneyBuilder } from './ProgramJourneyBuilder';
import { createProgramVersionNode, getProgram, getProgramVersionMap, updateProgramVersionMap, updateProgramVersionNode } from '@/features/curriculum/api';
import { Loader2, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ProgramPhase, ProgramVersion, ProgramVersionNode } from './types';

const pickLocaleText = (value: any): string => {
    if (typeof value === 'string') return value;
    if (value && typeof value === 'object') {
        return value.en || value.ru || value.kk || Object.values(value)[0] || '';
    }
    return '';
};

const normalizePhases = (phases: any[]): ProgramPhase[] => {
    if (!Array.isArray(phases) || phases.length === 0) {
        return [
            { id: 'I', title: 'Phase I', order: 1, color: '#6366f1', position: { x: 50, y: 50 } },
            { id: 'II', title: 'Phase II', order: 2, color: '#ec4899', position: { x: 450, y: 50 } },
        ];
    }
    return phases.map((p: any, idx: number) => ({
        id: String(p.id || p.module_key || idx + 1),
        title: pickLocaleText(p.title) || String(p.title || `Phase ${idx + 1}`),
        description: pickLocaleText(p.description) || undefined,
        color: String(p.color || '#6366f1'),
        order: Number(p.order ?? idx + 1),
        position: {
            x: Number(p.position?.x ?? 50 + idx * 250),
            y: Number(p.position?.y ?? 50),
        },
    }));
};

const normalizeNodes = (nodes: any[]): ProgramVersionNode[] => {
    if (!Array.isArray(nodes)) return [];
    return nodes.map((n: any) => ({
        id: String(n.id),
        program_version_id: n.program_version_id || n.journey_map_id,
        parent_node_id: n.parent_node_id || undefined,
        slug: String(n.slug || ''),
        type: (n.type || 'form') as any,
        title: pickLocaleText(n.title) || String(n.title || 'Untitled'),
        description: pickLocaleText(n.description) || undefined,
        module_key: String(n.module_key || 'I'),
        coordinates: {
            x: Number(n.coordinates?.x ?? 0),
            y: Number(n.coordinates?.y ?? 0),
        },
        config: n.config || {},
        prerequisites: Array.isArray(n.prerequisites) ? n.prerequisites : [],
        points: typeof n.points === 'number' ? n.points : undefined,
    }));
};

export const ProgramBuilderPage = () => {
    const { programId } = useParams();
    const navigate = useNavigate();

    const { data: program, isLoading: isLoadingProgram } = useQuery({
        queryKey: ['program', programId],
        queryFn: () => getProgram(programId as string),
        enabled: !!programId,
    });

    const mapQuery = useQuery({
        queryKey: ['program-map', programId],
        queryFn: () => getProgramVersionMap(programId as string),
        enabled: !!programId,
        select: (data: any): ProgramVersion => ({
            id: String(data.id),
            program_id: String(data.program_id),
            title: pickLocaleText(data.title) || 'Program Version',
            version: String(data.version || '0.0.0'),
            phases: normalizePhases(data.phases || []),
            nodes: normalizeNodes(data.nodes || []),
            edges: [],
        }),
    });
    const map = mapQuery.data;
    const isLoadingMap = mapQuery.isLoading;

    const hasNodeChanges = (a: ProgramVersionNode, b: ProgramVersionNode): boolean => {
        const stable = (v: any) => JSON.stringify(v ?? null);
        return (
            a.slug !== b.slug ||
            a.type !== b.type ||
            a.title !== b.title ||
            (a.description || '') !== (b.description || '') ||
            a.module_key !== b.module_key ||
            stable(a.coordinates) !== stable(b.coordinates) ||
            stable(a.config) !== stable(b.config) ||
            stable(a.prerequisites || []) !== stable(b.prerequisites || []) ||
            (a.parent_node_id || '') !== (b.parent_node_id || '')
        );
    };

    const handleSave = async (updatedMap: ProgramVersion) => {
        if (!programId || !map) return;

        const originalById = new Map(map.nodes.map((n) => [n.id, n]));
        const creates = updatedMap.nodes.filter((n) => !originalById.has(n.id));
        const updates = updatedMap.nodes.filter((n) => {
            const prev = originalById.get(n.id);
            return prev && hasNodeChanges(n, prev);
        });
        const phasesChanged =
            JSON.stringify(updatedMap.phases || []) !== JSON.stringify(map.phases || []);

        if (creates.length === 0 && updates.length === 0 && !phasesChanged) return;

        const ops: Promise<any>[] = [];

        if (phasesChanged) {
            ops.push(updateProgramVersionMap(programId, { phases: updatedMap.phases || [] }));
        }

        for (const n of creates) {
            ops.push(
                createProgramVersionNode(programId, {
                    parent_node_id: n.parent_node_id ?? null,
                    slug: n.slug,
                    type: n.type,
                    title: n.title,
                    description: n.description ?? null,
                    module_key: n.module_key,
                    coordinates: n.coordinates,
                    config: n.config ?? {},
                    prerequisites: n.prerequisites ?? [],
                })
            );
        }

        for (const n of updates) {
            ops.push(
                updateProgramVersionNode(programId, n.id, {
                    parent_node_id: n.parent_node_id ?? null,
                    slug: n.slug,
                    type: n.type,
                    title: n.title,
                    description: n.description ?? null,
                    module_key: n.module_key,
                    coordinates: n.coordinates,
                    config: n.config ?? {},
                    prerequisites: n.prerequisites ?? [],
                })
            );
        }

        await Promise.all(ops);

        await mapQuery.refetch();
    };

    if (isLoadingProgram || isLoadingMap) {
        return <div className="flex items-center justify-center h-screen"><Loader2 className="animate-spin" /></div>;
    }

    if (!program) return <div>Program not found</div>;

    return (
        <div className="flex flex-col h-screen bg-slate-50">
             <div className="h-14 bg-white border-b px-6 flex items-center gap-4">
                 <Button variant="ghost" size="icon" onClick={() => navigate('/admin/programs')}>
                     <ArrowLeft size={16} />
                 </Button>
                 <div>
                     <h1 className="font-bold text-slate-900">{program.title}</h1>
                     <div className="text-xs text-slate-500">Program Builder</div>
                 </div>
             </div>
             
             <div className="flex-1 overflow-hidden p-8">
                 {map ? (
                   <ProgramJourneyBuilder 
                      initialMap={map}
                      onSave={handleSave}
                      onNavigate={(path) => navigate(path)}
                   />
                 ) : (
                   <div className="p-8 text-slate-500">No program version found.</div>
                 )}
             </div>
        </div>
    );
};
