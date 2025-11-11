import { User } from './user.model';
import { Resource } from './resource.model';

export type BookingStatus = 'pending' | 'confirmed' | 'cancelled' | 'completed';

export interface Booking {
  id: number;
  user_id: number;
  user: User;
  resource_id: number;
  resource: Resource;
  start_datetime: string;
  end_datetime: string;
  status: BookingStatus;
  total_price: number;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface CreateBookingRequest {
  resource_id: number;
  start_datetime: string;
  end_datetime: string;
  notes?: string;
}

export interface BookingListItem {
  id: number;
  user_id: number;
  user_name: string;
  resource_id: number;
  resource_name: string;
  start_datetime: string;
  end_datetime: string;
  status: BookingStatus;
  total_price: number;
  created_at: string;
}