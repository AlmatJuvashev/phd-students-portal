import { api } from '@/api/client';

// Types
export interface Building {
  id: string;
  tenant_id: string;
  name: string;
  address: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Room {
  id: string;
  building_id: string;
  name: string;
  capacity: number;
  floor: number;
  department_id?: string;
  type: 'lecture_hall' | 'lab' | 'office' | 'classroom' | 'seminar_room' | 'simulation_center' | 'conference_room' | 'study_hall' | 'computer_lab' | 'medical_clinic';
  features: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateBuildingRequest {
  name: string;
  address?: string;
  description?: string;
}

export interface UpdateBuildingRequest extends CreateBuildingRequest {
  id: string;
}

export interface CreateRoomRequest {
  building_id: string;
  name: string;
  capacity: number;
  floor?: number;
  type: string;
  features?: string;
}

export interface UpdateRoomRequest extends CreateRoomRequest {
  id: string;
}

// Buildings API
export const getBuildings = async (): Promise<Building[]> => {
  const res = await api.get('/resources/buildings');
  return res.data;
};

export const createBuilding = async (data: CreateBuildingRequest): Promise<Building> => {
  const res = await api.post('/resources/buildings', data);
  return res.data;
};

export const updateBuilding = async (id: string, data: Omit<UpdateBuildingRequest, 'id'>): Promise<Building> => {
  const res = await api.put(`/resources/buildings/${id}`, data);
  return res.data;
};

export const deleteBuilding = async (id: string): Promise<void> => {
  await api.delete(`/resources/buildings/${id}`);
};

// Rooms API
export const getRooms = async (buildingId?: string): Promise<Room[]> => {
  const params = buildingId ? `?building_id=${buildingId}` : '';
  const res = await api.get(`/resources/rooms${params}`);
  return res.data;
};

export const createRoom = async (data: CreateRoomRequest): Promise<Room> => {
  const res = await api.post('/resources/rooms', data);
  return res.data;
};

export const updateRoom = async (id: string, data: Omit<UpdateRoomRequest, 'id'>): Promise<Room> => {
  const res = await api.put(`/resources/rooms/${id}`, data);
  return res.data;
};

export const deleteRoom = async (id: string): Promise<void> => {
  await api.delete(`/resources/rooms/${id}`);
};
