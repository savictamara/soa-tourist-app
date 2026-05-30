import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CheckTourExecutionLocationRequest, CheckTourExecutionLocationResponse, CreateKeyPointRequest, CreateReviewRequest, Review, CreateTourRequest, StartTourExecutionRequest, Tour, TourExecution, UpdateDurationsRequest, UpdateKeyPointRequest, UpdatePriceRequest } from '../models/tour.model';
import { KeyPoint } from '../models/key-point.model';


@Injectable({
  providedIn: 'root'
})
export class TourApiService {
  private readonly baseUrl = 'http://localhost:8080/api/tours';

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

  getTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(this.baseUrl);
  }

  getPublishedTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.baseUrl}/published`);
  }

  getAvailableTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.baseUrl}/available`);
  }

  getActiveExecution(touristId: string): Observable<TourExecution> {
    return this.http.get<TourExecution>(`${this.baseUrl}/executions/active/${touristId}`);
  }

  getLatestExecution(tourId: string, touristId: string): Observable<TourExecution> {
    return this.http.get<TourExecution>(`${this.baseUrl}/${tourId}/executions/latest/${touristId}`);
  }

  startTourExecution(tourId: string, request: StartTourExecutionRequest): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.baseUrl}/${tourId}/executions/start`, request);
  }

  checkExecutionLocation(executionId: string, request: CheckTourExecutionLocationRequest): Observable<CheckTourExecutionLocationResponse> {
    return this.http.post<CheckTourExecutionLocationResponse>(`${this.baseUrl}/executions/${executionId}/check-location`, request);
  }

  completeExecution(executionId: string): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.baseUrl}/executions/${executionId}/complete`, {});
  }

  abandonExecution(executionId: string): Observable<TourExecution> {
    return this.http.post<TourExecution>(`${this.baseUrl}/executions/${executionId}/abandon`, {});
  }

  getAllTours(): Observable<Tour[]> {
    return this.getTours();
  }

  getTourById(tourId: string): Observable<Tour> {
    return this.getTour(tourId);
  }

  addKeyPoint(tourId: string, request: CreateKeyPointRequest): Observable<KeyPoint> {
    return this.http.post<KeyPoint>(`${this.baseUrl}/${tourId}/key-points`, request);
  }

  getKeyPoints(tourId: string): Observable<KeyPoint[]> {
    return this.http.get<KeyPoint[]>(`${this.baseUrl}/${tourId}/key-points`);
  }

  updateKeyPoint(tourId: string, keyPointId: string, request: UpdateKeyPointRequest): Observable<KeyPoint> {
    return this.http.put<KeyPoint>(`${this.baseUrl}/${tourId}/key-points/${keyPointId}`, request);
  }

  deleteKeyPoint(tourId: string, keyPointId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${tourId}/key-points/${keyPointId}`);
  }

  updateDurations(tourId: string, request: UpdateDurationsRequest): Observable<Tour> {
    return this.http.put<Tour>(`${this.baseUrl}/${tourId}/durations`, request);
  }

  updatePrice(tourId: string, request: UpdatePriceRequest): Observable<Tour> {
    return this.http.put<Tour>(`${this.baseUrl}/${tourId}/price`, request);
  }

  publishTour(tourId: string): Observable<Tour> {
    return this.http.post<Tour>(`${this.baseUrl}/${tourId}/publish`, {});
  }

  archiveTour(tourId: string): Observable<Tour> {
    return this.http.post<Tour>(`${this.baseUrl}/${tourId}/archive`, {});
  }

  reactivateTour(tourId: string): Observable<Tour> {
    return this.http.post<Tour>(`${this.baseUrl}/${tourId}/reactivate`, {});
  }

  addReview(tourId: string, request: CreateReviewRequest): Observable<Review> {
    return this.http.post<Review>(`${this.baseUrl}/${tourId}/reviews`, request);
  }

  getReviews(tourId: string): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.baseUrl}/${tourId}/reviews`);
  }
}
