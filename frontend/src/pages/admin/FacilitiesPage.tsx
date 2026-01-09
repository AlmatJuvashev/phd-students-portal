import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Building, MapPin, Plus, MoreHorizontal, Trash2, Edit, 
  Loader2, Users, Layers, Wrench, Check
} from 'lucide-react';
import { toast } from 'sonner';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { DropdownMenu, DropdownItem } from '@/components/ui/dropdown-menu';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { 
  getBuildings, createBuilding, updateBuilding, deleteBuilding,
  getRooms, createRoom, updateRoom, deleteRoom,
  Building as BuildingType, Room, CreateBuildingRequest, CreateRoomRequest
} from '@/features/facilities/api';

// Adding explicit comment to force HMR refresh and clarify types
// BuildingType is the interface, Building is the Lucide icon

const ROOM_TYPES = [
  { value: 'classroom', label: 'Classroom', icon: '🏫' },
  { value: 'lecture_hall', label: 'Lecture Hall', icon: '🎓' },
  { value: 'lab', label: 'Laboratory', icon: '🔬' },
  { value: 'office', label: 'Office', icon: '💼' },
  { value: 'seminar_room', label: 'Seminar Room', icon: '👥' },
  { value: 'simulation_center', label: 'Simulation Center', icon: '🩺' },
  { value: 'conference_room', label: 'Conference Room', icon: '🗣️' },
  { value: 'study_hall', label: 'Study Hall', icon: '📚' },
  { value: 'computer_lab', label: 'Computer Lab', icon: '💻' },
  { value: 'medical_clinic', label: 'Medical Clinic', icon: '🏥' },
];

