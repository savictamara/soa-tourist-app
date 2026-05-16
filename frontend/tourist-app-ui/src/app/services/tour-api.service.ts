import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CreateKeyPointRequest, CreateTourRequest, Tour } from '../models/tour.model';
import { KeyPoint } from '../models/key-point.model';

@Injectable({
  providedIn: 'root'
})
export class TourApiService {
  private readonly baseUrl = 'http://localhost:8085/api/tours';

  constructor(private readonly http: HttpClient) {}

  createTour(request: CreateTourRequest): Observable<Tour> {
    console.log('POST create tour', this.baseUrl, request);
    return this.http.post<Tour>(this.baseUrl, request);
  }

  getTour(tourId: string): Observable<Tour> {
    return this.http.get<Tour>(`${this.baseUrl}/${tourId}`);
  }

  getToursByAuthor(authorId: string): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.baseUrl}/author/${authorId}`);
  }

  addKeyPoint(tourId: string, request: CreateKeyPointRequest): Observable<KeyPoint> {
    return this.http.post<KeyPoint>(`${this.baseUrl}/${tourId}/key-points`, request);
  }

  getKeyPoints(tourId: string): Observable<KeyPoint[]> {
    return this.http.get<KeyPoint[]>(`${this.baseUrl}/${tourId}/key-points`);
  }
}
