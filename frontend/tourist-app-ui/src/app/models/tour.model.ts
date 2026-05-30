import { KeyPoint } from './key-point.model';

export interface Tour {
  id: string;
  authorId: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string[];
  status: string;
  publishedAt?: string;
  archivedAt?: string;
  reactivatedAt?: string;
  lengthKm: number;
  durations: TourDuration[];
  price: number;
  keyPoints: KeyPoint[];
  reviews: Review[];
  createdAt: string;
  updatedAt: string;
}

export interface TourDuration {
  transportType: 'walking' | 'bicycle' | 'car';
  minutes: number;
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

export interface CompletedKeyPoint {
  keyPointId: string;
  keyPointName: string;
  reachedAt: string;
}

export interface TourExecution {
  id: string;
  tourId: string;
  tourName: string;
  touristId: string;
  status: 'active' | 'completed' | 'abandoned';
  startedAt: string;
  completedAt?: string;
  abandonedAt?: string;
  lastActivityAt: string;
  startLatitude: number;
  startLongitude: number;
  currentLatitude: number;
  currentLongitude: number;
  completedKeyPoints: CompletedKeyPoint[];
}

export interface StartTourExecutionRequest {
  touristId: string;
  latitude: number;
  longitude: number;
}

export interface CheckTourExecutionLocationRequest {
  latitude: number;
  longitude: number;
}

export interface CheckTourExecutionLocationResponse {
  execution: TourExecution;
  keyPointReached: boolean;
  reachedKeyPoint?: CompletedKeyPoint;
  distanceMeters: number;
  lastActivityAt: string;
  completedCount: number;
  totalKeyPointCount: number;
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

export interface UpdateKeyPointRequest {
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  imageUrl: string;
}

export interface UpdateDurationsRequest {
  durations: TourDuration[];
}

export interface UpdatePriceRequest {
  price: number;
}

export interface CreateReviewRequest {
  rating: number;
  comment: string;
  touristId: string;
  touristUsername: string;
  visitedDate: string;
  images: string[];
}