export const FacilitiesPage: React.FC = () => {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  
  const [activeTab, setActiveTab] = useState('buildings');
  const [selectedBuilding, setSelectedBuilding] = useState<string>('all');
  
  // Building modal state
  const [buildingModal, setBuildingModal] = useState(false);
  const [editingBuilding, setEditingBuilding] = useState<BuildingType | null>(null);
  const [buildingForm, setBuildingForm] = useState<CreateBuildingRequest>({ name: '', address: '', description: '' });
  
  // Room modal state
  const [roomModal, setRoomModal] = useState(false);
  const [editingRoom, setEditingRoom] = useState<Room | null>(null);
  const [roomForm, setRoomForm] = useState<CreateRoomRequest>({ building_id: '', name: '', capacity: 30, floor: 1, type: 'classroom' });

  // Queries
  const { data: buildings = [], isLoading: loadingBuildings } = useQuery({
    queryKey: ['buildings'],
    queryFn: getBuildings,
  });

  const { data: rooms = [], isLoading: loadingRooms } = useQuery({
    queryKey: ['rooms', selectedBuilding],
    queryFn: () => getRooms(selectedBuilding !== 'all' ? selectedBuilding : undefined),
  });

  // Building Mutations
  const createBuildingMutation = useMutation({
    mutationFn: createBuilding,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buildings'] });
      setBuildingModal(false);
      toast.success(t('facilities.building_created', 'Building created'));
    },
    onError: () => toast.error(t('facilities.create_error', 'Failed to create building')),
  });

  const updateBuildingMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateBuildingRequest }) => updateBuilding(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buildings'] });
      setBuildingModal(false);
      setEditingBuilding(null);
      toast.success(t('facilities.building_updated', 'Building updated'));
    },
  });

  const deleteBuildingMutation = useMutation({
    mutationFn: deleteBuilding,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buildings'] });
      toast.success(t('facilities.building_deleted', 'Building deleted'));
    },
  });

  // Room Mutations
  const createRoomMutation = useMutation({
    mutationFn: createRoom,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
      setRoomModal(false);
      toast.success(t('facilities.room_created', 'Room created'));
    },
    onError: () => toast.error(t('facilities.create_error', 'Failed to create room')),
  });

  const updateRoomMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateRoomRequest }) => updateRoom(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
      setRoomModal(false);
      setEditingRoom(null);
      toast.success(t('facilities.room_updated', 'Room updated'));
    },
  });

  const deleteRoomMutation = useMutation({
    mutationFn: deleteRoom,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
      toast.success(t('facilities.room_deleted', 'Room deleted'));
    },
  });

  // Handlers
  const handleOpenBuildingCreate = () => {
    setEditingBuilding(null);
    setBuildingForm({ name: '', address: '', description: '' });
    setBuildingModal(true);
  };

  const handleOpenBuildingEdit = (b: BuildingType) => {
    setEditingBuilding(b);
    setBuildingForm({ name: b.name, address: b.address, description: b.description });
    setBuildingModal(true);
  };

  const handleBuildingSubmit = () => {
    if (editingBuilding) {
      updateBuildingMutation.mutate({ id: editingBuilding.id, data: buildingForm });
    } else {
      createBuildingMutation.mutate(buildingForm);
    }
  };

  const handleOpenRoomCreate = () => {
    setEditingRoom(null);
    setRoomForm({ building_id: buildings[0]?.id || '', name: '', capacity: 30, floor: 1, type: 'classroom' });
    setRoomModal(true);
  };

  const handleOpenRoomEdit = (r: Room) => {
    setEditingRoom(r);
    setRoomForm({ building_id: r.building_id, name: r.name, capacity: r.capacity, floor: r.floor, type: r.type });
    setRoomModal(true);
  };

  const handleRoomSubmit = () => {
    if (editingRoom) {
      updateRoomMutation.mutate({ id: editingRoom.id, data: roomForm });
    } else {
      createRoomMutation.mutate(roomForm);
    }
  };

  // Stats
  const totalCapacity = rooms.reduce((acc, r) => acc + r.capacity, 0);

  return (
    <div className="space-y-6 p-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-black text-slate-900">{t('facilities.title', 'Facilities Management')}</h1>
          <p className="text-slate-500 text-sm">{t('facilities.subtitle', 'Manage campus buildings and rooms')}</p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="bg-gradient-to-br from-indigo-50 to-white border-indigo-100">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-indigo-100 rounded-xl"><Building className="text-indigo-600" size={24} /></div>
              <div>
                <p className="text-2xl font-black text-slate-900">{buildings.length}</p>
                <p className="text-xs text-slate-500 font-bold">{t('facilities.buildings', 'Buildings')}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-emerald-50 to-white border-emerald-100">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-emerald-100 rounded-xl"><MapPin className="text-emerald-600" size={24} /></div>
              <div>
                <p className="text-2xl font-black text-slate-900">{rooms.length}</p>
                <p className="text-xs text-slate-500 font-bold">{t('facilities.rooms', 'Rooms')}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-amber-50 to-white border-amber-100">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-amber-100 rounded-xl"><Users className="text-amber-600" size={24} /></div>
              <div>
                <p className="text-2xl font-black text-slate-900">{totalCapacity}</p>
                <p className="text-xs text-slate-500 font-bold">{t('facilities.total_capacity', 'Total Capacity')}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="buildings" className="gap-2"><Building size={16} /> {t('facilities.buildings', 'Buildings')}</TabsTrigger>
            <TabsTrigger value="rooms" className="gap-2"><MapPin size={16} /> {t('facilities.rooms', 'Rooms')}</TabsTrigger>
          </TabsList>
          {activeTab === 'buildings' ? (
            <Button onClick={handleOpenBuildingCreate} className="gap-2"><Plus size={16} /> {t('facilities.add_building', 'Add Building')}</Button>
          ) : (
            <Button onClick={handleOpenRoomCreate} className="gap-2" disabled={buildings.length === 0}><Plus size={16} /> {t('facilities.add_room', 'Add Room')}</Button>
          )}
        </div>

        {/* Buildings Tab */}
        <TabsContent value="buildings">
          {loadingBuildings ? (
            <div className="flex items-center justify-center py-12"><Loader2 className="animate-spin text-indigo-600" size={32} /></div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {buildings.map((b) => (
                <Card key={b.id} className="group hover:shadow-lg hover:border-indigo-200 transition-all">
                  <CardContent className="pt-6">
                    <div className="flex justify-between items-start mb-4">
                      <div className="w-12 h-12 rounded-xl bg-indigo-50 flex items-center justify-center">
                        <Building className="text-indigo-600" size={24} />
                      </div>
                      <DropdownMenu trigger={<Button variant="ghost" size="sm"><MoreHorizontal size={16} /></Button>}>
                        <DropdownItem onClick={() => handleOpenBuildingEdit(b)}>
                          <Edit size={14} className="mr-2" /> {t('common.edit', 'Edit')}
                        </DropdownItem>
                        <DropdownItem onClick={() => deleteBuildingMutation.mutate(b.id)}>
                          <Trash2 size={14} className="mr-2" /> {t('common.delete', 'Delete')}
                        </DropdownItem>
                      </DropdownMenu>
                    </div>
                    <h3 className="font-bold text-lg text-slate-900 mb-1">{b.name}</h3>
                    <p className="text-sm text-slate-500 mb-3 line-clamp-1">{b.address || t('facilities.no_address', 'No address')}</p>
                    <div className="flex items-center gap-2">
                      <Badge variant={b.is_active ? 'default' : 'secondary'} className={b.is_active ? 'bg-emerald-100 text-emerald-700' : ''}>
                        {b.is_active ? t('common.active', 'Active') : t('common.inactive', 'Inactive')}
                      </Badge>
                    </div>
                  </CardContent>
                </Card>
              ))}
              {buildings.length === 0 && (
                <div className="col-span-full py-12 text-center text-slate-400 border-2 border-dashed border-slate-200 rounded-xl">
                  {t('facilities.no_buildings', 'No buildings yet. Add your first building.')}
                </div>
              )}
            </div>
          )}
        </TabsContent>

        {/* Rooms Tab */}
        <TabsContent value="rooms">
          <div className="mb-4">
            <Select value={selectedBuilding} onValueChange={setSelectedBuilding}>
              <SelectTrigger className="w-60"><SelectValue placeholder="Filter by building" /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{t('common.all', 'All Buildings')}</SelectItem>
                {buildings.map(b => <SelectItem key={b.id} value={b.id}>{b.name}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>

          {loadingRooms ? (
            <div className="flex items-center justify-center py-12"><Loader2 className="animate-spin text-indigo-600" size={32} /></div>
          ) : (
            <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-slate-50 border-b border-slate-200">
                  <tr>
                    <th className="text-left px-4 py-3 font-bold text-slate-600">{t('facilities.room_name', 'Room')}</th>
                    <th className="text-left px-4 py-3 font-bold text-slate-600">{t('facilities.building', 'Building')}</th>
                    <th className="text-left px-4 py-3 font-bold text-slate-600">{t('facilities.type', 'Type')}</th>
                    <th className="text-left px-4 py-3 font-bold text-slate-600">{t('facilities.floor', 'Floor')}</th>
                    <th className="text-left px-4 py-3 font-bold text-slate-600">{t('facilities.capacity', 'Capacity')}</th>
                    <th className="text-right px-4 py-3 font-bold text-slate-600">{t('common.actions', 'Actions')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {rooms.map((r) => {
                    const building = buildings.find(b => b.id === r.building_id);
                    const roomType = ROOM_TYPES.find(t => t.value === r.type);
                    return (
                      <tr key={r.id} className="hover:bg-slate-50 transition-colors">
                        <td className="px-4 py-3 font-bold text-slate-900">{r.name}</td>
                        <td className="px-4 py-3 text-slate-600">{building?.name || '-'}</td>
                        <td className="px-4 py-3">
                          <Badge variant="outline" className="gap-1">
                            <span>{roomType?.icon}</span> {roomType?.label || r.type}
                          </Badge>
                        </td>
                        <td className="px-4 py-3 text-slate-600">{r.floor}</td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-1 text-slate-600">
                            <Users size={14} /> {r.capacity}
                          </div>
                        </td>
                        <td className="px-4 py-3 text-right">
                          <DropdownMenu trigger={<Button variant="ghost" size="sm"><MoreHorizontal size={16} /></Button>}>
                            <DropdownItem onClick={() => handleOpenRoomEdit(r)}>
                              <Edit size={14} className="mr-2" /> {t('common.edit', 'Edit')}
                            </DropdownItem>
                            <DropdownItem onClick={() => deleteRoomMutation.mutate(r.id)}>
                              <Trash2 size={14} className="mr-2" /> {t('common.delete', 'Delete')}
                            </DropdownItem>
                          </DropdownMenu>
                        </td>
                      </tr>
                    );
                  })}
                  {rooms.length === 0 && (
                    <tr>
                      <td colSpan={6} className="px-4 py-12 text-center text-slate-400">
                        {t('facilities.no_rooms', 'No rooms found')}
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </TabsContent>
      </Tabs>

      {/* Building Modal */}
      <Dialog open={buildingModal} onOpenChange={setBuildingModal}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{editingBuilding ? t('facilities.edit_building', 'Edit Building') : t('facilities.add_building', 'Add Building')}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div>
              <Label>{t('facilities.name', 'Name')}</Label>
              <Input value={buildingForm.name} onChange={(e) => setBuildingForm({ ...buildingForm, name: e.target.value })} placeholder="Main Campus" />
            </div>
            <div>
              <Label>{t('facilities.address', 'Address')}</Label>
              <Input value={buildingForm.address} onChange={(e) => setBuildingForm({ ...buildingForm, address: e.target.value })} placeholder="123 University Ave" />
            </div>
            <div>
              <Label>{t('facilities.description', 'Description')}</Label>
              <Textarea value={buildingForm.description} onChange={(e) => setBuildingForm({ ...buildingForm, description: e.target.value })} placeholder="Description..." />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setBuildingModal(false)}>{t('common.cancel', 'Cancel')}</Button>
            <Button onClick={handleBuildingSubmit} disabled={createBuildingMutation.isPending || updateBuildingMutation.isPending}>
              {(createBuildingMutation.isPending || updateBuildingMutation.isPending) && <Loader2 className="animate-spin mr-2" size={16} />}
              {editingBuilding ? t('common.save', 'Save') : t('common.create', 'Create')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Room Modal */}
      <Dialog open={roomModal} onOpenChange={setRoomModal}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{editingRoom ? t('facilities.edit_room', 'Edit Room') : t('facilities.add_room', 'Add Room')}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div>
              <Label>{t('facilities.building', 'Building')}</Label>
              <Select value={roomForm.building_id} onValueChange={(v) => setRoomForm({ ...roomForm, building_id: v })}>
                <SelectTrigger><SelectValue placeholder="Select building" /></SelectTrigger>
                <SelectContent>
                  {buildings.map(b => <SelectItem key={b.id} value={b.id}>{b.name}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label>{t('facilities.room_name', 'Room Name')}</Label>
                <Input value={roomForm.name} onChange={(e) => setRoomForm({ ...roomForm, name: e.target.value })} placeholder="Room 101" />
              </div>
              <div>
                <Label>{t('facilities.floor', 'Floor')}</Label>
                <Input type="number" value={roomForm.floor} onChange={(e) => setRoomForm({ ...roomForm, floor: parseInt(e.target.value) || 1 })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label>{t('facilities.type', 'Type')}</Label>
                <Select value={roomForm.type} onValueChange={(v) => setRoomForm({ ...roomForm, type: v })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {ROOM_TYPES.map(t => <SelectItem key={t.value} value={t.value}>{t.icon} {t.label}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label>{t('facilities.capacity', 'Capacity')}</Label>
                <Input type="number" value={roomForm.capacity} onChange={(e) => setRoomForm({ ...roomForm, capacity: parseInt(e.target.value) || 1 })} />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setRoomModal(false)}>{t('common.cancel', 'Cancel')}</Button>
            <Button onClick={handleRoomSubmit} disabled={createRoomMutation.isPending || updateRoomMutation.isPending}>
              {(createRoomMutation.isPending || updateRoomMutation.isPending) && <Loader2 className="animate-spin mr-2" size={16} />}
              {editingRoom ? t('common.save', 'Save') : t('common.create', 'Create')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default FacilitiesPage;
