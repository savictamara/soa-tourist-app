import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { AuthResponse } from '../models/auth-response.model';
import { User } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private static readonly userStorageKey = 'tourist-app-current-user';
  private static readonly tokenStorageKey = 'tourist-app-token';
  private readonly currentUserSubject = new BehaviorSubject<User | null>(this.readUser());

  readonly currentUser$ = this.currentUserSubject.asObservable();

  get currentUser(): User | null {
    return this.currentUserSubject.value;
  }

  get token(): string | null {
    return localStorage.getItem(AuthStateService.tokenStorageKey);
  }

  get isAuthenticated(): boolean {
    return !!this.token && !!this.currentUser;
  }

  setSession(response: AuthResponse): void {
    localStorage.setItem(AuthStateService.userStorageKey, JSON.stringify(response.user));
    localStorage.setItem(AuthStateService.tokenStorageKey, response.token);
    this.currentUserSubject.next(response.user);
  }

  clear(): void {
    localStorage.removeItem(AuthStateService.userStorageKey);
    localStorage.removeItem(AuthStateService.tokenStorageKey);
    this.currentUserSubject.next(null);
  }

  private readUser(): User | null {
    const storedUser = localStorage.getItem(AuthStateService.userStorageKey);

    if (!storedUser) {
      return null;
    }

    try {
      return JSON.parse(storedUser) as User;
    } catch {
      localStorage.removeItem(AuthStateService.userStorageKey);
      return null;
    }
  }
}
