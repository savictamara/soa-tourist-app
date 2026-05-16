export interface KeyPoint {
  id: string;
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  imageUrl?: string;
  order: number;
  createdAt: string;
  updatedAt: string;
}
