import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { RegisterRequest } from '../models/register-request.model';
import { User } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class StakeholdersApiService {
  private readonly baseUrl = 'http://localhost:5001/api';

  constructor(private readonly http: HttpClient) {}

  register(request: RegisterRequest): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/auth/register`, request);
  }

  getUsers(): Observable<User[]> {
    return this.http.get<User[]>(`${this.baseUrl}/users`);
  }
}
