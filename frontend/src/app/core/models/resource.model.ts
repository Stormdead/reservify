export interface Resource {
  id: number;
  name: string;
  description: string;
  capacity: number;
  price_per_hour: number;
  category: string;
  image_url: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateResourceRequest {
  name: string;
  description: string;
  capacity: number;
  price_per_hour: number;
  category: string;
  image_url: string;
}