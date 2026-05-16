import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { FollowRequest, FollowerUser } from '../models/follower.model';

@Injectable({
  providedIn: 'root'
})
export class FollowerApiService {
  private readonly baseUrl = 'http://localhost:8080/api/followers';

  constructor(private readonly http: HttpClient) {}

  createUser(user: FollowerUser): Observable<FollowerUser> {
    return this.http.post<FollowerUser>(`${this.baseUrl}/users`, user);
  }

  syncUser(user: FollowerUser): Observable<FollowerUser> {
    return this.http.post<FollowerUser>(`${this.baseUrl}/sync-user`, user);
  }

  getUsers(): Observable<FollowerUser[]> {
    return this.http.get<FollowerUser[]>(`${this.baseUrl}/users`);
  }

  follow(request: FollowRequest): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.baseUrl}/follow`, request);
  }

  unfollow(request: FollowRequest): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.baseUrl}/follow`, { body: request });
  }

  getFollowing(userId: string): Observable<FollowerUser[]> {
    return this.http.get<FollowerUser[]>(`${this.baseUrl}/${userId}/following`);
  }

  getRecommendations(userId: string): Observable<FollowerUser[]> {
    return this.http.get<FollowerUser[]>(`${this.baseUrl}/${userId}/recommendations`);
  }

  getFollowedAuthorIds(userId: string): Observable<string[]> {
    return this.http.get<string[]>(`${this.baseUrl}/${userId}/followed-author-ids`);
  }
}
