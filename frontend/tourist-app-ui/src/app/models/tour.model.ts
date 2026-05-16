import { KeyPoint } from './key-point.model';

export interface Tour {
  id: string;
  authorId: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string[];
  status: string;
  price: number;
  keyPoints: KeyPoint[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateTourRequest {
  authorId: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string[];
}

export interface CreateKeyPointRequest {
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  imageUrl?: string;
}
