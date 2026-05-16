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
  reviews: Review[];
  createdAt: string;
  updatedAt: string;
}

export interface Review {
  id: string;
  rating: number;
  comment: string;
  touristId: string;
  touristUsername: string;
  visitedDate: string;
  commentDate: string;
  images: string[];
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

export interface CreateReviewRequest {
  rating: number;
  comment: string;
  touristId: string;
  touristUsername: string;
  visitedDate: string;
  images: string[];
}
